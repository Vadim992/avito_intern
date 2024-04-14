package postgres

import (
	"github.com/Vadim992/avito/internal/dto"
	"github.com/Vadim992/avito/internal/storage"
)

type DBModel interface {
	GetUserBanner(tagId, featureId, role int) (*dto.GetBanner, error)
	GetBanners(whereStmt, limitOffsetStmt string) ([]dto.GetBanner, error)
	InsertBanner(banner dto.PostPatchBanner) (int, error)
	UpdateBannerId(id int, banner dto.PostPatchBanner) (*storage.UpdateDeleteFromDB, error)
	DeleteBanner(id int) (*storage.UpdateDeleteFromDB, error)
}
