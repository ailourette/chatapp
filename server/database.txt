package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserName  string
	Password  string
	FirstName string
	LastName  string
}

// Insert new user account information into database
func insertRecord(db *sql.DB, userName, passWord, firstName, lastName string) {
	_, err := db.Exec("INSERT INTO MYSTOREDB.Users VALUES (?,?,?,?)", userName, passWord, firstName, lastName)
	if err != nil {
		Error.Println("Error inserting record into database")
		fmt.Println(err)
	} else {
		fmt.Println("New user account information added to database successfully")
	}
}

// Delete user account information from database
func deleteRecord(db *sql.DB, userName string) int {
	results, err := db.Exec("DELETE FROM MYSTOREDB.Users where Username=?", userName)
	if err != nil {
		Error.Println("Error deleting record from database")
		fmt.Println(err)
		return 0
	} else {
		rows, _ := results.RowsAffected()
		return int(rows)
	}
}

// Update user account password in database
func changePasswordRecord(db *sql.DB, userName string, password []byte) {
	results, err := db.Query("SELECT * FROM MYSTOREDB.Users where Username=?", userName)
	if err != nil {
		fmt.Println(err)
	} else {
		for results.Next() {
			var person User
			err = results.Scan(&person.UserName, &person.Password, &person.FirstName, &person.LastName)
			if err != nil {
				fmt.Println(err)
			} else {
				_, err := db.Exec("UPDATE MYSTOREDB.Users set password = ?, firstname = ?, lastname = ? where Username=?", password, &person.FirstName, &person.LastName, &person.UserName)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

// get the hashed password of the user in string type
func getPasswordOfUser(db *sql.DB, userName string) string {
	results, err := db.Query("SELECT * FROM MYSTOREDB.Users where Username=?", userName)
	if err != nil {
		fmt.Println(err)
		return ""
	} else {
		for results.Next() {
			var person User
			err = results.Scan(&person.UserName, &person.Password, &person.FirstName, &person.LastName)
			if err != nil {
				fmt.Println(err)
				return ""
			} else {
				return person.Password
			}
		}
	}
	return ""
}

// get the first name of user in string type
func getFirstNameOfUser(db *sql.DB, userName string) string {
	results, err := db.Query("SELECT * FROM MYSTOREDB.Users where Username=?", userName)
	if err != nil {
		fmt.Println(err)
		return ""
	} else {
		for results.Next() {
			var person User
			err = results.Scan(&person.UserName, &person.Password, &person.FirstName, &person.LastName)
			if err != nil {
				fmt.Println(err)
				return ""
			} else {
				return person.FirstName
			}
		}
	}
	return ""
}

// hash the given password using bcrypt()
func hashPassword(password string) []byte {
	if hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost); err != nil {
		fmt.Println(err)
		return nil
	} else {
		return hash
	}
}

//                   saved in the db        user supplied
func verifyPassword(hashedPassword []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil
}

// Open data base to insert new user account information
func userSignupDataBase(un, pw, fn, ln string) {
	db, err := sql.Open("mysql", "user1:password@tcp(127.0.0.1:3306)/MYSTOREDB")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened.")
	}
	defer db.Close()

	//----adding new user---
	userName := un
	password := pw
	firstName := fn
	lastName := ln
	insertRecord(db, userName, string(hashPassword(password)), firstName, lastName)
}

// Open data base to remove user account information
func userDeleteDataBase(userName string) int {
	db, err := sql.Open("mysql", "user1:password@tcp(127.0.0.1:3306)/MYSTOREDB")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened.")
	}
	defer db.Close()

	//----deleting user---
	return deleteRecord(db, userName)
}

// Open data base to update user account password
func userChangePasswordDataBase(userName string, password string) {
	hashedPassword := hashPassword(password)
	db, err := sql.Open("mysql", "user1:password@tcp(127.0.0.1:3306)/MYSTOREDB")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened.")
	}
	defer db.Close()

	//----change password of user---
	changePasswordRecord(db, userName, hashedPassword)
}

// Checks if the entered username and password match an account in database. If yes login successful, else unsuccessful
func authenticatingUserFromDataBase(un string, pw string) bool {
	db, err := sql.Open("mysql", "user1:password@tcp(127.0.0.1:3306)/MYSTOREDB")
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened.")
	}
	defer db.Close()

	//---authenticating user---
	userName := un
	password := pw

	// retrieve the user's saved password (in string); hashed
	userSavedPassword := getPasswordOfUser(db, userName)

	// the password saved in the db the user's supplied password

	if verifyPassword([]byte(userSavedPassword), password) {
		return true
	} else {
		return false
	}
}
