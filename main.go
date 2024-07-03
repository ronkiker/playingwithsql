package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ronkiker/playingwithsql/blob/dev/internal/database"
)

type Config struct {
	DB *database.Queries
}

func main() {
	const root = "."
	godotenv.Load()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal("PORT failed to set")
	}
	db, err := sql.Open("postgres", "DB_URL")
	if err != nil {
		log.Fatal(errors.New("unable to open database"))
	}
	dbQueries := database.New(db)

	config := Config{
		DB: dbQueries,
	}
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
	v1Router.HandleFunc("/healthz", HandlerReadiness)
	v1Router.HandleFunc("/err", HandlerError)

	v1Router.Post("/users", config.HandleUserCreate)

	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Printf("Server starting on port %v \n", port)
	log.Fatal(srv.ListenAndServe())
}
