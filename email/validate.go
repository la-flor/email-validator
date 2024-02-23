package email

import (
	"errors"
	"log"
	"net"
)

// Since not all servers allow address verification for
// security reasons, we can only test if an address is
// invalid and cannot guarantee that the address is valid.
//
// a true response means it is invalid
// a false response means it has not been proven to be invalid
func CheckIfInvalid(email string) (bool, string) {
	parsedEmail, err := Parse(email)

	if err != nil {
		log.Println(err)
		return true, "Invalid email input"
	}

	validMX, _, _ := checkMXValidity(parsedEmail.Domain)
	
	if !validMX {
		return true, "Invalid MX records"
	}

	return false, "Email has not been proved invalid"
}

func checkMXValidity(host string) (bool, string, uint16) {
	mxrecords, err := net.LookupMX(host)

	if err != nil || len(mxrecords) < 1 {
		log.Println(errors.New("unable to get MX records"))
		return false, "", 0
	}

	Host := mxrecords[0].Host
	Pref := mxrecords[0].Pref

	return true, Host, Pref
}
