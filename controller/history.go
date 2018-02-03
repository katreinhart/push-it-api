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

	var _workouts []model.CompletedWorkout

	_workouts, err = model.History(uid)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(_workouts)

	handleErrorAndRespond(js, err, w)
}

// FetchSavedWorkouts returns all workouts that have not been completed
func FetchSavedWorkouts(w http.ResponseWriter, r *http.Request) {
	uid, err := GetUIDFromBearerToken(r)
	if err != nil {
		handleErrorAndRespond(nil, model.ErrorForbidden, w)
		return
	}

	workouts, err := model.FetchSavedWorkouts(uid)
	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	js, err := json.Marshal(workouts)

	handleErrorAndRespond(js, err, w)
}
