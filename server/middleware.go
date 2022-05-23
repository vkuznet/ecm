package main

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	utils "github.com/vkuznet/ecm/utils"
)

// SigningKey returns unique signing key
func SigningKey() string {
	return fmt.Sprintf("%d", utils.MacAddress())
}

// isAuthorized middleware check incoming HTTP request token
// and authorize the user based on SigningKey
func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return SigningKey(), nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {

			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

// GenerateJWT generates JWT token
// Client should use the following format to place HTTP request:
//
//      token, err := GenerateJWT()
//      client := &http.Client{}
//      req, _ := http.NewRequest("GET", someURL, nil)
//      req.Header.Set("Token", token)
//
func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["client"] = "ECM client"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(SigningKey())

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
