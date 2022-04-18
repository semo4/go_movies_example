package models

type User struct {
	Id             int     `json:"id"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	Email          string  `json:"email"`
	Password       string  `json:"password"`
	FavoriteMovies []int64 `json:"faviorate_movies"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FavoriteMovie struct {
	FavoriteMovies []int64 `json:"faviorate_movies"`
}
