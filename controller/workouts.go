package controller

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/katreinhart/push-it-api/model"
)

// CreateWorkout instantiates a new workout object with data from request body
func CreateWorkout(w http.ResponseWriter, r *http.Request) {
	uid, err := GetUIDFromBearerToken(r)
	fmt.Println("UID", uid)
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

// AddExerciseSet adds a set to the exercise
func AddExerciseSet(w http.ResponseWriter, r *http.Request) {
	// get workout ID and exercise ID from vars
	vars := mux.Vars(r)
	id, _ := vars["id"]
	eid, _ := vars["eid"]

	// get the body from the request
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	js, err := model.AddExerciseSet(id, eid, b)

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
