package scripts

import (
	"fmt"

	"github.com/Jibaru/env0/pkg/auth"
)

// LoadAndValidateToken loads and validates the authentication token
func LoadAndValidateToken() (string, error) {
	token, err := auth.LoadToken()
	if err != nil {
		return "", fmt.Errorf("authentication required")
	}
	if token == "" {
		return "", fmt.Errorf("authentication token expired")
	}
	return token, nil
}
