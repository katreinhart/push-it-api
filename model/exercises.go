package model

import (
	"encoding/json"
)

// ListExercises returns all available exercises
func ListExercises() ([]byte, error) {

	var exercises []exercise
	// var _exercises []transformedExercise

	db.Find(&exercises)

	if len(exercises) == 0 {
		// return not found error
	}

	// transform exercises

	// marshal into json

	// return js, err
	return []byte(""), nil
}

// func FetchSingleExercise(id uint) ([]byte, error) {

// }

// CreateExercise adds an exercise to the database.
func CreateExercise(b []byte) ([]byte, error) {

	// create data
	var exercise exercise

	err := json.Unmarshal(b, &exercise)
	if err != nil {
		// handle error case
	}

	db.Save(&exercise)

	// Return a success message (maybe edit later to return the exercise?)
	return []byte("{\"message\": \"Exercise successfully added\"}"), nil
}
