package model

import (
	"strconv"
)

// History retrieves all of the given user's completed workouts and returns nested objects
func History(uid string) ([]CompletedWorkout, error) {

	var workouts []WorkoutModel
	var _workouts []CompletedWorkout

	var exercises []WorkoutExercise
	var sets []WorkoutExerciseSet

	// find all workouts in db with user uid
	db.Where("user_id = ? AND completed = ?", uid, true).Find(&workouts)
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
			_exercises = append(_exercises, TransformedWorkoutExercise{WorkoutID: strID, ExerciseID: ex.ID, ExerciseName: exerciseName, GoalSets: ex.GoalSets, GoalRepsPerSet: ex.GoalRepsPerSet})
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

// FetchSavedWorkouts returns all workouts that have not been completed
func FetchSavedWorkouts(uid string) ([]SavedWorkout, error) {
	var workouts []WorkoutModel
	var _workouts []SavedWorkout
	var exercises []WorkoutExercise

	// find all workouts in db with user uid
	db.Where("user_id = ? AND completed = ?", uid, false).Find(&workouts)
	if len(workouts) == 0 {
		return nil, ErrorNotFound
	}

	// find all associated exercises and sets associated with each workout
	for _, wko := range workouts {
		var _exercises []TransformedWorkoutExercise

		var strID = strconv.Itoa(int(wko.ID))
		db.Find(&exercises, "workout_id = ?", wko.ID)
		for _, ex := range exercises {
			var exerciseName, _ = getExerciseName(ex.ExerciseID)
			_exercises = append(_exercises, TransformedWorkoutExercise{WorkoutID: strID, ExerciseID: ex.ID, ExerciseName: exerciseName, GoalWeight: ex.GoalWeight, GoalSets: ex.GoalSets, GoalRepsPerSet: ex.GoalRepsPerSet})
		}
		workoutID := strconv.Itoa(int(wko.ID))
		_workouts = append(_workouts, SavedWorkout{UserID: uid, CreatedAt: wko.CreatedAt, WorkoutID: workoutID, Exercises: _exercises})
	}

	return _workouts, nil
}
