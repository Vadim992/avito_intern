package req

import (
	"errors"
	"fmt"
	"github.com/Vadim992/avito/internal/dto"
	"github.com/Vadim992/avito/internal/mws"
	"github.com/Vadim992/avito/internal/postgres"
	"github.com/Vadim992/avito/internal/storage"
	"net/http"
	"time"
)

const bannerTimeConstraint = 10 * time.Second

func GetUsersBanner(db postgres.DBModel, inMemory storage.Storage,
	role int, w http.ResponseWriter, r *http.Request) error {
	queryParams := r.URL.Query()

	tagId, err := hasRequiredIntQuery(queryParams, tagIdQuery)

	if err != nil {
		return err
	}

	featureId, err := hasRequiredIntQuery(queryParams, featureIdQuery)

	if err != nil {
		return err
	}

	useLastRevision, err := hasBoolQuery(queryParams, useLastRevisionQuery)

	if err != nil {
		return err
	}

	if useLastRevision {
		res, err := db.GetUserBanner(tagId, featureId, role)

		if err != nil {
			return err
		}

		err = sentDataToFront(res.Content, w, http.StatusOK)
		if err != nil {
			return err
		}

		return nil
	}

	content, err := getDataFromStorage(inMemory, tagId, featureId, role)

	if err != nil {
		if !errors.Is(err, OldDataErr) {
			return err
		}

		fmt.Println(tagId, featureId, role)
		res, errDB := db.GetUserBanner(tagId, featureId, role)

		if errDB != nil {
			return errDB
		}

		updateInMemory(inMemory, res)

		content = *res.Content
	}

	err = sentDataToFront(content, w, http.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

func getDataFromStorage(inMemory storage.Storage, tagId, featureId, role int) (dto.BannerContent, error) {
	bannerInfo, err := inMemory.Get(tagId, featureId)

	if err != nil {
		return dto.BannerContent{}, err
	}

	err = validateBanner(bannerInfo, role)

	if err != nil {
		return dto.BannerContent{}, err
	}

	return bannerInfo.Content, nil
}

func validateBanner(bannerInfo *storage.BannerInfo, role int) error {
	if !bannerInfo.IsActive && role == mws.USER {
		return PermissionErr
	}

	t := bannerInfo.UpdatedAt
	t = t.Add(bannerTimeConstraint)

	if !t.After(time.Now()) {
		return OldDataErr
	}

	return nil
}

func updateInMemory(inMemory storage.Storage, banner *dto.GetBanner) {
	var b dto.PostPatchBanner

	b.Content = banner.Content
	b.IsActive = banner.IsActive

	inMemory.Update(*banner.BannerId, nil, 0, &dto.PostPatchBanner{}, &b)
}
