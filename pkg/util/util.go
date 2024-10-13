package util

import (
	"errors"
	"net"
	"regexp"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func BcryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareBcryptPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func WrapErr(err error, msg string) error {
	return errors.New(err.Error() + ": " + msg)
}

func ValidationErrs(errs []error) error {
	msg := "validation errors: "

	if len(errs) == 0 {
		return errors.New("no validation errors, but ValidationErrs() was called")
	}

	for _, err := range errs {
		msg += err.Error() + ", "
	}

	msg = msg[:len(msg)-2]

	return errors.New(msg)
}

// IsValidHostname checks if the input is a valid domain name or IP address.
func IsValidHostname(hostname string) bool {
	// Check if it's a valid IP address
	if net.ParseIP(hostname) != nil {
		return true
	}

	// Regex for valid domain name (including subdomains)
	// Domain name should be between 1 and 253 characters and can include letters, numbers, and hyphens.
	// It must not start or end with a hyphen, and each label should be between 1 and 63 characters.
	return regexp.MustCompile(
		`^([a-zA-Z0-9]{1}[a-zA-Z0-9-]{0,61}[a-zA-Z0-9]{1}\.)+[a-zA-Z]{2,}$`,
	).
		MatchString(hostname)
}

// IsValidPort checks if the input is a valid port number.
func IsValidPort(port string) bool {
	i, err := strconv.Atoi(port)
	if err != nil {
		return false
	}

	return i > 0 && i <= 65535
}

func IsValidEmail(email string) bool {
	return regexp.MustCompile(
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
	).
		MatchString(email)
}
