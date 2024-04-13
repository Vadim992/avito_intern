package internal

import (
	"database/sql"
	"errors"
	"github.com/Vadim992/avito/internal/mws"
	"github.com/Vadim992/avito/internal/postgres"
	"github.com/Vadim992/avito/internal/req"
	"github.com/Vadim992/avito/internal/storage"
	"io"
	"net/http"
)

type App struct {
	db              postgres.DBModel
	inMemoryStorage storage.Storage
	tokenMap        map[string]int
}

func NewApp(db postgres.DBModel, inMemoryStorage storage.Storage, m map[string]int) *App {
	return &App{
		db:              db,
		inMemoryStorage: inMemoryStorage,
		tokenMap:        m,
	}
}

func (app *App) GetUserBanner(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get(mws.TokenHeader)
	role := app.tokenMap[token]
	err := req.GetUsersBanner(app.db, app.inMemoryStorage, role, w, r)

	if err != nil {
		switch {
		case errors.Is(err, req.QueryDataErr):
			ClientErr(w, http.StatusBadRequest, err)
		case errors.Is(err, req.PermissionErr):
			ClientErr(w, http.StatusForbidden, err)
		case errors.Is(err, sql.ErrNoRows):
			ClientErr(w, http.StatusNotFound, err)
		case errors.Is(err, storage.ReadErr):
			ClientErr(w, http.StatusNotFound, err)
		default:
			ServerErr(w, err)
		}
	}
}

func (app *App) GetBanners(w http.ResponseWriter, r *http.Request) {
	err := req.GetBanners(app.db, w, r)

	if err != nil {
		switch {
		case errors.Is(err, req.QueryDataErr):
			ClientErr(w, http.StatusBadRequest, err)
		case errors.Is(err, sql.ErrNoRows):
			ClientErr(w, http.StatusNotFound, err)
		default:
			ServerErr(w, err)
		}
	}
}

func (app *App) PostBanners(w http.ResponseWriter, r *http.Request) {
	err := req.PostBanner(app.db, app.inMemoryStorage, w, r)

	if err != nil {
		err = req.ValidateDriverErrors(err)

		switch {
		case errors.Is(err, io.EOF):
			ClientErr(w, http.StatusBadRequest, err)
		case errors.Is(err, req.EmptyFieldErr):
			ClientErr(w, http.StatusBadRequest, err)
		case errors.Is(err, postgres.EmptyArrTagIdsErr):
			ClientErr(w, http.StatusBadRequest, err)
		case errors.Is(err, postgres.InsertUpdateBannerErr):
			ClientErr(w, http.StatusBadRequest, err)
		default:
			ServerErr(w, err)
		}
	}

}

func (app *App) PatchBannerId(w http.ResponseWriter, r *http.Request) {
	err := req.PatchBannerId(app.db, app.inMemoryStorage, r)

	if err != nil {
		err = req.ValidateDriverErrors(err)

		switch {
		case errors.Is(err, req.PathIdErr):
			ClientErr(w, http.StatusBadRequest, err)
		case errors.Is(err, io.EOF):
			ClientErr(w, http.StatusBadRequest, err)
		case errors.Is(err, sql.ErrNoRows):
			ClientErr(w, http.StatusNotFound, err)
		case errors.Is(err, postgres.InsertUpdateBannerErr):
			ClientErr(w, http.StatusBadRequest, err)
		default:
			ServerErr(w, err)
		}

	}
}
func (app *App) DeleteBannerId(w http.ResponseWriter, r *http.Request) {
	err := req.DeleteBannerId(app.db, app.inMemoryStorage, w, r)

	if err != nil {
		switch {
		case errors.Is(err, req.PathIdErr):
			ClientErr(w, http.StatusBadRequest, err)
		case errors.Is(err, sql.ErrNoRows):
			ClientErr(w, http.StatusNotFound, err)
		default:
			ServerErr(w, err)
		}

	}
}
