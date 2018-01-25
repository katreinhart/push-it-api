package model

import (
	"encoding/json"
	"errors"
	"strconv"
)

// History retrieves all of the given user's completed workouts and returns nested objects
func History(uid string) ([]byte, error) {
	var workouts []workoutModel
	var _workouts []completedWorkout

	var exercises []workoutExercise
	var sets []workoutExerciseSet

	// find all workouts in db with user uid
	db.Find(&workouts, "user_id = ?", uid)
	if len(workouts) == 0 {
		return nil, errors.New("Not found")
	}

	// find all associated exercises and sets associated with each workout
	for _, wko := range workouts {
		var _exercises []transformedWorkoutExercise
		var _sets []transformedWorkoutSet

		db.Find(&exercises, "workout_id = ?", wko.ID)
		for _, ex := range exercises {
			var strID = strconv.Itoa(int(wko.ID))
			var exerciseName, _ = getExerciseName(ex.ExerciseID)

			db.Find(&sets, "workout_exercise_id = ?", strID)
			for _, set := range sets {
				_sets = append(_sets, transformedWorkoutSet{ExerciseName: exerciseName, Weight: set.Weight, RepsAttempted: set.RepsAttempted, RepsCompleted: set.RepsCompleted})
			}
			_exercises = append(_exercises, transformedWorkoutExercise{WorkoutID: strID, ExerciseID: ex.ID, ExerciseName: exerciseName, GoalSets: ex.GoalRepsPerSet, GoalRepsPerSet: ex.GoalRepsPerSet})
		}
		_workouts = append(_workouts, completedWorkout{User: uid, Start: wko.Start, End: wko.End, Rating: wko.Rating, Comments: wko.Comments, Exercises: _exercises, Sets: _sets})
	}

	js, err := json.Marshal(_workouts)

	return js, err
}
