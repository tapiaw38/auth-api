package utils

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

const (
	HASH_COST = 8
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandomString returns a random string of the given length
func RandomString(lenght int) string {
	b := make([]rune, lenght)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// ValidateEmail validates an email
func ValidateEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(regex).MatchString(email)
}

// ValidatePassword validates a password
func ValidatePassword(password string) bool {
	regex := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`
	return regexp.MustCompile(regex).MatchString(password)
}

// ConvertInterfaceToString converts an interface to a string
func ConvertInterfaceToString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case float64:
		return fmt.Sprintf("%f", v), nil
	default:
		return "", fmt.Errorf("cannot convert %v to string", v)
	}
}

// Generate a token
func GenerateToken() (string, error) {
	// Create a random byte slice.
	tokenBytes := make([]byte, 16)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	// Encode the byte slice to a hexadecimal string.
	token := hex.EncodeToString(tokenBytes)

	return token, nil
}

// HashPassword hashes a password
func HashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), HASH_COST)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
