package mTool

import "regexp"

//[1-30]@[1-10].[1-4]
func VerifyEmail(email string) bool {
	regular := "^[0-9a-zA-Z]{1,30}@[0-9a-zA-Z]{1,10}.[a-z]{1,4}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(email)
}

// 1[10]
func VerifyPhoneNumber(phoneNumber string) bool {
	regular := "^1[0-9]{10}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(phoneNumber)
}
