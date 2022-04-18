package models

type Movie struct {
	Id               int     `json:"id"`
	PosterPath       string  `json:"poster_path"`
	Adult            bool    `json:"adult"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
	OriginalTitle    string  `json:"original_title"`
	OriginalLanguage string  `json:"original_language"`
	Title            string  `json:"title"`
	BackdropPath     string  `json:"backdrop_path"`
	Popularity       float64 `json:"popularity"`
	VoteCount        int     `json:"vote_count"`
	Video            bool    `json:"video"`
	VoteAverage      float64 `json:"vote_average"`
	GenreIds         []int64 `json:"genre_ids"`
}
