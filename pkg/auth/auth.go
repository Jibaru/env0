package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Auth represents the authentication data
type Auth struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// User represents the authenticated user data
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Save persists the authentication data to the user's home directory
func Save(auth Auth) error {
	home, err := getHomeDir()
	if err != nil {
		return err
	}

	cfgPath := filepath.Join(home, ".env0_cfg")
	if err := os.MkdirAll(cfgPath, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	data, err := json.Marshal(auth)
	if err != nil {
		return fmt.Errorf("failed to marshal auth data: %v", err)
	}

	if err := os.WriteFile(filepath.Join(cfgPath, "auth.json"), data, 0600); err != nil {
		return fmt.Errorf("failed to write auth file: %v", err)
	}

	return nil
}

// Load reads the saved authentication data
func Load() (*Auth, error) {
	home, err := getHomeDir()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath.Join(home, ".env0_cfg", "auth.json"))
	if err != nil {
		return nil, fmt.Errorf("no auth data found: %v", err)
	}

	var auth Auth
	if err := json.Unmarshal(data, &auth); err != nil {
		return nil, fmt.Errorf("invalid auth data: %v", err)
	}

	if auth.Token == "" {
		return nil, fmt.Errorf("no valid token found")
	}

	return &auth, nil
}

// LoadToken is a convenience function that returns just the token
// Useful for backward compatibility and simple token checks
func LoadToken() (string, error) {
	auth, err := Load()
	if err != nil {
		return "", err
	}
	return auth.Token, nil
}

// getHomeDir returns the user's home directory in a cross-platform way
func getHomeDir() (string, error) {
	// Try USERPROFILE for Windows first
	if home := os.Getenv("USERPROFILE"); home != "" {
		return home, nil
	}
	// Try HOME for Unix-like systems
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}
	return "", fmt.Errorf("unable to determine user home directory")
}
