package model

import (
	"encoding/json"
	"errors"
)

// ListExercises returns all available exercises
func ListExercises() ([]byte, error) {

	var exercises []exercise
	var _exercises []transformedExercise

	db.Find(&exercises)

	if len(exercises) == 0 {
		return nil, errors.New("Not found")
	}

	// transform exercises
	for _, item := range exercises {
		_exercises = append(_exercises, transformedExercise{ID: item.ID, Name: item.Name, Link: item.Link})
	}

	// marshal into json
	js, err := json.Marshal(_exercises)

	// return js, err
	return js, err
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
