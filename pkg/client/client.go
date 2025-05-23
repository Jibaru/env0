package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

const baseURL = "https://env0-api.vercel.app"

// ClientError wraps HTTP status codes and underlying errors
type ClientError struct {
	Status int
	Err    error
}

func (e *ClientError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("status %d: %v", e.Status, e.Err)
	}
	return fmt.Sprintf("status %d", e.Status)
}

// Client defines the API methods
type Client interface {
	Signup(ctx context.Context, username, email, password string) error
	Login(ctx context.Context, usernameOrEmail, password string) error
	CreateApp(ctx context.Context, name string) (ownerName string, err error)
	GetApp(ctx context.Context, fullAppName string) (envs map[string]map[string]interface{}, err error)
	UpdateApp(ctx context.Context, fullAppName string, envs map[string]map[string]interface{}) error
	AddUser(ctx context.Context, fullAppName, username string) error
	RemoveUser(ctx context.Context, fullAppName, username string) error
}

// client is the concrete implementation
type client struct {
	token string
}

// New returns a new API client. Pass empty token for unauthenticated calls.
func New(token string) Client {
	return &client{token: token}
}

// internal doRequest helper
func (c *client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, []byte, error) {
	url := baseURL + path
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		reader = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Add("Authorization", c.token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	return resp, data, nil
}

// Signup registers a new user
func (c *client) Signup(ctx context.Context, username, email, password string) error {
	body := map[string]string{"username": username, "email": email, "password": password}
	resp, data, err := c.doRequest(ctx, http.MethodPost, "/api/v1/register", body)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusCreated {
		return nil
	}
	// try parse error message
	var res map[string]interface{}
	_ = json.Unmarshal(data, &res)
	if msg, ok := res["error"].(string); ok {
		return &ClientError{Status: resp.StatusCode, Err: errors.New(msg)}
	}
	return &ClientError{Status: resp.StatusCode}
}

// Login authenticates and saves token to config
func (c *client) Login(ctx context.Context, usernameOrEmail, password string) error {
	body := map[string]string{"emailOrUsername": usernameOrEmail, "password": password}
	resp, data, err := c.doRequest(ctx, http.MethodPost, "/api/v1/login", body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		var res map[string]interface{}
		_ = json.Unmarshal(data, &res)
		if msg, ok := res["error"].(string); ok {
			return &ClientError{Status: resp.StatusCode, Err: errors.New(msg)}
		}
		return &ClientError{Status: resp.StatusCode}
	}
	type LoginResp struct {
		Token string `json:"token"`
		User  struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"user"`
	}
	var res LoginResp
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	// save auth
	if err := SaveAuth(res.Token); err != nil {
		return err
	}
	c.token = res.Token
	return nil
}

// CreateApp creates a new app, returns ownerName
func (c *client) CreateApp(ctx context.Context, name string) (string, error) {
	body := map[string]string{"name": name}
	resp, data, err := c.doRequest(ctx, http.MethodPost, "/api/v1/apps", body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusCreated {
		var res map[string]interface{}
		_ = json.Unmarshal(data, &res)
		if msg, ok := res["error"].(string); ok {
			return "", &ClientError{Status: resp.StatusCode, Err: errors.New(msg)}
		}
		return "", &ClientError{Status: resp.StatusCode}
	}
	var res map[string]interface{}
	_ = json.Unmarshal(data, &res)
	owner, _ := res["ownerName"].(string)
	return owner, nil
}

// GetApp retrieves app environments
func (c *client) GetApp(ctx context.Context, fullAppName string) (map[string]map[string]interface{}, error) {
	path := "/api/v1/apps/" + url.PathEscape(fullAppName)
	resp, data, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		var res map[string]interface{}
		_ = json.Unmarshal(data, &res)
		if msg, ok := res["error"].(string); ok {
			return nil, &ClientError{Status: resp.StatusCode, Err: errors.New(msg)}
		}
		return nil, &ClientError{Status: resp.StatusCode}
	}
	var res struct {
		Envs map[string]map[string]interface{} `json:"envs"`
	}
	_ = json.Unmarshal(data, &res)
	return res.Envs, nil
}

// UpdateApp pushes environment changes
func (c *client) UpdateApp(ctx context.Context, fullAppName string, envs map[string]map[string]interface{}) error {
	path := "/api/v1/apps/" + url.PathEscape(fullAppName)
	body := map[string]interface{}{"envs": envs}
	resp, data, err := c.doRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		var res map[string]interface{}
		_ = json.Unmarshal(data, &res)
		if msg, ok := res["error"].(string); ok {
			return &ClientError{Status: resp.StatusCode, Err: errors.New(msg)}
		}
		return &ClientError{Status: resp.StatusCode}
	}
	return nil
}

// AddUser adds a user to the app
func (c *client) AddUser(ctx context.Context, fullAppName, username string) error {
	path := "/api/v1/apps/" + url.PathEscape(fullAppName) + "/users/" + url.PathEscape(username)
	resp, data, err := c.doRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		var res map[string]interface{}
		_ = json.Unmarshal(data, &res)
		if msg, ok := res["error"].(string); ok {
			return &ClientError{Status: resp.StatusCode, Err: errors.New(msg)}
		}
		return &ClientError{Status: resp.StatusCode}
	}
	return nil
}

// RemoveUser removes a user from the app
func (c *client) RemoveUser(ctx context.Context, fullAppName, username string) error {
	path := "/api/v1/apps/" + url.PathEscape(fullAppName) + "/users/" + url.PathEscape(username)
	resp, data, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		var res map[string]interface{}
		_ = json.Unmarshal(data, &res)
		if msg, ok := res["error"].(string); ok {
			return &ClientError{Status: resp.StatusCode, Err: errors.New(msg)}
		}
		return &ClientError{Status: resp.StatusCode}
	}
	return nil
}

// SaveAuth persists the token to the user's home directory
func SaveAuth(token string) error {
	cfgPath := filepath.Join(os.Getenv("HOME"), ".env0_cfg")
	if err := os.MkdirAll(cfgPath, 0700); err != nil {
		return err
	}
	auth := map[string]string{"token": token}
	data, _ := json.Marshal(auth)
	return os.WriteFile(filepath.Join(cfgPath, "auth.json"), data, 0600)
}

// LoadToken reads the saved token
func LoadToken() (string, error) {
	data, err := os.ReadFile(filepath.Join(os.Getenv("HOME"), ".env0_cfg", "auth.json"))
	if err != nil {
		return "", err
	}
	var auth map[string]string
	if err := json.Unmarshal(data, &auth); err != nil {
		return "", err
	}
	return auth["token"], nil
}
