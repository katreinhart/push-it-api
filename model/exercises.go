package model

// ListExercises returns all available exercises
func ListExercises() ([]TransformedExercise, error) {

	var exs []Exercise
	var _exs []TransformedExercise

	db.Find(&exs)

	if len(exs) == 0 {
		return nil, ErrorNotFound
	}

	// transform exercises
	for _, item := range exs {
		_exs = append(_exs, TransformedExercise{ID: item.ID, Name: item.Name, Link: item.Link})
	}

	return _exs, nil
}

// FetchSingleExercise finds exercise matching id in database and returns JSON
func FetchSingleExercise(id string) (TransformedExercise, error) {
	var ex Exercise

	db.First(&ex, id)

	if ex.ID == 0 {
		return TransformedExercise{}, ErrorNotFound
	}

	_ex := TransformedExercise{ID: ex.ID, Name: ex.Name, Link: ex.Link}

	return _ex, nil
}

// CreateExercise adds an exercise to the database.
func CreateExercise(ex Exercise) TransformedExercise {

	db.Save(&ex)

	// Return a success message (maybe edit later to return the exercise?)
	return TransformedExercise{ID: ex.ID, Name: ex.Name, Link: ex.Link}
}

func getExerciseID(exerciseName string) (uint, error) {
	var dbExercise Exercise

	db.Find(&dbExercise, "name = ?", exerciseName)
	if dbExercise.ID == 0 {
		return 0, ErrorNotFound
	}

	return dbExercise.ID, nil
}

func getExerciseName(eid uint) (string, error) {
	var exercise Exercise
	db.Find(&exercise, "id = ?", eid)

	if exercise.ID == 0 {
		return "", ErrorNotFound
	}

	return exercise.Name, nil
}
