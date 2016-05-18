package api

import (

		//"fmt"
		"regexp"
		//"strings"
		"testing"

)

// Tests to see if the length of the token is 32

func TestLenToken(t *testing.T){

	var token1, token2, token3 string 

	token1 = "d1g5dvc05kfjgt03s0v0hlfe42svt52s" 
	token2 = "0000000000000000000000000000000" 
	token3 = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"	

	if len(token1) != 32 {

		t.Error("Token 1's size is not 32")

	}

	if len(token2) == 32 {

		t.Error("Token 2's size is not 32")

	}

	if ((len(token3) > 32) || (len(token3) < 32)) {

		t.Error("Token 3's size is not 32")

	}

}

// Tests to see if the token contains valid characters

func TestValCharToken(t *testing.T){

	var token11, token12, token13 string

	token11 = "s3a2edfg35f6fb43ft4f6fw3rfdr3fg5"
	token12 = "__-___--____-___///////___--"
	token13 = "as____---dsd__23232---///sdsd"

	if check1, _ := regexp.MatchString("^[a-zA-Z0-9]*$", token11); !check1 {

		t.Error("Token has invalid characters")

	}

	if check2, _ := regexp.MatchString("^[a-zA-Z0-9]*$", token12); check2 {

		t.Error("Token has invalid characters")

	}

	if check3, _ := regexp.MatchString("^[a-zA-Z0-9]*$", token13); check3 {

		t.Error("Token has invalid characters")

	}

}