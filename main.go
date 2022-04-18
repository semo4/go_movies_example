package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/semo4/go_movies_example/api"
	"github.com/semo4/go_movies_example/auth"
)

func main() {

	var authMiddleware auth.AuthMiddleware
	router := mux.NewRouter()

	router.HandleFunc("/api/movies", api.GetMovies).Methods(http.MethodGet)
	router.HandleFunc("/api/movies/{id}", api.GetMovie).Methods(http.MethodGet)
	router.HandleFunc("/api/login", api.Login).Methods(http.MethodPost)
	router.HandleFunc("/api/register", api.Register).Methods(http.MethodPost)

	router.Handle("/api/favorite_movies", authMiddleware.IsAuthorized(http.HandlerFunc(api.FavoriteMovies))).Methods(http.MethodGet)
	router.Handle("/api/favorite_movies/{id}", authMiddleware.IsAuthorized(http.HandlerFunc(api.FavoriteMovie))).Methods(http.MethodGet)
	router.Handle("/api/add_favorite_movies", authMiddleware.IsAuthorized(http.HandlerFunc(api.AddFavoriteMovies))).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(":8000", router))

}
