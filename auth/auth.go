package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type AuthMiddleware struct{}

var jwtSecretKey string

func GenerateJWTToken(userEmail string) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error :: ", err.Error())
	}
	jwtSecretKey = os.Getenv("JWT_SECRET_KEY")

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = userEmail
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	to := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := to.SignedString([]byte(jwtSecretKey))

	// tokenString, err := token.SignedString([]byte(jwtSecretKey))

	if err != nil {
		fmt.Printf("Somthing went wrong:: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func validateToken(tokenString string) (*jwt.Token, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error :: ", err.Error())
	}
	jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		fmt.Printf("Somthing went wrong:: %s", err.Error())
		return nil, err
	}

	return token, nil

}

var email string
var ok bool

func GetEmail() string {
	return email
}

func (amw *AuthMiddleware) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userOk := false
		authorization := r.Header.Get("Authorization")
		if authorization == " " {
			http.Error(w, "Invalid authorization", http.StatusUnauthorized)
			return

		}
		tokenParts := strings.Split(authorization, " ")
		if len(tokenParts) < 2 {
			http.Error(w, "unauthorized, you must pass an Authorization header with a valid JWT bearer token", http.StatusUnauthorized)
			return
		}

		token, err := validateToken(tokenParts[1])
		if err != nil {
			fmt.Printf("%s", err.Error())
			http.Error(w, "you must pass an Authorization header with a valid JWT bearer token", http.StatusUnauthorized)
			return
		}

		if email, ok = token.Claims.(jwt.MapClaims)["email"].(string); ok {
			if strings.HasSuffix(email, "@hotmail.com") || strings.HasSuffix(email, "@gmail.com") {
				userOk = true
			}
		}
		if !userOk {
			http.Error(w, "that user does not have access to call this service", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
