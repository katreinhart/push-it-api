package model

import (
	"encoding/json"
	"errors"
)

// CreateWorkout instantiates a new workout object with data from request body
func CreateWorkout(uid string, b []byte) ([]byte, error) {
	var workout workoutModel

	err := json.Unmarshal(b, &workout)

	if err != nil {
		return nil, err
	}
	workout.User = uid

	db.Save(&workout)

	return []byte("{\"message\": \"Workout created successfully.\"}"), nil
}

// AddExerciseToWorkout takes a wid which is existing workout and adds a new exercise to it.
func AddExerciseToWorkout(wid string, b []byte) ([]byte, error) {
	var exercisePosted workoutExerciseAsPosted
	var exercise workoutExercise

	err := json.Unmarshal(b, &exercisePosted)

	if err != nil {
		return nil, err
	}

	exerciseID, err := getExerciseID(exercisePosted.ExerciseName)
	if err != nil {
		return nil, err
	}

	exercise.ExerciseID = exerciseID
	exercise.WorkoutID = wid
	exercise.GoalSets = exercisePosted.GoalSets
	exercise.GoalRepsPerSet = exercisePosted.GoalRepsPerSet

	db.Save(&exercise)

	return []byte("{\"message\": \"Exercise added successfully.\"}"), nil
}

// GetWorkout returns the given workout based on workout ID (wid)
func GetWorkout(wid string) ([]byte, error) {
	var workout workoutModel

	db.First(&workout, "id = ?", wid)

	if workout.ID == 0 {
		return nil, errors.New("Not found")
	}

	return json.Marshal(workout)
}
