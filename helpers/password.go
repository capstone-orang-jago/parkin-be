package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

// HassPass hashes the given plaintext password using bcrypt with a cost of 12.
// Returns the hashed password as a string, or an error if hashing fails.
func HassPass(pass string) (string, error) {
	pw := []byte(pass)
	salt := 12
	Hass, err := bcrypt.GenerateFromPassword(pw, salt)
	if err != nil {
		return "", err
	}
	return string(Hass), nil
}

// ComparePass compares a bcrypt hashed password with its possible plaintext equivalent.
// Returns true if the password matches, or false and an error if it does not.
func ComparePass(hashedPassword, plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return false, err
	}
	return true, nil
}
