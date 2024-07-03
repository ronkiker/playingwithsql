package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (db *Config) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to parse request body")
		return
	}

	fmt.Printf("Name Received: %v \n", params.Name)
	response, err := controller.CreateUsers(params.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to create user")
	}
	respondWithJSON(w, http.StatusCreated, response)

}
