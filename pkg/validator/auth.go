package validator

import (
	"fmt"
	"strings"
)

// SignUpRequest represents the signup request structure
type SignUpRequest struct {
	Email      string `json:"email" form:"email" validate:"required,email"`
	Password   string `json:"password" form:"password" validate:"required,min=8"`
	Name       string `json:"name" form:"name" validate:"required"`
	GivenName  string `json:"given_name" form:"given_name"`
	FamilyName string `json:"family_name" form:"family_name"`
}

// LoginRequest represents the login request structure
type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required"`
}

// ValidateSignUp validates signup request fields
func ValidateSignUp(req SignUpRequest) error {
	var errors []string

	// Validate email
	if err := IsValidString(req.Email, "email"); err != nil {
		errors = append(errors, err.Error())
	} else if !IsValidEmail(req.Email) {
		errors = append(errors, "email must be a valid email address")
	}

	// Validate password
	if err := IsValidString(req.Password, "password"); err != nil {
		errors = append(errors, err.Error())
	} else if err := IsValidPassword(req.Password); err != nil {
		errors = append(errors, err.Error())
	}

	// Validate name
	if err := IsValidString(req.Name, "name"); err != nil {
		errors = append(errors, err.Error())
	} else if err := IsValidName(req.Name); err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf("signup failed: %s", strings.Join(errors, "; "))
	}

	return nil
}

// ValidateLogin validates login request fields
func ValidateLogin(req LoginRequest) error {
	var errors []string

	// Validate email
	if err := IsValidString(req.Email, "email"); err != nil {
		errors = append(errors, err.Error())
	} else if !IsValidEmail(req.Email) {
		errors = append(errors, "email must be a valid email address")
	}

	// Validate password
	if err := IsValidString(req.Password, "password"); err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf("login failed: %s", strings.Join(errors, "; "))
	}

	return nil
}
