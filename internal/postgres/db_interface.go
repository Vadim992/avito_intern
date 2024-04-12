package postgres

import "github.com/Vadim992/avito/internal/dto"

type DBModel interface {
	GetUserBanner(tagId, featureId, role int) (*dto.GetBanner, error)
	GetBanners(whereStmt, limitOffsetStmt string) ([]dto.GetBanner, error)
	InsertBanner(banner dto.PostPatchBanner) (int, error)
	UpdateBannerId(id int, banner dto.PostPatchBanner) (int, []int64, *dto.PostPatchBanner, error)
	DeleteBanner(id int) (int, []int64, error)
}
