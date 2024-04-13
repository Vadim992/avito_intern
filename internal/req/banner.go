package req

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Vadim992/avito/internal/dto"
	"github.com/Vadim992/avito/internal/postgres"
	"github.com/Vadim992/avito/internal/storage"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	equal = "="
	space = " "
)

func GetBanners(db postgres.DBModel, w http.ResponseWriter, r *http.Request) error {
	queryParams := r.URL.Query()

	var whereStmt strings.Builder

	tagId, ok, err := hasIntQuery(queryParams, tagIdQuery)

	if ok {
		if err != nil {
			return err
		}

		postgres.CreateWhereReq(&whereStmt, tagIdQuery, strconv.Itoa(tagId), equal)
	}

	featureId, ok, err := hasIntQuery(queryParams, featureIdQuery)

	if ok {
		if err != nil {
			return err
		}

		postgres.CreateWhereReq(&whereStmt, featureIdQuery, strconv.Itoa(featureId), equal)
	}

	whereStmtStr := strings.TrimSpace(whereStmt.String())

	var limitOffsetStmt strings.Builder

	limit, ok, err := hasIntQuery(queryParams, limitQuery)

	if ok {
		if err != nil {
			return err
		}

		postgres.CreateLimitOffsetReq(&limitOffsetStmt, limitQuery, limit)
	}

	offset, ok, err := hasIntQuery(queryParams, offsetQuery)

	if ok {
		if err != nil {
			return err
		}

		postgres.CreateLimitOffsetReq(&limitOffsetStmt, offsetQuery, offset)
	}

	limitOffsetStmtStr := strings.TrimSpace(limitOffsetStmt.String())

	data, err := db.GetBanners(r.Context(), whereStmtStr, limitOffsetStmtStr)

	if err != nil {
		fmt.Println(errors.Is(err, context.Canceled))
		return err
	}

	err = sentDataToFront(data, w, http.StatusOK)
	if err != nil {
		return err
	}

	return nil
}

func PostBanner(db postgres.DBModel, inMemory storage.Storage, w http.ResponseWriter, r *http.Request) error {
	var postPatchBanner dto.PostPatchBanner

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&postPatchBanner); err != nil {
		return err
	}

	if err := checkFillPostBanner(postPatchBanner); err != nil {
		return err
	}

	if len(postPatchBanner.TagIds) == 0 {
		return postgres.EmptyArrTagIdsErr
	}

	postPatchBanner.TagIds = deleteDuplicateFromTagIds(postPatchBanner.TagIds)

	bannerId, err := db.InsertBanner(postPatchBanner)

	if err != nil {
		return err
	}

	err = inMemory.Save(bannerId, &postPatchBanner)

	if err != nil {
		return err
	}

	bannerIdStruct := dto.NewBannerId(bannerId)

	err = sentDataToFront(bannerIdStruct, w, http.StatusCreated)
	if err != nil {
		return err
	}

	return nil
}

func checkFillPostBanner(banner interface{}) error {
	val := reflect.ValueOf(banner)

	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).IsNil() {
			return EmptyFieldErr
		}

		if val.Field(i).Kind() == reflect.Ptr {
			field := val.Field(i).Elem()

			if field.Kind() == reflect.Struct {
				err := checkFillPostBanner(field.Interface())

				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
