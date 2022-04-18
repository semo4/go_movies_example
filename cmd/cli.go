package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"

	"github.com/lib/pq"

	"github.com/semo4/go_movies_example/models"
)

var db *sql.DB

var dbErr error

var syncDataBaseConnection sync.Once

func init() {
	syncDataBaseConnection.Do(func() {

		err := godotenv.Load(".env")
		if err != nil {
			fmt.Println("Error :: ", err.Error())
		}
		host := os.Getenv("HOST")
		port, _ := strconv.ParseInt(os.Getenv("PORT"), 10, 64)
		user := os.Getenv("USER")
		password := os.Getenv("PASSWORD")
		dbname := os.Getenv("NAME")

		conn := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, int(port), user, password, dbname)
		db, dbErr = sql.Open("postgres", conn)

		if dbErr != nil {
			panic(dbErr)
		}
		if dbErr = db.Ping(); dbErr != nil {
			panic(dbErr)
		}
		fmt.Println("successfully connect to DB")
	})
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error :: ", err.Error())
	}
	movieSecretKey := flag.String("movieSecretKey", os.Getenv("MOVIE_API_KEY"), "The Api Key to use Movie Api")
	pageNumber := flag.Int("pageNumber", 5, "The number of page you need to fetch from api")

	flag.Parse()
	// var wg sync.WaitGroup
	// throttleChan := make(chan bool, 10)
	if *movieSecretKey != "" {
		for i := 1; i <= *pageNumber; i++ {
			// throttleChan <- true
			// wg.Add(1)

			// go func() {

			fmt.Println("Start with Page number ", i)
			client := http.Client{}

			link := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s&language=en-US&page=%s", *movieSecretKey, strconv.Itoa(i))

			req, err := http.NewRequest(http.MethodGet, link, nil)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			if err != nil {
				fmt.Printf("Error in Request:: %s", err.Error())
			}
			response, err := client.Do(req)

			if err != nil {
				fmt.Printf("Error When get response :: %s", err.Error())
			}

			body, err := ioutil.ReadAll(response.Body)

			if err != nil {
				fmt.Printf("Error In read Response :: %s", err.Error())
			}

			var moviesResponse models.MoviesData

			err = json.Unmarshal(body, &moviesResponse)
			if err != nil {
				fmt.Printf("Error in Convert to Json :: %s\n", err.Error())
			}

			InsertData(moviesResponse.Results)
			// }()

			// <-throttleChan
			// wg.Done()
		}
	}
	// wg.Wait()
}

func InsertData(result []models.Movie) {
	var wg sync.WaitGroup
	throttleChan := make(chan bool, 10)
	for _, item := range result {
		throttleChan <- true
		wg.Add(1)

		go func() {

			sqlStatement := `INSERT INTO movies (
				poster_path, adult, overview, release_date, genre_id, id, original_title, original_language, title, backdrop_path, popularity, vote_count, video, vote_average)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);`
			// id := 0
			_, dbErr := db.Exec(sqlStatement, item.PosterPath, item.Adult, item.Overview, item.ReleaseDate, pq.Array(item.GenreIds), item.Id, item.OriginalTitle, item.OriginalLanguage, item.Title, item.BackdropPath, item.Popularity, item.VoteCount, item.Video, item.VoteAverage)
			// QueryRow(sqlStatement, item.PosterPath, item.Adult, item.Overview, item.ReleaseDate, item.GenreId, item.Id, item.OriginalTitle, item.OriginalLanguage, item.Title, item.BackdropPath, item.Popularity, item.VoteCount, item.Video, item.VoteAverage).Scan(&id)
			if dbErr != nil {
				panic(dbErr.Error())
			}
			fmt.Println("New record ID is: ", item.Id)
		}()
		<-throttleChan
		wg.Done()

	}
	wg.Wait()
}
