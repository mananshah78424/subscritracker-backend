package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Common validation functions that can be used across the application

// IsValidEmail validates email format using regex
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidPassword validates password strength
func IsValidPassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	// Check for at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}

	// Check for at least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}

	// Check for at least one number
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("password must contain at least one number")
	}

	return nil
}

// IsValidName validates name format
func IsValidName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if len(strings.TrimSpace(name)) < 2 {
		return fmt.Errorf("name must be at least 2 characters long")
	}

	// Check for valid characters (letters, spaces, hyphens, apostrophes)
	if !regexp.MustCompile(`^[a-zA-Z\s\-']+$`).MatchString(name) {
		return fmt.Errorf("name contains invalid characters")
	}

	return nil
}

// IsValidString checks if a string is not empty after trimming
func IsValidString(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

// IsValidLength checks if a string meets length requirements
func IsValidLength(value, fieldName string, min, max int) error {
	length := len(strings.TrimSpace(value))
	if length < min {
		return fmt.Errorf("%s must be at least %d characters long", fieldName, min)
	}
	if max > 0 && length > max {
		return fmt.Errorf("%s must be no more than %d characters long", fieldName, max)
	}
	return nil
}
