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
	workout.UserID = uid

	db.Save(&workout)

	js, err := json.Marshal(&workout)

	return js, err
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
		var _sets []transformedWorkoutSet

		db.Find(&dbExercises, "workout_id = ?", wid)

		for _, ex := range dbExercises {
			db.Find(&dbSets, "workout_exercise_id = ?", ex.ID)
			for _, set := range dbSets {
				name, _ := getExerciseName(ex.ID)
				_sets = append(_sets, transformedWorkoutSet{ExerciseName: name, Weight: set.Weight, RepsAttempted: set.RepsAttempted, RepsCompleted: set.RepsCompleted})
			}
			exName, err := getExerciseName(ex.ExerciseID)
			if err != nil {
				panic("exercise not found")
			}
			_exercises = append(_exercises, transformedWorkoutExercise{WorkoutID: ex.WorkoutID, ExerciseID: ex.ExerciseID, ExerciseName: exName, GoalSets: ex.GoalSets, GoalRepsPerSet: ex.GoalRepsPerSet})
		}

		completed.Sets = _sets
		completed.Exercises = _exercises

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

// MarkWorkoutAsCompleted updates workout given supplied body b
func MarkWorkoutAsCompleted(uid string, id string, b []byte) ([]byte, error) {
	var workout workoutModel
	var workoutUpdate updateWorkoutModel

	db.First(&workout, "id = ?", id)
	if workout.UserID != uid {
		return nil, errors.New("Forbidden")
	}

	err := json.Unmarshal(b, &workoutUpdate)

	if err != nil {
		return nil, err
	}

	db.Model(&workout).Update("completed", workoutUpdate.Completed)
	db.Model(&workout).Update("rating", workoutUpdate.Rating)
	db.Model(&workout).Update("comments", workoutUpdate.Comments)

	js, err := json.Marshal(workout)

	return js, err
}
