package dto

import "time"

type BannerContent struct {
	Title *string `json:"title"`
	Text  *string `json:"text"`
	Url   *string `json:"url"`
}

func NewBannerContent(title, text, url string) BannerContent {
	return BannerContent{
		Title: &title,
		Text:  &text,
		Url:   &url,
	}
}

type PostPatchBanner struct {
	FeatureId *int           `json:"feature_id"`
	TagIds    []int64        `json:"tag_ids"`
	Content   *BannerContent `json:"content"`
	IsActive  *bool          `json:"is_active"`
}

type GetBanner struct {
	BannerId *int `json:"banner_id"`
	PostPatchBanner
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type BannerId struct {
	BannerId int `json:"banner_id"`
}

func NewBannerId(bannerId int) *BannerId {
	return &BannerId{
		BannerId: bannerId,
	}
}
