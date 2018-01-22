package controller

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/katreinhart/push-it-api/model"
)

// GetSecondaryGoals responds to requests for user's secondary goals
func GetSecondaryGoals(w http.ResponseWriter, r *http.Request) {
	// parse user ID from token
	uid, err := GetUIDFromBearerToken(r)
	if err != nil {
		// handle error
	}
	// Fetch secondary goals from model
	js, err := model.GetSecondaryGoals(uid)

	handleErrorAndRespond(js, err, w)
}

// PostSecondaryGoals responds to requests to post or update user's secondary goals
func PostSecondaryGoals(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post secondary goals CONTROLLER function")
	uid, err := GetUIDFromBearerToken(r)
	if err != nil {
		fmt.Println("Error with bearer token")
		// handle error
	}

	// get the body from the request
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	fmt.Println("Got stuff from body I think")
	// Fetch secondary goals from model
	js, err := model.PostSecondaryGoals(uid, b)
	handleErrorAndRespond(js, err, w)
}
