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

// FetchSingleExercise finds exercise matching id in database and returns JSON
func FetchSingleExercise(id string) ([]byte, error) {
	var exercise exercise

	db.First(&exercise, id)

	if exercise.ID == 0 {
		return nil, errors.New("Not found")
	}

	_exercise := transformedExercise{ID: exercise.ID, Name: exercise.Name, Link: exercise.Link}

	js, err := json.Marshal(_exercise)

	return js, err
}

// CreateExercise adds an exercise to the database.
func CreateExercise(b []byte) ([]byte, error) {

	// create data
	var exercise exercise

	err := json.Unmarshal(b, &exercise)
	if err != nil {
		return nil, errors.New("Something went wrong")
	}

	db.Save(&exercise)

	// Return a success message (maybe edit later to return the exercise?)
	return []byte("{\"message\": \"Exercise successfully added\"}"), nil
}
