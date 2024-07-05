package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/ronkiker/playingwithsql/blob/dev/internal/database"
)

func (cfg *config) HandlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unable to parse request body")
		return
	}

	follow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})
	if err != nil {
		respondWithError(w, 400, "unable to create feed follow")
		return
	}
	respondWithJSON(w, 201, follow)
}

func (cfg *config) HandlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting feed follows for user: %v \n", err))
		return
	}
	respondWithJSON(w, 201, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (cfg *config) HandlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollow := chi.URLParam(r, "feedFollowId")
	ffID, err := uuid.Parse(feedFollow)
	if err != nil {
		respondWithError(w, 400, "invalid feed follow id")
		return
	}
	er := cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     ffID,
		UserID: user.ID,
	})
	if er != nil {
		respondWithError(w, 400, fmt.Sprintf("Error locating feed follow: %v \n", err))
		return
	}
	respondWithJSON(w, 201, struct{}{})
}
