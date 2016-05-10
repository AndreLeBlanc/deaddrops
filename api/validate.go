package api

import (
	"regexp"
)

//Validates that a token string is correct length and contains correct characters for a md5 hash.
func ValidateToken(token string) bool {
	if len(token) != 32 {
		return false
	} else if match, _ := regexp.MatchString("^[a-zA-Z0-9]*$", token); !match {
		return match
	} else {
		return true
	}

}

//Validates that the file is allowed to be stored in a deadrop.
func ValidateFile( /*file*/ ) bool {
	//TODO: file validation, ex. not too big
	return true
}
