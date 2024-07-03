package main

import "net/http"

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Status string `json:"status"`
	}
	respondWithJSON(w, 200, response{
		Status: "ok",
	})
}
