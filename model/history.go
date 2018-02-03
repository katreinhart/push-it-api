package model

import (
	"encoding/json"
	"errors"
	"strconv"
)

// History retrieves all of the given user's completed workouts and returns nested objects
func History(uid string) ([]CompletedWorkout, error) {

	var workouts []WorkoutModel
	var _workouts []CompletedWorkout

	var exercises []WorkoutExercise
	var sets []WorkoutExerciseSet

	// find all workouts in db with user uid
	db.Find(&workouts, "user_id = ?", uid).Where("completed = ?", true)
	if len(workouts) == 0 {
		return nil, ErrorNotFound
	}

	// find all associated exercises and sets associated with each workout
	for _, wko := range workouts {
		var _exercises []TransformedWorkoutExercise
		var _sets []WorkoutSet

		var strID = strconv.Itoa(int(wko.ID))
		db.Find(&exercises, "workout_id = ?", wko.ID)
		for _, ex := range exercises {
			var exerciseName, _ = getExerciseName(ex.ExerciseID)
			_exercises = append(_exercises, TransformedWorkoutExercise{WorkoutID: strID, ExerciseID: ex.ID, ExerciseName: exerciseName, GoalSets: ex.GoalRepsPerSet, GoalRepsPerSet: ex.GoalRepsPerSet})
		}
		db.Find(&sets, "workout_id = ?", strID)
		for _, set := range sets {
			_sets = append(_sets, WorkoutSet{Exercise: set.ExerciseName, Weight: set.Weight, RepsAttempted: set.RepsAttempted, RepsCompleted: set.RepsCompleted})
		}
		workoutID := strconv.Itoa(int(wko.ID))
		_workouts = append(_workouts, CompletedWorkout{User: uid, WorkoutID: workoutID, Start: wko.Start, End: wko.End, Rating: wko.Rating, Comments: wko.Comments, Exercises: _exercises, Sets: _sets})
	}

	return _workouts, nil
}

func FetchSavedExercises(uid string) ([]byte, error) {
	var workouts []WorkoutModel
	var _workouts []savedWorkout
	var exercises []WorkoutExercise

	// find all workouts in db with user uid
	db.Find(&workouts, "user_id = ?", uid).Where("completed = ?", false)
	if len(workouts) == 0 {
		return nil, errors.New("Not found")
	}

	// find all associated exercises and sets associated with each workout
	for _, wko := range workouts {
		var _exercises []TransformedWorkoutExercise

		var strID = strconv.Itoa(int(wko.ID))
		db.Find(&exercises, "workout_id = ?", wko.ID)
		for _, ex := range exercises {
			var exerciseName, _ = getExerciseName(ex.ExerciseID)
			_exercises = append(_exercises, TransformedWorkoutExercise{WorkoutID: strID, ExerciseID: ex.ID, ExerciseName: exerciseName, GoalSets: ex.GoalSets, GoalRepsPerSet: ex.GoalRepsPerSet})
		}
		workoutID := strconv.Itoa(int(wko.ID))
		_workouts = append(_workouts, savedWorkout{User: uid, WorkoutID: workoutID, Exercises: _exercises})
	}

	js, err := json.Marshal(_workouts)

	return js, err
}
