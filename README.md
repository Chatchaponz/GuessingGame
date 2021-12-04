# GuessingGame

Golang -> backend
Package that I use
- go get github.com/gin-gonic/gin |---------- for routing
- go get github.com/dgrijalva/jwt-go |------- for JWT token
- go get github.com/gin-contrib/cors |------- for CORS config

React.js -> frontend
library that I add in
- nmp i react-router-dom |------- for routing
- npm i axios |------------------ for http request

# How to run
+ in folder "backend" run file main.go (go run main.go).
+ in folder "frontend" run (npm start).

# Game description
Test User : 
+ username: testuser
+ password: 1234

Features : 
+ Need to login before play. (for authentication and get JWT token)
+ Each token will last 20 minutes. As long as token are not expire player can play.
+ If token expire during the game, some certain action will bring player back to the login page like push guess button or refresh page.
+ Game will always told you wheather your guess is too high or low from the hidden number.
+ Every response, request, error between frontend and backend will be log in Console for learning and testing proposes.
+ Player can always logout anytime during the game.

I use around 20++ hours to learn everything from the scratch and there are a lot to learn more :)
