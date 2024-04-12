package req

import (
	"github.com/Vadim992/avito/internal/dto"
	"github.com/Vadim992/avito/internal/storage"
	"time"
)

type MockDB struct {
	BannersTable     map[storage.SearchIds]int
	BannersDataTable map[int]BannersDataDB
}

type BannersDataDB struct {
	BannerId int
	dto.BannerContent
	CreatedAt time.Time
	UpdatedAt time.Time
}
