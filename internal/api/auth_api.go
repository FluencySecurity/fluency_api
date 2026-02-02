package api

import (
	"fmt"

	fluencyAPI "github.com/SecurityDo/fluency_api/api"
	"github.com/SecurityDo/fluency_api/model"
	fluencyModel "github.com/SecurityDo/fluency_api/model"
)

func (c *Client) AddUser(name, displayName, role, org string) error {

	// Use structured logging
	//c.Logger.Info("adding user",
	//	"name", name,
	//	"role", role,
	//)

	authService := fluencyAPI.NewAuthService(c.fluencyClient)

	err := authService.AddUser(&fluencyAPI.AddUserRequest{
		User: &fluencyModel.UserEntry{
			Username:     name,
			Email:        name,
			FirstName:    displayName,
			Roles:        []string{role},
			Organization: org,
		},
	})
	if err != nil {
		c.Logger.Error("failed to add user", "error", err, "name", name, "role", role)
		return fmt.Errorf("failed to add user: %w", err)
	}
	return nil
}

func (c *Client) DeleteUser(username string) (err error) {

	// Use structured logging

	authService := fluencyAPI.NewAuthService(c.fluencyClient)

	err = authService.DeleteUser(username)
	if err != nil {
		c.Logger.Error("failed to delete user", "error", err, "name", username)
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (c *Client) ListUser() (users []*model.UserEntry, err error) {

	// Use structured logging

	authService := fluencyAPI.NewAuthService(c.fluencyClient)

	users, err = authService.ListUser()
	if err != nil {
		c.Logger.Error("failed to list user", "error", err)
		return nil, fmt.Errorf("failed to list user: %w", err)
	}
	return users, nil
}

func (c *Client) AddToken(name, displayName, role string) (token string, err error) {

	// Use structured logging
	//c.Logger.Info("adding user",
	//	"name", name,
	//	"role", role,
	//)

	authService := fluencyAPI.NewAuthService(c.fluencyClient)

	token, err = authService.AddToken(name, displayName, role)
	if err != nil {
		c.Logger.Error("failed to add token", "error", err, "name", name, "role", role)
		return "", fmt.Errorf("failed to add token: %w", err)
	}
	return token, nil
}

func (c *Client) DeleteToken(name string) (err error) {

	// Use structured logging

	authService := fluencyAPI.NewAuthService(c.fluencyClient)

	err = authService.DeleteToken(name)
	if err != nil {
		c.Logger.Error("failed to delete token", "error", err, "name", name)
		return fmt.Errorf("failed to delete token: %w", err)
	}
	return nil
}

func (c *Client) ListToken() (tokens []*model.ApiTokenEntry, err error) {

	// Use structured logging

	authService := fluencyAPI.NewAuthService(c.fluencyClient)

	tokens, err = authService.ListToken()
	if err != nil {
		c.Logger.Error("failed to list token", "error", err)
		return nil, fmt.Errorf("failed to list token: %w", err)
	}
	return tokens, nil
}
