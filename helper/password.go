package helper

import "golang.org/x/crypto/bcrypt"

func HashPassword(plaintextPassword string) (string, error) {
	// Generate a new salt with a cost of 12
	salt, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return "", err
	}

	// Convert the salt to a string and return it
	return string(salt), nil
}
