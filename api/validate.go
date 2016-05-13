package api

import (
	"fmt"
	"regexp"
	"strings"
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

func ValidateFileName(fileName string) bool {
	valid, err := regexp.Compile("^([\\w]+\\.[a-z]+)$")
	if err != nil {
		fmt.Println(err)
		return false
	}
	if valid.MatchString(fileName) {
		return true
	} else {
		return false
	}
}

//Splits a string at every "/" and trims leading and trailing "/"
func ParseURL(path string) []string {
	var parsedURL []string
	parsedURL = append(parsedURL, path)
	fmt.Println(parsedURL)

	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	if strings.HasSuffix(path, "/") {
		path = path[:(len(path) - 1)]
	}
	parsedURL = append(parsedURL, strings.Split(path, "/")...)
	fmt.Println(parsedURL)
	return parsedURL
}

//Returns the part of an urlSubStr that should contain a token
func GetToken(urlSubStr []string) string {
	return urlSubStr[2]
}

//Retunrs the part of an urlSubStr that should contain a file name
func GetFilename(urlSubStr []string) string {
	return urlSubStr[3]
}
