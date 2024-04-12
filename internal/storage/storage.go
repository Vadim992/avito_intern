package storage

import (
	"database/sql"
	"fmt"
	"github.com/Vadim992/avito/internal/dto"
	"sync"
	"time"
)

type Storage interface {
	Save(bannerId int, banner *dto.PostPatchBanner)
	Get(tagId, featureId int) (*BannerInfo, error)
	Update(id int, oldTags []int64, oldFeatureId int,
		reqBanner, banner *dto.PostPatchBanner)
	Delete(bannerId, featureId int, tags []int64)
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
	Content   dto.BannerContent
	IsActive  bool
	UpdatedAt time.Time
}

func NewBannersInfo(content dto.BannerContent, isActive bool) *BannerInfo {
	return &BannerInfo{
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

func (s *InMemoryStorage) Save(bannerId int, banner *dto.PostPatchBanner) {
	//don't need check exist these data in storage or not, because check it earlier, when add data to DB
	// Save Must be called AFTER insert in DB
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Banners[bannerId] = NewBannersInfo(*banner.Content, *banner.IsActive)

	for _, tagId := range banner.TagIds {
		searchIds := NewSearchIds(int(tagId), *banner.FeatureId)

		s.SearchStorage[searchIds] = bannerId
	}
}

func (s *InMemoryStorage) Get(tagId, featureId int) (*BannerInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	searchIds := NewSearchIds(tagId, featureId)
	val, ok := s.SearchStorage[searchIds]

	if !ok {
		fmt.Println(val, ok)
		return nil, sql.ErrNoRows
	}

	return s.Banners[val], nil
}

func (s *InMemoryStorage) Update(id int, oldTags []int64, oldFeatureId int,
	reqBanner, banner *dto.PostPatchBanner) {
	// all checks were when update data in DB
	// Update InMemory MUST be after update db
	s.mu.Lock()
	defer s.mu.Unlock()

	s.deleteFromSearchStorage(oldTags, oldFeatureId)

	switch {
	case reqBanner.TagIds != nil && reqBanner.FeatureId != nil:
		s.fillSearchStorage(reqBanner.TagIds, *reqBanner.FeatureId, id)
	case reqBanner.TagIds != nil:
		s.fillSearchStorage(reqBanner.TagIds, oldFeatureId, id)
	case reqBanner.FeatureId != nil:
		s.fillSearchStorage(oldTags, *reqBanner.FeatureId, id)
	}

	item := s.Banners[id]
	if banner.IsActive != nil {
		item.IsActive = *banner.IsActive
	}

	if banner.Content != nil {
		item.Content = *banner.Content
	}

	item.UpdatedAt = time.Now()
	s.Banners[id] = item
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

func (s *InMemoryStorage) Delete(bannerId, featureId int, tags []int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.Banners, bannerId)

	for _, tagId := range tags {
		delete(s.SearchStorage, NewSearchIds(int(tagId), featureId))
	}
}
