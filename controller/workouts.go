package controller

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katreinhart/push-it-api/model"
)

// CreateWorkout instantiates a new workout object with data from request body
func CreateWorkout(w http.ResponseWriter, r *http.Request) {
	uid, err := GetUIDFromBearerToken(r)
	if err != nil {
		handleErrorAndRespond(nil, errors.New("Forbidden"), w)
	}

	// get the body from the request
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	js, err := model.CreateWorkout(uid, b)

	handleErrorAndRespond(js, err, w)
}

// GetWorkout handles GET one http requests
func GetWorkout(w http.ResponseWriter, r *http.Request) {
	// get workout ID from vars
	vars := mux.Vars(r)
	id, _ := vars["id"]

	js, err := model.GetWorkout(id)

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

	js, err := model.AddExerciseToWorkout(id, b)

	handleErrorAndRespond(js, err, w)
}
