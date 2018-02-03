package model

import (
	"strconv"
)

// CreateWorkout instantiates a new workout object with data from request body
func CreateWorkout(wk WorkoutModel) WorkoutModel {

	db.Save(&wk)

	return wk
}

// GetWorkout returns the given workout based on workout ID (wid)
func GetWorkout(wid string) (WorkoutModel, error) {
	var wk WorkoutModel

	db.First(&wk, "id = ?", wid)

	if wk.ID == 0 {
		return WorkoutModel{}, ErrorNotFound
	}

	return wk, nil
}

// GetCompletedWorkout returns the workout with nested exercises and sets.
func GetCompletedWorkout(wk WorkoutModel) (CompletedWorkout, error) {

	var completed CompletedWorkout

	var dbExercises []WorkoutExercise
	var dbSets []WorkoutExerciseSet

	var _exercises []TransformedWorkoutExercise
	var _sets []WorkoutSet

	db.Find(&dbExercises, "workout_id = ?", wk.ID)

	for _, ex := range dbExercises {
		db.Find(&dbSets, "workout_exercise_id = ?", ex.ID)
		for _, set := range dbSets {
			name, _ := getExerciseName(ex.ID)
			_sets = append(_sets, WorkoutSet{Exercise: name, Weight: set.Weight, RepsAttempted: set.RepsAttempted, RepsCompleted: set.RepsCompleted})
		}
		exName, err := getExerciseName(ex.ExerciseID)
		if err != nil {
			panic("exercise not found")
		}
		_exercises = append(_exercises, TransformedWorkoutExercise{WorkoutID: ex.WorkoutID, ExerciseID: ex.ExerciseID, ExerciseName: exName, GoalSets: ex.GoalSets, GoalRepsPerSet: ex.GoalRepsPerSet})
	}

	workoutID := strconv.Itoa(int(wk.ID))
	completed.WorkoutID = workoutID
	completed.Sets = _sets
	completed.Exercises = _exercises

	return completed, nil
}

// AddExerciseToWorkout takes a wid which is existing workout and adds a new exercise to it.
func AddExerciseToWorkout(wid string, ep WorkoutExerciseAsPosted) (WorkoutExercise, error) {

	var ex WorkoutExercise

	exerciseID, err := getExerciseID(ep.ExerciseName)

	if err != nil {
		return WorkoutExercise{}, err
	}

	ex.ExerciseID = exerciseID
	ex.GoalWeight = ep.GoalWeight
	ex.WorkoutID = wid
	ex.GoalSets = ep.GoalSets
	ex.GoalRepsPerSet = ep.GoalRepsPerSet

	db.Save(&ex)

	return ex, nil
}

// AddExerciseSet adds a set of the given exercise to the workout in question.
func AddExerciseSet(wid string, wsp WorkoutSet) (WorkoutExerciseSet, error) {
	var newExSet WorkoutExerciseSet

	newExSet.WorkoutID = wid
	newExSet.ExerciseName = wsp.Exercise
	newExSet.Weight = wsp.Weight
	newExSet.RepsAttempted = wsp.RepsAttempted
	newExSet.RepsCompleted = wsp.RepsCompleted

	db.Save(&newExSet)

	return newExSet, nil
}

// MarkWorkoutAsCompleted updates workout given supplied body b
func MarkWorkoutAsCompleted(uid string, id string, uw UpdateWorkout) (WorkoutModel, error) {
	var workout WorkoutModel

	db.First(&workout, "id = ?", id)
	if workout.UserID != uid {
		return WorkoutModel{}, ErrorForbidden
	}

	db.Model(&workout).Update("completed", uw.Completed)
	db.Model(&workout).Update("rating", uw.Rating)
	db.Model(&workout).Update("comments", uw.Comments)

	return workout, nil
}

// UpdateWorkoutTimestamps updates database entry for given workout with started and completed timestamps
func UpdateWorkoutTimestamps(id string, ts PostedTimestamps) (WorkoutModel, error) {
	var workout WorkoutModel

	db.First(&workout, "id = ?", id)
	if workout.ID == 0 {
		return WorkoutModel{}, ErrorNotFound
	}

	db.Model(&workout).Update("start", ts.StartedAt)
	db.Model(&workout).Update("end", ts.FinishedAt)

	return workout, nil
}
