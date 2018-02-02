package model

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser creates a new user instance in the database
func CreateUser(u UserModel) (TransformedUser, error) {

	// declare data structures to be used
	var dbUser UserModel

	// See if the user exists in the database. If so, return an error (no duplicates allowed)
	db.First(&dbUser, "email = ?", u.Email)
	if dbUser.ID != 0 {
		return TransformedUser{}, ErrorUserExists
	}

	// Hash the user's password using bcrypt helper function and handle any error
	hash, err := hashPassword(u.Password)
	if err != nil {
		return TransformedUser{}, ErrorInternalServer
	}

	// Overwrite the user's password with the hashed version (no plaintext storage of passwords)
	u.Password = hash

	// save the user in the DB
	db.Save(&u)

	// Get the user back from the database so you have the correct ID
	db.First(&dbUser, "email = ?", u.Email)

	// create and sign the JWT
	t, err := createAndSignJWT(dbUser)

	// Handle error in JWT creation/signing
	if err != nil {
		return TransformedUser{}, ErrorInternalServer
	}

	// create transformed version of user structure, marshal it into JSON and return
	_user := TransformedUser{ID: dbUser.ID, Email: u.Email, Goal: u.Goal, Level: u.Level, Token: t}
	return _user, nil
}

// LoginUser takes info from controller and returns a token if user is who they claim to be
func LoginUser(u UserModel) (TransformedUser, error) {

	// Declare data types and unmarshal JSON into user struct
	var dbUser UserModel

	db.First(&dbUser, "email = ?", u.Email)

	// handle user not found error
	if dbUser.ID == 0 {
		return TransformedUser{}, ErrorNotFound
	}

	// See if password matches the hashed password from the database
	match := checkPasswordHash(u.Password, dbUser.Password)
	if !match {
		return TransformedUser{}, ErrorForbidden
	}
	// Create and sign JWT; handle any error
	t, err := createAndSignJWT(dbUser)
	if err != nil {
		return TransformedUser{}, ErrorInternalServer
	}
	// create transmission friendly user struct
	_user := TransformedUser{ID: dbUser.ID, Email: dbUser.Email, Name: dbUser.Name, Goal: dbUser.Goal, Level: dbUser.Level, Token: t}

	return _user, nil
}

// SetUserInfo takes info from the onboarding screen and updates the database.
func SetUserInfo(uid uint, u UserModel) (TransformedUser, error) {

	// Declare data types and unmarshal JSON into user struct
	var dbUser UserModel

	db.First(&dbUser, "email = ?", u.Email)

	// handle user not found error
	if dbUser.ID == 0 {
		fmt.Println("user not found")
		return TransformedUser{}, ErrorNotFound
	}
	// Compare to uid from JWT and make sure user matches.
	if dbUser.ID != uid {
		fmt.Println("wrong user")
		return TransformedUser{}, ErrorForbidden
	}

	// update first name, level, and goal in the database
	dbUser.Level = u.Level
	dbUser.Goal = u.Goal
	dbUser.Name = u.Name

	db.Save(&dbUser)

	// create transmission friendly user struct
	_user := TransformedUser{ID: dbUser.ID, Email: dbUser.Email, Name: u.Name, Level: u.Level, Goal: u.Goal}

	return _user, nil
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
func createAndSignJWT(user UserModel) (string, error) {

	// create the expiration time, build claim, create and sign token, and return.
	// token expires in 30 days (720 hours)
	e := time.Now().Add(time.Hour * 720).Unix()
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
