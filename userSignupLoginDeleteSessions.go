package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	uuid "github.com/satori/go.uuid"
)

var (
	mutex         sync.Mutex // Concurrency
	tpl           *template.Template
	userFirstName string
)

// Declaration of type patient
type user struct {
	UserName  string `validate:"required"`
	Password  string `validate:"required"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Language  string `validate:"required"`
}

// Declaration of dbPatients and dbSessions map
var (
	dbUsers    = map[string]user{}
	dbSessions = map[string]string{}
)

/*
Go In Action 2 Done
This function allows patient to sign up for a account
Error handling is implemented to takes care of invalid user input
Mutex lock and unlock is used to ensure no two users modify the patient account map at the same time
*/
func userSignup(res http.ResponseWriter, req *http.Request) {
	// Get cookie
	c, err := req.Cookie("session")
	if err != nil {
		sID := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(res, c)
	}

	// If patient exists already, get user
	var p user
	if un, ok := dbSessions[c.Value]; ok {
		p = dbUsers[un]
	}

	// Process form submission
	if req.Method == http.MethodPost { // Error handling for no input
		un := req.FormValue("username")
		pw := req.FormValue("password")
		firstName := req.FormValue("firstname")
		lastName := req.FormValue("lastname")
		language := req.FormValue("language")

		// Trim space of username Input
		un = strings.TrimSpace(un)

		mutex.Lock()
		p = user{un, pw, firstName, lastName, language}
		validate := validator.New()
		err := validate.Struct(p)
		if err != nil {
			io.WriteString(res, `
			<html>
			<meta http-equiv='refresh' content='5; url=/userSignup '/>
			Please fill in all fields!<br>
			You will be redirected shortly in 5 seconds...<br>
			</html>
			`)
			return
		}
		dbSessions[c.Value] = un
		dbUsers[un] = p

		// Checks if user enter correct email address format
		result := isEmailValid(un)
		if !result {
			io.WriteString(res, `
			<html>
			<meta http-equiv='refresh' content='5; url=/ '/>
			Incorrect Email address format!<br>
			You will be redirected shortly in 5 seconds...<br>
			</html>
		`)
			return
		}
		userSignupDataBase(un, pw, firstName, lastName)
		sessionWriteCsv()
		mutex.Unlock()
		io.WriteString(res, `
 			<html>
 			<meta http-equiv='refresh' content='5; url=/userLogin '/>
			 You have successfully signed up! <br>
			 You will be redirected shortly in 5 seconds...<br>
 			</html>
 		`)
		return

	}
	tpl.ExecuteTemplate(res, "userSignup.gohtml", p)
}

/*
Go In Action 2 Done
This function allows user to login to their account
Error handling is implemented to takes care of invalid user input
*/
func userLogin(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		pw := req.FormValue("password")

		// Trim space of username Input
		un = strings.TrimSpace(un)

		// Checks if user enter correct email address format
		result := isEmailValid(un)
		if !result {
			io.WriteString(res, `
			<html>
			<meta http-equiv='refresh' content='5; url=/ '/>
			Incorrect Email address format!<br>
			You will be redirected shortly in 5 seconds...<br>
			</html>
		`)
			return
		}

		result = authenticatingUserFromDataBase(un, pw)

		// Check if username exist in database
		if !result {
			io.WriteString(res, `
 			<html>
 			<meta http-equiv='refresh' content='5; url=/ '/>
			 Incorrect Username or Password <br>
			 You will be redirected shortly in 5 seconds...<br>
 			</html>
 		`)
		}

		// Save user information in global variable for use in making/edit or cancel appointment
		db, err := sql.Open("mysql", "user1:password@tcp(127.0.0.1:3306)/MYSTOREDB")
		if err != nil {
			panic(err.Error())
		} else {
			fmt.Println("Database opened.")
		}
		defer db.Close()

		userFirstName = getFirstNameOfUser(db, un)

		// Create a session
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(res, c)
		dbSessions[c.Value] = un
		http.Redirect(res, req, "/userLoginSuccess", http.StatusSeeOther) // Redirect to user feature page
		return
	}
	tpl.ExecuteTemplate(res, "userLogin.gohtml", nil)
}

/*
This function allows admin to delete all sessions
No error handling is implemented because we do not expect user input
*/
// func deleteSessions(res http.ResponseWriter, req *http.Request) {
// 	for key := range dbSessions {
// 		delete(dbSessions, key)
// 	}
// 	sessionWriteCsv()
// 	io.WriteString(res, `
// 		<html>
// 		<meta http-equiv='refresh' content='5; url=/adminLoginSuccess '/>
// 		Sessions deleted<br>
// 		You will be redirected shortly in 5 seconds...<br>
// 		</html>
// 		`)
// }

/*
Go In Action 2 Done
This function allows admin to delete a patient's account
Error handling is implemented to takes care of invalid user input
Mutex lock and unlock is used to ensure no two users modify the patient account map at the same time
*/
// func deleteUsers(res http.ResponseWriter, req *http.Request) {
// 	// Get user name input, if name match delete user else user not found
// 	tpl.ExecuteTemplate(res, "deleteUsers.gohtml", nil)

// 	if req.Method == http.MethodPost {
// 		patientInputUserName := req.FormValue("username")
// 		validate := validator.New() // The patient signup func works without issues
// 		err := validate.Var(patientInputUserName, "required")

// 		if err != nil {
// 			http.Error(res, "Please enter patient name!", http.StatusForbidden)
// 			return
// 		} else {
// 			patientExist := patientDeleteDataBase(patientInputUserName)
// 			if patientExist != 0 {
// 				io.WriteString(res, `
// 					<html>
// 					<meta http-equiv='refresh' content='5; url=/adminLoginSuccess '/>
// 					User account deleted<br>
// 					You will be redirected shortly in 5 seconds...<br>
// 					</html>
// 					`)
// 			} else {
// 				io.WriteString(res, `
// 				<html>
// 				<meta http-equiv='refresh' content='5; url=/adminLoginSuccess '/>
// 			  	User not found!<br>
// 				You will be redirected shortly in 5 seconds...<br>
// 				</html>
// 			`)
// 			}
// 		}
// 	}
// }

// Get username, new password. Then search for user using username. Hash password then write to database.
func userChangePassword(res http.ResponseWriter, req *http.Request) {
	//var db *sql.DB
	if req.Method == http.MethodPost { // Error handling for no input
		un := req.FormValue("username")
		pw := req.FormValue("password")

		// Trim space of username Input
		un = strings.TrimSpace(un)

		if un == "" && pw == "" {
			http.Error(res, "Please enter your username(Email) and new password!", http.StatusForbidden)
			return
		} else {
			userChangePasswordDataBase(un, pw)
			sessionWriteCsv()
			io.WriteString(res, `
 			<html>
			 <meta http-equiv='refresh' content='5; url=/userLogin '/>
			 Your password is updated successfully! <br>
			 You will be redirected shortly in 5 seconds...<br>
 			</html>
 		`)
			return
		}
	}
	tpl.ExecuteTemplate(res, "userChangePassword.gohtml", nil)
}

// This function checks if user is already logged in
func alreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("session")
	if err != nil {
		return false
	}
	username := dbSessions[myCookie.Value]
	_, ok := dbUsers[username]
	return ok
}

/*
This function allows user to logout of their account
Sessions and cookie will be deleted
*/
func logout(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("session")

	// Delete the session
	delete(dbSessions, myCookie.Value)

	// Remove the cookie
	myCookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)
	io.WriteString(res, `
		<html>
		<meta http-equiv='refresh' content='5; url=/logout '/>
		You have successfully logged out.<br>
		You will be redirected shortly in 5 seconds...
		</html>
		`)
}
