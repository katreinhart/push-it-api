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
	for _, i := range exs {
		_exs = append(_exs, TransformedExercise{ID: i.ID, Name: i.Name, Link: i.Link})
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
	var ex Exercise

	db.Find(&ex, "name = ?", exerciseName)
	if ex.ID == 0 {
		return 0, ErrorNotFound
	}

	return ex.ID, nil
}

func getExerciseName(eid uint) (string, error) {
	var ex Exercise
	db.Find(&ex, "id = ?", eid)

	if ex.ID == 0 {
		return "", ErrorNotFound
	}

	return ex.Name, nil
}
