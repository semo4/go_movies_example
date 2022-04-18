package models

type MoviesData struct {
	Page         int64   `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int64   `json:"total_pages"`
	TotalResults int64   `json:"total_results"`
}
