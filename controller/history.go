package controller

import (
	"encoding/json"
	"net/http"

	"github.com/katreinhart/push-it-api/model"
)

// History retrieves all completed workouts for the user in question
func History(w http.ResponseWriter, r *http.Request) {
	// Get user ID from token
	uid, err := GetUIDFromBearerToken(r)

	var _workouts []CompletedWorkout

	_workouts, err := model.History(uid)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(_workouts)

	handleErrorAndRespond(js, err, w)
}

func FetchSavedExercises(w http.ResponseWriter, r *http.Request) {
	uid, err := GetUIDFromBearerToken(r)

	js, err := model.History(uid)

	handleErrorAndRespond(js, err, w)
}
