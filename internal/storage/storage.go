package storage

import (
	"errors"
	"github.com/Vadim992/avito/internal/dto"
	"sync"
	"time"
)

var (
	SaveErr   = errors.New("failed save data in storage")
	ReadErr   = errors.New("no data in storage")
	UpdateErr = errors.New("cannot update data in storage")
	DeleteErr = errors.New("cannot delete data from storage")
)

type UpdateDeleteFromDB struct {
	BannerId      int
	FeatureId     *int
	TagIds        []int64
	ReqBanner     *dto.PostPatchBanner
	UpdatedBanner *dto.PostPatchBanner
}

func NewUpdateDeleteFromDB(bannerId int, featureId *int, tagIds []int64,
	reqBanner, updatedBanner *dto.PostPatchBanner) *UpdateDeleteFromDB {
	return &UpdateDeleteFromDB{
		BannerId:      bannerId,
		FeatureId:     featureId,
		TagIds:        tagIds,
		ReqBanner:     reqBanner,
		UpdatedBanner: updatedBanner,
	}
}

type Storage interface {
	Save(bannerId int, banner *dto.PostPatchBanner) error
	Get(tagId, featureId int) (*BannerInfo, error)
	Update(banner *UpdateDeleteFromDB) error
	Delete(banner *UpdateDeleteFromDB) error
}

func NewStorage() Storage {
	return NewInMemoryStorage()
}

type SearchIds struct {
	TagId     int
	FeatureId int
}

func NewSearchIds(tagId, featureId int) SearchIds {
	return SearchIds{
		TagId:     tagId,
		FeatureId: featureId,
	}
}

type BannerInfo struct {
	BannerId  int
	Content   dto.BannerContent
	IsActive  bool
	UpdatedAt time.Time
}

func NewBannersInfo(id int, content dto.BannerContent, isActive bool) *BannerInfo {
	return &BannerInfo{
		BannerId:  id,
		Content:   content,
		UpdatedAt: time.Now(),
		IsActive:  isActive,
	}
}

type InMemoryStorage struct {
	SearchStorage map[SearchIds]int
	Banners       map[int]*BannerInfo
	mu            sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		SearchStorage: make(map[SearchIds]int),
		Banners:       make(map[int]*BannerInfo),
	}
}

func (s *InMemoryStorage) Save(bannerId int, banner *dto.PostPatchBanner) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if bannerId <= 0 || banner == nil {
		return SaveErr
	}

	if banner.Content == nil || banner.IsActive == nil {
		return SaveErr
	}

	if banner.FeatureId == nil || len(banner.TagIds) == 0 {
		return SaveErr
	}

	s.Banners[bannerId] = NewBannersInfo(bannerId, *banner.Content, *banner.IsActive)

	for _, tagId := range banner.TagIds {
		searchIds := NewSearchIds(int(tagId), *banner.FeatureId)

		s.SearchStorage[searchIds] = bannerId
	}

	return nil
}

func (s *InMemoryStorage) Get(tagId, featureId int) (*BannerInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	searchIds := NewSearchIds(tagId, featureId)
	val, ok := s.SearchStorage[searchIds]

	if !ok {
		return nil, ReadErr
	}

	return s.Banners[val], nil
}

func validateStructForUpdate(storageStruct *UpdateDeleteFromDB) error {
	if storageStruct == nil {
		return UpdateErr
	}

	if storageStruct.BannerId <= 0 {
		return UpdateErr
	}

	if storageStruct.FeatureId != nil {
		if *storageStruct.FeatureId <= 0 {
			return UpdateErr
		}
	}

	return nil
}

func (s *InMemoryStorage) Update(storageStruct *UpdateDeleteFromDB) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validateStructForUpdate(storageStruct); err != nil {
		return err
	}

	reqBanner := storageStruct.ReqBanner
	if reqBanner != nil && len(storageStruct.TagIds) != 0 && storageStruct.FeatureId != nil {
		switch {
		case reqBanner.TagIds != nil && reqBanner.FeatureId != nil:
			s.deleteFromSearchStorage(storageStruct.TagIds, *storageStruct.FeatureId)

			s.fillSearchStorage(reqBanner.TagIds, *reqBanner.FeatureId, storageStruct.BannerId)
		case reqBanner.TagIds != nil:
			s.deleteFromSearchStorage(storageStruct.TagIds, *storageStruct.FeatureId)

			s.fillSearchStorage(reqBanner.TagIds, *storageStruct.FeatureId, storageStruct.BannerId)
		case reqBanner.FeatureId != nil:
			s.deleteFromSearchStorage(storageStruct.TagIds, *storageStruct.FeatureId)

			s.fillSearchStorage(storageStruct.TagIds, *reqBanner.FeatureId, storageStruct.BannerId)
		}

	}

	item := s.Banners[storageStruct.BannerId]

	updatedBanner := storageStruct.UpdatedBanner

	if updatedBanner != nil {

		if updatedBanner.IsActive != nil {
			item.IsActive = *updatedBanner.IsActive
		}

		if updatedBanner.Content != nil {
			item.Content = *updatedBanner.Content
		}

		item.UpdatedAt = time.Now()
		s.Banners[storageStruct.BannerId] = item
	}

	return nil
}

func (s *InMemoryStorage) deleteFromSearchStorage(oldTags []int64, oldFeatureId int) {
	for _, tagId := range oldTags {
		delete(s.SearchStorage, NewSearchIds(int(tagId), oldFeatureId))
	}
}

func (s *InMemoryStorage) fillSearchStorage(tagIds []int64, featureId, id int) {
	for _, tagId := range tagIds {
		s.SearchStorage[NewSearchIds(int(tagId), featureId)] = id
	}
}

func validateStructForDelete(storageStruct *UpdateDeleteFromDB) error {
	if storageStruct == nil {
		return DeleteErr
	}

	if storageStruct.BannerId <= 0 || storageStruct.FeatureId == nil || *storageStruct.FeatureId <= 0 {
		return DeleteErr
	}

	if len(storageStruct.TagIds) == 0 {
		return DeleteErr
	}

	return nil
}

func (s *InMemoryStorage) Delete(storageStruct *UpdateDeleteFromDB) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validateStructForDelete(storageStruct); err != nil {
		return err
	}

	delete(s.Banners, storageStruct.BannerId)

	for _, tagId := range storageStruct.TagIds {
		delete(s.SearchStorage, NewSearchIds(int(tagId), *storageStruct.FeatureId))
	}

	return nil
}
