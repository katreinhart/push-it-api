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
	var completed completedWorkout

	db.First(&workout, "id = ?", wid)

	if workout.ID == 0 {
		return nil, errors.New("Not found")
	}

	if workout.Completed {
		var dbExercises []workoutExercise
		var dbSets []workoutExerciseSet

		var _exercises []transformedWorkoutExercise
		var _sets []workoutSetAsPosted

		db.Find(&dbExercises, "workout_id = ?", wid)

		for _, ex := range dbExercises {
			db.Find(&dbSets, "workout_exercise_id = ?", ex.ID)
			for _, set := range dbSets {
				_sets = append(_sets, workoutSetAsPosted{RepsAttempted: set.RepsAttempted, RepsCompleted: set.RepsCompleted})
			}
			exName, err := getExerciseName(ex.ExerciseID)
			if err != nil {
				panic("exercise not found")
			}
			_exercises = append(_exercises, transformedWorkoutExercise{WorkoutID: ex.WorkoutID, ExerciseID: ex.ExerciseID, ExerciseName: exName, GoalSets: ex.GoalSets, GoalRepsPerSet: ex.GoalRepsPerSet})
		}

		js, err := json.Marshal(completed)

		return js, err
	}

	return json.Marshal(workout)
}

// AddExerciseSet adds a set of the given exercise to the workout in question.
func AddExerciseSet(wid string, eid string, b []byte) ([]byte, error) {

	var newSet workoutSetAsPosted
	var newExSet workoutExerciseSet

	err := json.Unmarshal(b, &newSet)

	if err != nil {
		return nil, err
	}

	newExSet.WorkoutExerciseID = eid
	newExSet.Weight = newSet.Weight
	newExSet.RepsAttempted = newSet.RepsAttempted
	newExSet.RepsCompleted = newSet.RepsCompleted

	db.Save(&newExSet)

	return []byte("{\"message\": \"Set added successfully.\"}"), nil
}
