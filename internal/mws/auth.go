package mws

import (
	"net/http"
)

const TokenHeader = "token"

func Auth(roles []int, tokenMap map[string]int, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(TokenHeader)

		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		isValidRole := validateRole(roles, tokenMap, token)

		if !isValidRole {
			w.WriteHeader(http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, r)
	}
}

func validateRole(roles []int, tokenMap map[string]int, token string) bool {
	curRole := tokenMap[token]

	for _, role := range roles {
		if curRole == role {
			return true
		}
	}

	return false
}
