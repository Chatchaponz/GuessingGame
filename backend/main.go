// go init backend
// go mod tidy
// go get github.com/gin-gonic/gin
// go get github.com/dgrijalva/jwt-go
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID       uint64 `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Guess struct {
	ID          uint64 `json:"user_id"`
	GuessNumber int64  `json:"guess_number"`
}

// var currentToken string
var hiddenNumber int64

var ourUser = User{ID: 1, Username: "testuser", Password: "1234"}

var router = gin.Default()

// Extract token from request header
func ExtractToken(request *http.Request) string {
	bearToken := request.Header.Get("Authorization")

	//fmt.Println(bearToken) // test print

	strArr := strings.Split(bearToken, " ")

	//fmt.Println(strArr) // test print

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

func getTokenData(request *http.Request) (uint64, error) {
	token, err := VerifyToken(request)
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		user_id, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)

		if err != nil {
			return 0, err
		}

		return user_id, nil
	}
	return 0, err
}

// Middleware for authentication token
func tokenMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := TokenValid(context.Request)
		if err != nil {
			context.JSON(http.StatusUnauthorized, err.Error())
			context.Abort()
			return
		}

		context.Next()
	}
}

func DoGuess(context *gin.Context) {
	var thisGuess *Guess

	if err := context.ShouldBindJSON(&thisGuess); err != nil {
		context.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}

	user_id, err := getTokenData(context.Request)
	if err != nil {
		context.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	// compare id in token with id of requested guess and id of our user
	if user_id != thisGuess.ID || user_id != ourUser.ID {
		context.JSON(http.StatusUnauthorized, "Unauthorized - ID doesn't match")
		return
	}

	var result string
	status := http.StatusOK

	if thisGuess.GuessNumber < hiddenNumber {
		result = "Too low"
	} else if thisGuess.GuessNumber > hiddenNumber {
		result = "Too high"
	} else {
		result = "Correct"
		hiddenNumber = rand.Int63n(101) // random new number
		status = http.StatusCreated     // http 201
	}

	context.JSON(status, result)
}

func Login(context *gin.Context) {
	var requestUser User

	if err := context.ShouldBindJSON(&requestUser); err != nil {
		context.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	// compare requested user with our user
	if ourUser.Username != requestUser.Username || ourUser.Password != requestUser.Password {
		context.JSON(http.StatusUnauthorized, "Please provide valid details")
		return
	}

	token, err := CreateToken(ourUser.ID)

	//currentToken = token

	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	context.JSON(http.StatusOK, token)
}

func CreateToken(userID uint64) (string, error) {
	var err error
	// Access Token
	os.Setenv("SECRET", "iwtptits")
	tokenClaims := jwt.MapClaims{}
	tokenClaims["authorized"] = true
	tokenClaims["user_id"] = userID
	tokenClaims["exp"] = time.Now().Add(time.Minute * 30).Unix() // time.Now().Add(time.Hour * 24).Unix()
	thisToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	token, err := thisToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return "", err
	}

	return token, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	hiddenNumber = rand.Int63n(101) // 0 <= n < 100

	router.POST("/login", Login)
	router.POST("/guess", tokenMiddleware(), DoGuess)

	log.Fatal(router.Run(":8080"))
}
