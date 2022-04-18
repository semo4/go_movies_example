package utils

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

var DB *sql.DB
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
		DB, dbErr = sql.Open("postgres", conn)

		if dbErr != nil {
			panic(dbErr)
		}
		if err = DB.Ping(); err != nil {
			panic(err)
		}
		fmt.Println("successfully connect to DB")
	})
}
