package email

import (
	"errors"
	"log"
	"net"
	"net/smtp"
	"strconv"
	"time"
)



var mxToServerSMTPAddress = map[string]string{
	"google": "smtp.gmail.com",
	"smtpin.zoho.com": "smtp.zoho.com",
	"mx.zoho.com": "smtp.zoho.com",
	"mta6.am0.yahoodns.net": "smtp.mail.yahoo.com",
	"o2": "smtp.o2.ie",
	"mx-aol.mail.gm0.yahoodns.net": "smtp.aol.com",
	"mx0b-00191d01.pphosted.com": "smtp.att.yahoo.com",
	"hotmail-com.olc.protection.outlook.com": "smtp.live.com",
	"mxb-00143702.gslb.pphosted.com": "smtp.comcast.net",
	"mxa-0024a201.gslb.pphosted.com": "outgoing.verizon.net",
}

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

	validMX, Host, _ := checkMXValidity(parsedEmail.Domain)
	
	if !validMX {
		return true, "Invalid MX records"
	}

	invalidHost := checkHostIfInvalid(Host)

	if (invalidHost) {
		return true, "Invalid host"
	}

	invalidEmail, _ := invalidFullEmail(Host, email)

	if (invalidEmail) {
		return true, "Email is invalid"
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

func checkHostIfInvalid(Host string) bool {
	client, err := dialTimeout(Host)

	if err != nil {
		return true
	}

	defer client.Close()

	return false
}

func dialTimeout(Host string) (*smtp.Client, error) {
	address := net.JoinHostPort(Host,strconv.Itoa(25))

	timeout := time.Second * 5

	conn, err := net.DialTimeout("tcp", address, timeout)

	if err != nil {
		println("dialTimeout failed", conn, err)
		return nil, err
	}

	defer conn.Close()

	return smtp.NewClient(conn, Host)
}



func invalidFullEmail(Host string, email string) (bool, string) {
	client, err := dialTimeout(Host)
	if err != nil {
		return true, "dialTimeout failed"
	}

	err = client.Hello(Host)
	if err != nil {
		return true, "Hello ping failed"
	}

	serverMailAddress := mxToServerSMTPAddress[Host]

	if serverMailAddress == "" {
		return false, "Unable to identify server address"
	}

	err = client.Mail(serverMailAddress)
	if err != nil {
		return true, "Server mail address failure"
	}

	err = client.Rcpt(email)
	if err != nil {
		return true, "RCPT failure"
	}

	return false, "No failure detected"
}
