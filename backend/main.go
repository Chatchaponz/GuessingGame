// go init backend
// go mod tidy
// go get github.com/gin-gonic/gin
// go get github.com/dgrijalva/jwt-go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Guess struct {
	Username    string `json:"username"`
	GuessNumber int64  `json:"guess_number"`
}

var currentToken string

var ourUser = User{Username: "testuser", Password: "1234"}

var router = gin.Default()

// Extract token from request header
func ExtractToken(request *http.Request) string {
	bearToken := request.Header.Get("Authorization")

	fmt.Println(bearToken) // test print

	strArr := strings.Split(bearToken, " ")

	fmt.Println(strArr) // test print

	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func VerifyToken(request *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(request)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenValid(request *http.Request) error {
	token, err := VerifyToken(request)

	if err != nil {
		return err
	}

	// _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid
	if !token.Valid {
		return err
	}

	return nil
}

func getTokenData(request *http.Request) (string, error) {
	token, err := VerifyToken(request)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			return "", err
		}

		return username, nil
	}
	return "", err
}

func DoGuess(context *gin.Context) {
	var thisGuess *Guess

	if err := context.ShouldBindJSON(&thisGuess); err != nil {
		context.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}

	username, err := getTokenData(context.Request)
	if err != nil {
		context.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	thisGuess.Username = username
	thisGuess.GuessNumber = 3

	context.JSON(http.StatusOK, thisGuess)
}

func Login(context *gin.Context) {
	var requestUser User

	if err := context.ShouldBindJSON(&requestUser); err != nil {
		context.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	if ourUser.Username != requestUser.Username || ourUser.Password != requestUser.Password {
		context.JSON(http.StatusUnauthorized, "Please provide valid details")
		return
	}

	token, err := CreateToken(requestUser.Username)

	currentToken = token

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	context.JSON(http.StatusOK, currentToken)
}

func CreateToken(username string) (string, error) {
	var err error
	// Access Token
	os.Setenv("SECRET", "iwtptits")
	tokenClaims := jwt.MapClaims{}
	tokenClaims["authorized"] = true
	tokenClaims["username"] = username
	tokenClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	thisToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	token, err := thisToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", err
	}

	return token, nil
}

func main() {
	router.POST("/login", Login)
	router.POST("/guess", DoGuess)

	log.Fatal(router.Run(":4242"))
}
