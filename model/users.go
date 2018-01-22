package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user instance in the database
func CreateUser(b []byte) ([]byte, error) {

	// declare data structures to be used
	var user, dbUser userModel

	// Unmarshal the json body into the user data structure
	err := json.Unmarshal(b, &user)

	// handle error if any
	if err != nil {
		return []byte("{\"message\": \"Something went wrong.\"}"), err
	}

	// See if the user exists in the database. If so, return an error (no duplicates allowed)
	db.First(&dbUser, "email = ?", user.Email)
	if dbUser.ID != 0 {
		return []byte("{\"message\": \"User already exists in DB.\"}"), errors.New("User already exists")
	}

	// Hash the user's password using bcrypt helper function and handle any error
	hash, err := hashPassword(user.Password)
	if err != nil {
		return []byte("{\"message\": \"Sorry, something went wrong.\"}"), err
	}

	// Overwrite the user's password with the hashed version (no plaintext storage of passwords)
	user.Password = hash

	// save the user in the DB
	db.Save(&user)

	// Get the user back from the database so you have the correct ID
	db.First(&dbUser, "email = ?", user.Email)

	// create and sign the JWT
	t, err := createAndSignJWT(dbUser)

	// Handle error in JWT creation/signing
	if err != nil {
		fmt.Println(err.Error())
		return []byte("{\"message\": \"Sorry, something went wrong.\"}"), err
	}

	// create transformed version of user structure, marshal it into JSON and return
	_user := transformedUser{ID: user.ID, Email: user.Email, Goal: user.Goal, Level: user.Level, Token: t}
	js, err := json.Marshal(_user)
	return js, err
}

// LoginUser takes info from controller and returns a token if user is who they claim to be
func LoginUser(b []byte) ([]byte, error) {

	// Declare data types and unmarshal JSON into user struct
	var user, dbUser userModel
	err := json.Unmarshal(b, &user)

	if err != nil {
		return []byte("{\"message\": \"Something went wrong.\"}"), err
	}

	db.First(&dbUser, "email = ?", user.Email)

	// handle user not found error
	if dbUser.ID == 0 {
		return []byte("{\"message\": \"User not found in DB.\"}"), errors.New("Not found")
	}

	// See if password matches the hashed password from the database
	match := checkPasswordHash(user.Password, dbUser.Password)
	if !match {
		return []byte("{\"message\": \"Check your inputs and try again.\"}"), errors.New("Unauthorized")
	}
	// Create and sign JWT; handle any error
	t, err := createAndSignJWT(dbUser)
	if err != nil {
		return []byte("{\"message\": \"Something went wrong with JWT.\"}"), err
	}
	// create transmission friendly user struct
	_user := transformedUser{ID: dbUser.ID, Email: dbUser.Email, Name: dbUser.Name, Goal: dbUser.Goal, Level: dbUser.Level, Token: t}

	// marshal user into JSON and return
	js, err := json.Marshal(_user)
	return js, err
}

// SetUserInfo takes info from the onboarding screen and updates the database.
func SetUserInfo(b []byte) ([]byte, error) {

	// Declare data types and unmarshal JSON into user struct
	var user, dbUser userModel
	err := json.Unmarshal(b, &user)

	if err != nil {
		return []byte("{\"message\": \"Something went wrong.\"}"), err
	}

	db.First(&dbUser, "email = ?", user.Email)

	// handle user not found error
	if dbUser.ID == 0 {
		return []byte("{\"message\": \"User not found in DB.\"}"), errors.New("Not found")
	}
	// Here I should get JWT and make sure user matches. This should be done in the controller.

	// update first name, level, and goal in the database
	dbUser.Level = user.Level
	dbUser.Goal = user.Goal
	dbUser.Name = user.Name

	db.Save(&dbUser)

	// create transmission friendly user struct
	_user := transformedUser{ID: dbUser.ID, Email: dbUser.Email, Name: user.Name, Level: user.Level, Goal: user.Goal}

	// marshal user into JSON and return
	js, err := json.Marshal(_user)
	return js, err
}

// --------------------- Helper Functions ---------------------
// user login password helper functions
// from https://gowebexamples.com/password-hashing/
func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(b), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// JWT helper function
func createAndSignJWT(user userModel) (string, error) {

	// create the expiration time, build claim, create and sign token, and return.
	e := time.Now().Add(time.Hour * 24).Unix()
	c := CustomClaims{
		user.ID,
		jwt.StandardClaims{
			ExpiresAt: e,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	secret := []byte(os.Getenv("SECRET"))
	t, err := token.SignedString(secret)
	return t, err
}
