package req

import (
	"encoding/json"
	"github.com/Vadim992/avito/internal/dto"
	"github.com/Vadim992/avito/internal/postgres"
	"github.com/Vadim992/avito/internal/storage"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func getIdFromPath(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return 0, PathIdErr
	}

	if id < 1 {
		return 0, PathIdErr
	}

	return id, nil
}

func PatchBannerId(db postgres.DBModel, inMemory storage.Storage, r *http.Request) error {
	id, err := getIdFromPath(r)

	if err != nil {
		return err
	}

	var postPatchBanner dto.PostPatchBanner

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&postPatchBanner); err != nil {
		return err
	}

	oldFeatureId, oldTagIds, banner, err := db.UpdateBannerId(id, postPatchBanner)

	if err != nil {
		return err
	}

	inMemory.Update(id, oldTagIds, oldFeatureId, &postPatchBanner, banner)

	return nil
}

func DeleteBannerId(db postgres.DBModel, inMemory storage.Storage, r *http.Request) error {
	id, err := getIdFromPath(r)

	if err != nil {
		return err
	}

	featureId, tagIds, err := db.DeleteBanner(id)

	if err != nil {
		return err
	}

	inMemory.Delete(id, featureId, tagIds)

	return nil
}
