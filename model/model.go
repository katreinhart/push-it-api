package model

import (
	"fmt"
	"os"

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
		Goal     string `json:"primary_goal"`
	}

	transformedUser struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Goal  string `json:"goal"`
		Token string `json:"token"`
	}

	// workout struct {
	// 	gorm.Model
	// 	User  uint      `json:"uid"`
	// 	Start time.Time `json:"start_time"`
	// 	End   time.Time `json:"finish_time"`
	// }

	// exercise struct {
	// 	gorm.Model
	// 	Name string `json:"ex_name"`
	// 	Link string `json:"info_url"`
	// }

	// workoutExercise struct {
	// 	gorm.Model
	// 	WorkoutID  uint `json:"workout_id"`
	// 	ExerciseID uint `json:"exercise_id"`
	// }

	// workoutExerciseSet struct {
	// 	gorm.Model
	// 	WorkoutExerciseID uint `json:"workout_exercise_id"`
	// 	Weight            uint `json:"weight"`
	// 	RepsAttempted     uint `json:"reps_att"`
	// 	RepsCompleted     uint `json:"reps_comp"`
	// }

	// secondaryGoal struct {
	// 	gorm.Model
	// 	SetDate    time.Time `json:"set_date"`
	// 	GoalDate   time.Time `json:"goal_date"`
	// 	Exercise   uint      `json:"exercise_id"`
	// 	GoalWeight uint      `json:"goal_weight"`
	// }

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
}
