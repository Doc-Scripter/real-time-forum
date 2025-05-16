package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	// "errors"
)

const (
	saltLength = 16
)

// GenerateSalt generates a random salt
func GenerateSalt() (string, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// HashPassword hashes a password with a salt
func HashPassword(password, salt string) (string, error) {
	// Combine password and salt
	combined := []byte(password + salt)
	
	// Hash the combined string
	hash := sha256.Sum256(combined)
	
	// Return the base64 encoded hash
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

// VerifyPassword verifies a password against a stored hash
func VerifyPassword(password, storedHash, salt string) bool {
	// Hash the provided password with the stored salt
	hashedPassword, err := HashPassword(password, salt)
	if err != nil {
		return false
	}
	
	// Compare the hashes
	return hashedPassword == storedHash
}

// GenerateSaltedHash generates a salted hash for a password
func GenerateSaltedHash(password string) (string, string, error) {
	// Generate a salt
	salt, err := GenerateSalt()
	if err != nil {
		return "", "", err
	}
	
	// Hash the password with the salt
	hash, err := HashPassword(password, salt)
	if err != nil {
		return "", "", err
	}
	
	return hash, salt, nil
}
