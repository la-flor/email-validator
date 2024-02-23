package email

import (
	"strings"
	"errors"
)

type ParsedEmail struct {
	Address	string
	Domain	string
}

func Parse(email string) (ParsedEmail, error) {
	idx := strings.LastIndexByte(email, '@')
	var parsedEmail ParsedEmail

	if (idx < 0) {
		return parsedEmail, errors.New("email does not make address and domain distinction")
	}

	parsedEmail.Address = email[:idx]
	parsedEmail.Domain = email[idx+1:]

	return parsedEmail, nil
}
