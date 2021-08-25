package main

import (
	"regexp"
)

// Regular expression to check if email is of valid format
var emailRegex = regexp.MustCompile(
	"^[\\w!#$%&'*+/=?`{|}~^-]+(?:\\.[\\w!#$%&'*+/=?`{|}~^-]+)*@(?:[a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,6}$") // regular expression

// Checks if email is of valid format, return true if valid else false
func isEmailValid(e string) bool {
	if len(e) < 3 || len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
