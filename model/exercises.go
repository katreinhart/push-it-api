package model

// ListExercises returns all available exercises
func ListExercises() []byte {

	var exercises []exercise
	// var _exercises []transformedExercise

	db.Find(&exercises)

	if len(exercises) == 0 {
		// return not found error
	}

	// transform exercises

	// marshal into json

	// return js, err
	return []byte("")
}
