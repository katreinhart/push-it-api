package model

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"

	// importing postgres dialect for GORM
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// declare DB
var db *gorm.DB

type (
	userModel struct {
		gorm.Model
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Goal     string `json:"goal"`
		Level    string `json:"level"`
	}

	transformedUser struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Goal  string `json:"goal"`
		Level string `json:"level"`
		Token string `json:"token"`
	}

	workoutModel struct {
		gorm.Model
		UserID    string    `json:"user_id"`
		Start     time.Time `json:"start_time"`
		End       time.Time `json:"finish_time"`
		Rating    int       `json:"rating"`
		Comments  string    `json:"comments"`
		Completed bool      `json:"completed"`
	}

	updateWorkoutModel struct {
		ID        string `json:"id"`
		Completed bool   `json:"completed"`
		Rating    int    `json:"rating"`
		Comments  string `json:"comments"`
	}

	completedWorkout struct {
		User      string                       `json:"uid"`
		Start     time.Time                    `json:"start_time"`
		End       time.Time                    `json:"finish_time"`
		Rating    int                          `json:"rating"`
		Comments  string                       `json:"comments"`
		Exercises []transformedWorkoutExercise `json:"exercises"`
		Sets      []transformedWorkoutSet      `json:"sets"`
	}

	exercise struct {
		gorm.Model
		Name string `json:"ex_name"`
		Link string `json:"info_url"`
	}

	transformedExercise struct {
		ID   uint   `json:"id"`
		Name string `json:"ex_name"`
		Link string `json:"info_url"`
	}

	workoutExercise struct {
		gorm.Model
		WorkoutID      string `json:"workout_id"`
		ExerciseID     uint   `json:"exercise_id"`
		GoalSets       int    `json:"goal_sets"`
		GoalRepsPerSet int    `json:"goal_reps_per_set"`
	}

	transformedWorkoutExercise struct {
		WorkoutID      string `json:"workout_id"`
		ExerciseID     uint   `json:"exercise_id"`
		ExerciseName   string `json:"exercise_name"`
		GoalSets       int    `json:"goal_sets"`
		GoalRepsPerSet int    `json:"goal_reps_per_set"`
	}

	workoutExerciseAsPosted struct {
		ExerciseName   string `json:"exercise_name"`
		GoalSets       int    `json:"goal_sets"`
		GoalRepsPerSet int    `json:"goal_reps"`
	}

	workoutExerciseSet struct {
		gorm.Model
		WorkoutExerciseID string `json:"workout_exercise_id"`
		Weight            int    `json:"weight"`
		RepsAttempted     int    `json:"reps_att"`
		RepsCompleted     int    `json:"reps_comp"`
	}

	workoutSetAsPosted struct {
		Weight        int `json:"weight"`
		RepsAttempted int `json:"reps_att"`
		RepsCompleted int `json:"reps_comp"`
	}

	transformedWorkoutSet struct {
		ExerciseName  string `json:"exercise"`
		Weight        int    `json:"weight"`
		RepsAttempted int    `json:"reps_att"`
		RepsCompleted int    `json:"reps_comp"`
	}

	postedGoals struct {
		Goal1 postedBasicGoal `json:"goal1"`
		Goal2 postedBasicGoal `json:"goal2"`
	}

	postedBasicGoal struct {
		UserID     string `json:"uid"`
		GoalDate   string `json:"goal_date"`
		Exercise   string `json:"exercise"`
		GoalWeight string `json:"goal_weight"`
	}

	secondaryGoal struct {
		gorm.Model
		UserID     string    `json:"uid"`
		GoalDate   time.Time `json:"goal_date"`
		Exercise   string    `json:"exercise"`
		GoalWeight int       `json:"goal_weight"`
	}

	transformedGoal struct {
		GoalID     int       `json:"goal_id"`
		UserID     string    `json:"uid"`
		SetDate    time.Time `json:"set_date"`
		GoalDate   time.Time `json:"goal_date"`
		Exercise   string    `json:"exercise"`
		GoalWeight int       `json:"goal_weight"`
	}

	goalsResponse struct {
		Data []transformedGoal `json:"data"`
	}

	updatePrimaryGoal struct {
		Goal string `json:"goal"`
	}

	WeightPlate struct {
		Plate45 int `json:"45"`
		Plate35 int `json:"35"`
		Plate25 int `json:"25"`
		Plate10 int `json:"10"`
		Plate05 int `json:"05"`
		Plate02 int `json:"025"`
	}

	// CustomClaims for JWT handling
	CustomClaims struct {
		UID uint `json:"uid"`
		jwt.StandardClaims
	}
)

// init function runs at setup; connects to database
func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading env file")
	}

	hostname := os.Getenv("HOST")
	dbname := os.Getenv("DBNAME")
	uname := os.Getenv("DBUSER")
	password := os.Getenv("PASSWORD")

	dbString := "host=" + hostname + " user=" + uname + " dbname=" + dbname + " sslmode=disable password=" + password

	// var err error
	db, err = gorm.Open("postgres", dbString)
	if err != nil {
		fmt.Println(err.Error())
		panic("Unable to connect to DB")
	}

	db.AutoMigrate(&userModel{})
	db.AutoMigrate(&exercise{})
	db.AutoMigrate(&secondaryGoal{})
	db.AutoMigrate(&workoutModel{})
	db.AutoMigrate(&workoutExercise{})
	db.AutoMigrate(&workoutExerciseSet{})
}
