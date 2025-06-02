package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

// IsValid checks if the auth token is a valid JWT and not expired
func (a *Auth) IsValid() bool {
	if a == nil || a.Token == "" {
		return false
	}

	// Check if it's a JWT token (should have 3 parts separated by dots)
	parts := strings.Split(a.Token, ".")
	if len(parts) != 3 {
		return false
	}

	// Parse the token without verifying the signature
	// We only want to check the structure and expiration
	token, _, err := new(jwt.Parser).ParseUnverified(a.Token, jwt.MapClaims{})
	if err != nil {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	// Check expiration if present
	if exp, ok := claims["exp"].(float64); ok {
		expTime := time.Unix(int64(exp), 0)
		if time.Now().After(expTime) {
			return false
		}
	}

	return true
}

// IsAuthenticated returns true if the auth data exists and has a valid token
func (a *Auth) IsAuthenticated() bool {
	return a != nil && a.IsValid()
}

// HasUserInfo returns true if the auth data contains user information
func (a *Auth) HasUserInfo() bool {
	return a != nil && a.User.Username != "" && a.User.Email != ""
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

	// Check if token is valid
	if !auth.IsValid() {
		return nil, fmt.Errorf("token is invalid or expired")
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
