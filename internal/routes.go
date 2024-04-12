package internal

import (
	"github.com/Vadim992/avito/internal/mws"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
)

var pathRoles = map[string][]int{
	"/user_banner": {mws.ADMIN, mws.USER},
	"/banner":      {mws.ADMIN},
	"/banner/id":   {mws.ADMIN},
}

func (app *App) Routes() http.Handler {

	router := mux.NewRouter()

	router.HandleFunc("/user_banner", mws.Auth(pathRoles["/user_banner"],
		app.tokenMap, app.GetUserBanner)).Methods("GET")

	router.HandleFunc("/banner", mws.Auth(pathRoles["/banner"],
		app.tokenMap, app.GetBanners)).Methods("GET")
	router.HandleFunc("/banner", mws.Auth(pathRoles["/banner"],
		app.tokenMap, app.PostBanners)).Methods("POST")

	router.HandleFunc("/banner/{id:[0-9]+}", mws.Auth(pathRoles["/banner/id"],
		app.tokenMap, app.PatchBannerId)).Methods("PATCH")
	router.HandleFunc("/banner/{id:[0-9]+}", mws.Auth(pathRoles["/banner/id"],
		app.tokenMap, app.DeleteBannerId)).Methods("DELETE")

	mwChain := alice.New(mws.RecoverPanic, mws.LogRequest, mws.SetHeaders)

	return mwChain.Then(router)
}
