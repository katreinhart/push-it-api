package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katreinhart/push-it-api/model"
)

// CreateWorkout instantiates a new workout object with data from request body
func CreateWorkout(w http.ResponseWriter, r *http.Request) {
	uid, err := GetUIDFromBearerToken(r)

	if err != nil {
		handleErrorAndRespond(nil, model.ErrorForbidden, w)
	}

	// get the body from the request
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	var wk model.WorkoutModel
	err = json.Unmarshal(b, &wk)

	wk.UserID = uid

	_wk := model.CreateWorkout(wk)

	js, err := json.Marshal(_wk)

	handleErrorAndRespond(js, err, w)
}

// GetWorkout handles GET one http requests
func GetWorkout(w http.ResponseWriter, r *http.Request) {
	// get workout ID from vars
	vars := mux.Vars(r)
	id, _ := vars["id"]

	var wk model.WorkoutModel

	wk, err := model.GetWorkout(id)

	if err != nil {
		handleErrorAndRespond(nil, model.ErrorNotFound, w)
		return
	}

	if wk.Completed {
		cw, err := model.GetCompletedWorkout(wk)

		if err != nil {
			handleErrorAndRespond(nil, err, w)
			return
		}

		js, err := json.Marshal(cw)
		handleErrorAndRespond(js, err, w)
		return
	}

	js, err := json.Marshal(wk)

	handleErrorAndRespond(js, err, w)
}

// AddExerciseToWorkout handles http requests to POST workouts/id/exercises
func AddExerciseToWorkout(w http.ResponseWriter, r *http.Request) {
	// get workout ID from vars
	vars := mux.Vars(r)
	id, _ := vars["id"]

	// get the body from the request
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	var ep model.WorkoutExerciseAsPosted
	var e model.WorkoutExercise

	err := json.Unmarshal(b, &ep)

	if err != nil {
		handleErrorAndRespond(nil, model.ErrorBadRequest, w)
		return
	}

	e, err = model.AddExerciseToWorkout(id, ep)

	if err != nil {
		handleErrorAndRespond(nil, model.ErrorInternalServer, w)
		return
	}

	js, err := json.Marshal(e)

	handleErrorAndRespond(js, err, w)
}

// AddExerciseSet adds a set to the exercise
func AddExerciseSet(w http.ResponseWriter, r *http.Request) {
	// get workout ID and exercise ID from vars
	vars := mux.Vars(r)
	id, _ := vars["id"]

	// get the body from the request
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	var newSet model.WorkoutSet
	var newExSet model.WorkoutExerciseSet

	err := json.Unmarshal(b, &newSet)

	if err != nil {
		handleErrorAndRespond(nil, err, w)
		return
	}

	newExSet, err = model.AddExerciseSet(id, newSet)

	if err != nil {
		handleErrorAndRespond(nil, model.ErrorForbidden, w)
		return
	}
	js, err := json.Marshal(newExSet)

	handleErrorAndRespond(js, err, w)
}

// MarkWorkoutAsCompleted handles PUT requests to workout and marks as completed if appropriate
func MarkWorkoutAsCompleted(w http.ResponseWriter, r *http.Request) {

	// Get user ID from token
	uid, err := GetUIDFromBearerToken(r)
	if err != nil {
		handleErrorAndRespond(nil, errors.New("Forbidden"), w)
	}

	// get vars from url params
	vars := mux.Vars(r)
	id, _ := vars["id"]

	// get the body from the request
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	js, err := model.MarkWorkoutAsCompleted(uid, id, b)

	handleErrorAndRespond(js, err, w)
}

// UpdateWorkoutTimestamps handles requests to update start/end timestamps.
func UpdateWorkoutTimestamps(w http.ResponseWriter, r *http.Request) {
	// get workout ID from url params
	vars := mux.Vars(r)
	id, _ := vars["id"]

	// parse body from request
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	js, err := model.UpdateWorkoutTimestamps(id, b)

	handleErrorAndRespond(js, err, w)
}
