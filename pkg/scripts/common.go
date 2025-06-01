package scripts

import (
	"fmt"

	"github.com/Jibaru/env0/pkg/client"
)

// LoadAndValidateToken loads and validates the authentication token
func LoadAndValidateToken() (string, error) {
	token, err := client.LoadToken()
	if err != nil {
		return "", fmt.Errorf("authentication required")
	}
	if token == "" {
		return "", fmt.Errorf("authentication token expired")
	}
	return token, nil
}
