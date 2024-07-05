package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ronkiker/playingwithsql/blob/dev/internal/database"
)

type config struct {
	DB *database.Queries
}

func main() {
	const root = "."
	godotenv.Load()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal("PORT failed to set")
	}
	dbURL := os.Getenv("DB_URL")
	if len(dbURL) == 0 {
		log.Fatal("DB_URL failed to set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(errors.New("unable to open database"))
	}
	db := database.New(dbConn)
	config := config{
		DB: db,
	}

	go startScraper(db, 10, time.Minute)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", HandlerReadiness)
	v1Router.Get("/err", HandlerError)

	v1Router.Post("/feeds", config.authenticationService(config.HandlerCreateFeed))
	v1Router.Get("/feeds", config.HandlerGetFeeds)

	v1Router.Post("/users", config.HandleUserCreate)
	v1Router.Get("/users", config.authenticationService(config.HandleGetUser))

	v1Router.Post("/feed_follows", config.authenticationService(config.HandlerCreateFeedFollow))
	v1Router.Get("/feed_follows", config.authenticationService(config.HandlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowId}", config.authenticationService(config.HandlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Printf("Server starting on port %v \n", port)
	log.Fatal(srv.ListenAndServe())
}
