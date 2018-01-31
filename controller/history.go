package controller

import (
	"net/http"

	"github.com/katreinhart/push-it-api/model"
)

// History retrieves all completed workouts for the user in question
func History(w http.ResponseWriter, r *http.Request) {
	// Get user ID from token
	uid, err := GetUIDFromBearerToken(r)

	js, err := model.History(uid)

	handleErrorAndRespond(js, err, w)
}

func FetchSavedExercises(w http.ResponseWriter, r *http.Request) {
	uid, err := GetUIDFromBearerToken(r)

	js, err := model.History(uid)

	handleErrorAndRespond(js, err, w)
}
