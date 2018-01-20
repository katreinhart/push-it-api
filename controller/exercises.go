package controller

import (
	"bytes"
	"net/http"

	"github.com/katreinhart/push-it-api/model"
)

// ListExercises is index route for /api/exercises
func ListExercises(w http.ResponseWriter, r *http.Request) {
	js, err := model.ListExercises()
	handleErrorAndRespond(js, err, w)
}

// CreateExercise handles POST calls to /api/exercises
func CreateExercise(w http.ResponseWriter, r *http.Request) {
	// Create a new buffer to read the body, then parse into a []byte
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	b := []byte(buf.String())

	// Send the []byte b to the model and receive json and error
	js, err := model.CreateExercise(b)

	handleErrorAndRespond(js, err, w)
}
