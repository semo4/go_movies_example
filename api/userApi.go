package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/semo4/go_movies_example/auth"
	"github.com/semo4/go_movies_example/models"
	"github.com/semo4/go_movies_example/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.LoginUser
	var message models.ResponseModel
	_ = json.NewDecoder(r.Body).Decode(&user)
	_, err := utils.DB.Query("select email, password from users where email = $1 and password = $2", user.Email, user.Password)
	if err != nil {
		http.Error(w, "Your Account not found", http.StatusNotFound)
	}
	token, err := auth.GenerateJWTToken(user.Email)
	if err != nil {
		http.Error(w, "Error ", http.StatusInternalServerError)
	}
	message.Message = "LogIn Successfully"
	message.Token = token
	json.NewEncoder(w).Encode(message)
	// json.Marshal(&message)
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)

	row, err := utils.DB.Query("select email from users where email = $1", user.Email)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	if row.Next() {
		result := row.Scan(&user.Email)
		if result != nil {
			http.Error(w, err.Error(), 500)
		}
		w.WriteHeader(204)
		w.Write([]byte("User Exist try with new Account"))

	} else {
		_, err := utils.DB.Query("insert into users (id, first_name, last_name, email, password, faviorate_movies) values ($1,$2,$3,$4,$5,$6)", user.Id, user.FirstName, user.LastName, user.Email, user.Password, pq.Array(user.FavoriteMovies))
		if err != nil {
			http.Error(w, err.Error(), 500)
			// http.Error(w, "Check your Entries", http.StatusBadRequest)
		}
		w.WriteHeader(201)
		w.Write([]byte("User Created Successfully"))
	}

	// json.Marshal(user)
}
func FavoriteMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	var movies []models.Movie
	row, err := utils.DB.Query("select faviorate_movies from users where email = ($1)", auth.GetEmail())
	if err != nil {
		fmt.Printf(">>>>>>>>%v\n", err.Error())
		http.Error(w, "User not found", http.StatusInternalServerError)
	}

	if row.Next() {

		data := row.Scan(pq.Array(&user.FavoriteMovies))
		if data != nil {
			http.Error(w, "FavoriteMovies not found", http.StatusInternalServerError)
		}
		for _, item := range user.FavoriteMovies {
			movieResult, err := utils.DB.Query("select * from movies where id = $1 ", item)
			if err != nil {
				w.WriteHeader(500)
				http.Error(w, "not found", http.StatusInternalServerError)
			}
			for movieResult.Next() {
				var movie models.Movie
				err := movieResult.Scan(&movie.Id, &movie.PosterPath, &movie.Adult, &movie.Overview, &movie.ReleaseDate, &movie.OriginalTitle, &movie.OriginalLanguage, &movie.Title, &movie.BackdropPath, &movie.Popularity, &movie.VoteCount, &movie.Video, &movie.VoteAverage, pq.Array(&movie.GenreIds))
				if err != nil {
					fmt.Printf(">>>>>>>>%v\n", err.Error())
				}
				movies = append(movies, movie)
			}
		}
	}

	json.NewEncoder(w).Encode(movies)
	// json.Marshal(movies)
}

func FavoriteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result := strings.Contains(params["id"], ",")
	if result {
		ids := strings.Split(params["id"], ",")
		var moviesRes []models.Movie
		for _, item := range ids {
			movieResult, err := utils.DB.Query("select * from movies where id = $1 ", item)
			if err != nil {
				w.WriteHeader(500)
				http.Error(w, "not found", http.StatusInternalServerError)
			}
			for movieResult.Next() {
				var movie models.Movie
				err := movieResult.Scan(&movie.Id, &movie.PosterPath, &movie.Adult, &movie.Overview, &movie.ReleaseDate, &movie.OriginalTitle, &movie.OriginalLanguage, &movie.Title, &movie.BackdropPath, &movie.Popularity, &movie.VoteCount, &movie.Video, &movie.VoteAverage, pq.Array(&movie.GenreIds))
				if err != nil {
					fmt.Printf(">>>>>>>>%v\n", err.Error())
					// fmt.Fprintf(w, "error")
				}
				moviesRes = append(moviesRes, movie)
			}
		}
		json.NewEncoder(w).Encode(moviesRes)

	} else {
		res := models.Movie{}
		row := utils.DB.QueryRow("select * from movies where id = $1", params["id"]).Scan(&res.Id, &res.PosterPath, &res.Adult, &res.Overview, &res.ReleaseDate, &res.OriginalTitle, &res.OriginalLanguage, &res.Title, &res.BackdropPath, &res.Popularity, &res.VoteCount, &res.Video, &res.VoteAverage, pq.Array(&res.GenreIds))
		if row != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		w.WriteHeader(200)
		// w.Write([]byte("OK"))
		json.NewEncoder(w).Encode(res)
	}

}

func AddFavoriteMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var favoriteMovie models.FavoriteMovie
	var favorite []int64
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&favoriteMovie)
	userEmail := auth.GetEmail()

	row, err := utils.DB.Query("select faviorate_movies from users where email = ($1)", userEmail)
	if err != nil {
		fmt.Printf(">>>>>>>>%v\n", err.Error())
		http.Error(w, "User not found", http.StatusInternalServerError)
	}

	if row.Next() {
		data := row.Scan(pq.Array(&user.FavoriteMovies))
		if data != nil {
			fmt.Printf(">>>>>>>>%v\n", data.Error())
			http.Error(w, "FavoriteMovies not found", http.StatusInternalServerError)
		}
		favorite = append(favorite, user.FavoriteMovies...)
		favorite = append(favorite, favoriteMovie.FavoriteMovies...)

		_, err := utils.DB.Query("update users set faviorate_movies = $1 where email = $2", pq.Array(favorite), userEmail)
		if err != nil {
			http.Error(w, err.Error(), 500)
			// http.Error(w, "Check your Entries", http.StatusBadRequest)
		}
	}
	w.WriteHeader(201)
	w.Write([]byte("Favorite Movie Added Successfully"))

}
