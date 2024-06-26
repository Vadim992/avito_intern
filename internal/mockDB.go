package internal

import (
	"database/sql"
	"fmt"
	"github.com/Vadim992/avito/internal/dto"
	"github.com/Vadim992/avito/internal/mws"
	"github.com/Vadim992/avito/internal/postgres"
	"github.com/Vadim992/avito/internal/storage"
	"time"
)

type MockDB struct {
	BannersTable     map[storage.SearchIds]int
	BannersDataTable map[int]*BannersDataDB
}

type BannersDataDB struct {
	BannerId  int
	Content   dto.BannerContent
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBannersDataDB(id int, title, text, url string, isActive bool) *BannersDataDB {
	content := dto.NewBannerContent(title, text, url)

	return &BannersDataDB{
		BannerId:  id,
		Content:   content,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewMockDB() MockDB {
	bannersData := make(map[int]*BannersDataDB, 5)
	fillBannersData(bannersData, 5)

	bannersTable := make(map[storage.SearchIds]int, 25)
	fillBannersTable(bannersTable, len(bannersData), len(bannersData))

	return MockDB{
		BannersTable:     bannersTable,
		BannersDataTable: bannersData,
	}
}

func fillBannersData(m map[int]*BannersDataDB, num int) {
	for i := 1; i <= num; i++ {
		title := fmt.Sprintf("title%d", i)
		text := fmt.Sprintf("text%d", i)
		url := fmt.Sprintf("url%d", i)

		isActive := true

		if i%2 == 0 {
			isActive = false
		}

		b := NewBannersDataDB(i, title, text, url, isActive)

		m[i] = b
	}
}

func fillBannersTable(m map[storage.SearchIds]int, numTag, numFeature int) {
	for i := 1; i <= numFeature; i++ {
		for j := 1; j <= numTag; j++ {
			searchIds := storage.NewSearchIds(j, i)

			m[searchIds] = i
		}

	}
}

func (m *MockDB) FillInMemory(inMemory *storage.InMemoryStorage) {
	for key, val := range m.BannersTable {
		inMemory.SearchStorage[key] = val
	}

	for key, val := range m.BannersDataTable {
		oldContent := val.Content
		newContent := dto.NewBannerContent(*oldContent.Title,
			*oldContent.Text, *oldContent.Url)
		b := storage.NewBannersInfo(val.BannerId, newContent,
			val.IsActive)

		inMemory.Banners[key] = b
	}
}

func (m *MockDB) GetUserBanner(tagId, featureId, role int) (*dto.GetBanner, error) {
	searchIds := storage.NewSearchIds(tagId, featureId)

	bannerId, ok := m.BannersTable[searchIds]
	if !ok {
		return nil, sql.ErrNoRows
	}

	bannerData, ok := m.BannersDataTable[bannerId]
	if !ok {
		return nil, sql.ErrNoRows
	}

	if !bannerData.IsActive && role == mws.USER {
		return nil, postgres.PermissionErr
	}

	res := dto.GetBanner{BannerId: &bannerData.BannerId}
	res.IsActive = &bannerData.IsActive
	c := bannerData.Content
	content := dto.NewBannerContent(*c.Title, *c.Text, *c.Url)
	res.Content = &content

	return &res, nil
}

func (m *MockDB) GetBanners(whereStmt, limitOffsetStmt string) ([]dto.GetBanner, error) {
	return nil, nil
}

func (m *MockDB) InsertBanner(banner dto.PostPatchBanner) (int, error) {
	return 0, nil
}

func (m *MockDB) UpdateBannerId(id int, banner dto.PostPatchBanner) (*storage.UpdateDeleteFromDB, error) {
	return nil, nil
}
func (m *MockDB) DeleteBanner(id int) (*storage.UpdateDeleteFromDB, error) {
	return nil, nil
}
