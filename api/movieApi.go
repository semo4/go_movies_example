package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/semo4/go_movies_example/models"
	"github.com/semo4/go_movies_example/utils"
)

func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var moviesRes []models.Movie

	rows, err := utils.DB.Query("select * from movies")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		res := models.Movie{}
		err := rows.Scan(&res.Id, &res.PosterPath, &res.Adult, &res.Overview, &res.ReleaseDate, &res.OriginalTitle, &res.OriginalLanguage, &res.Title, &res.BackdropPath, &res.Popularity, &res.VoteCount, &res.Video, &res.VoteAverage, pq.Array(&res.GenreIds))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		moviesRes = append(moviesRes, res)

	}
	fmt.Printf("%+v", moviesRes)
	w.WriteHeader(200)
	// json.Marshal(moviesRes)
	json.NewEncoder(w).Encode(moviesRes)
}

func GetMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result := strings.Contains(params["id"], ",")
	if result {
		ids := strings.Split(params["id"], ",")
		rows, err := utils.DB.Query("select * from movies where id in %s", ids)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var moviesRes []models.Movie

		for rows.Next() {
			res := models.Movie{}
			err := rows.Scan(&res.Id, &res.PosterPath, &res.Adult, &res.Overview, &res.ReleaseDate, &res.OriginalTitle, &res.OriginalLanguage, &res.Title, &res.BackdropPath, &res.Popularity, &res.VoteCount, &res.Video, &res.VoteAverage, pq.Array(&res.GenreIds))
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), 500)
				return
			}
			moviesRes = append(moviesRes, res)

		}
		w.WriteHeader(200)
		w.Write([]byte("OK"))
		// json.Marshal(moviesRes)
		json.NewEncoder(w).Encode(moviesRes)

	} else {
		res := models.Movie{}
		row := utils.DB.QueryRow("select * from movies where id = $1", params["id"]).Scan(&res.Id, &res.PosterPath, &res.Adult, &res.Overview, &res.ReleaseDate, &res.OriginalTitle, &res.OriginalLanguage, &res.Title, &res.BackdropPath, &res.Popularity, &res.VoteCount, &res.Video, &res.VoteAverage, pq.Array(&res.GenreIds))
		if row != nil {
			log.Println(row)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.WriteHeader(200)
		// w.Write([]byte("OK"))
		json.NewEncoder(w).Encode(res)
		// json.Marshal(res)
	}
}
