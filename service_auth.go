package main

import (
	"net/http"

	"github.com/ronkiker/playingwithsql/blob/dev/authentication"
	"github.com/ronkiker/playingwithsql/blob/dev/internal/database"
)

type authenticationHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *config) authenticationService(handler authenticationHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := authentication.GetApiKey(r.Header)
		if err != nil {
			respondWithError(w, 404, err.Error())
			return
		}
		user, err := cfg.DB.GetUserByApiKey(r.Context(), token)
		if err != nil {
			respondWithError(w, 400, err.Error())
			return
		}
		handler(w, r, user)
	}
}
