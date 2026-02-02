package api

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SecurityDo/fluency_api/client"
	"github.com/SecurityDo/fluency_api/model"
)

type AuthService struct {
	client *client.FluencyClient
}

func NewAuthService(client *client.FluencyClient) *AuthService {
	return &AuthService{
		client: client,
	}
}

type AddUserRequest struct {
	User *model.UserEntry `json:"user"`
}

func (s *AuthService) AddUser(req *AddUserRequest) error {
	_, err := s.client.GenericCall("api/auth", "userAdd", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding user: %v\n", err.Error())
		return err
	}
	fmt.Fprintln(os.Stderr, "User added successfully")
	return nil
}

type ListUserResponse struct {
	Users []*model.UserEntry `json:"users"`
}

func (s *AuthService) ListUser() (users []*model.UserEntry, err error) {
	res, err := s.client.GenericCall("api/auth", "userList", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing users: %v\n", err.Error())
		return nil, err
	}
	var result ListUserResponse
	err = json.Unmarshal(res.GetBytes(), &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing user list response: %v\n", err.Error())
		return nil, err
	}
	return result.Users, nil
}

type GetUserRequest struct {
	Username string `json:"username"`
}

type GetUserResponse struct {
	User *model.UserEntry `json:"user"`
}

func (s *AuthService) GetUser(username string) (*model.UserEntry, error) {
	req := &GetUserRequest{
		Username: username,
	}

	res, err := s.client.GenericCall("api/auth", "getUser", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting user %s: %v\n", username, err.Error())
		return nil, err
	}
	var result GetUserResponse
	err = json.Unmarshal(res.GetBytes(), &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing get user response: %v\n", err.Error())
		return nil, err
	}
	return result.User, nil
}

type DeleteUserRequest struct {
	Username string `json:"username"`
}

func (s *AuthService) DeleteUser(username string) error {
	req := &DeleteUserRequest{
		Username: username,
	}

	_, err := s.client.GenericCall("api/auth", "userDelete", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting user %s: %v\n", username, err.Error())
		return err
	}
	fmt.Println("User deleted successfully")
	return nil
}

type tokenRequestArgs struct {
	Name  string               `json:"id,omitempty"`
	Entry *model.ApiTokenEntry `json:"entry,omitempty"`
	Flag  bool                 `json:"flag"`
}
type tokenRequest struct {
	Action string            `json:"action"`
	Args   *tokenRequestArgs `json:"args,omitempty"`
}

type AddTokenResponse struct {
	Token string `json:"token"`
}

func (s *AuthService) AddToken(name, description, role string) (token string, err error) {
	req := &tokenRequest{

		Action: "add",
		Args: &tokenRequestArgs{

			Entry: &model.ApiTokenEntry{
				Name:        name,
				Description: description,
				Roles:       []string{role},
			},
		},
	}

	res, err := s.client.GenericCall("api/auth", "api_token", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding token %s: %v\n", name, err.Error())
		return "", err
	}
	var result AddTokenResponse
	err = json.Unmarshal(res.GetBytes(), &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing add token response: %v\n", err.Error())
		return "", err
	}
	return result.Token, nil
}

func (s *AuthService) DeleteToken(name string) (err error) {
	req := &tokenRequest{
		Action: "delete",
		Args: &tokenRequestArgs{
			Entry: &model.ApiTokenEntry{
				Name: name,
			},
		},
	}
	_, err = s.client.GenericCall("api/auth", "api_token", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting token %s: %v\n", name, err.Error())
		return err
	}
	return nil
}

type ListTokenResponse struct {
	Entries []*model.ApiTokenEntry `json:"entries"`
}

func (s *AuthService) ListToken() (tokens []*model.ApiTokenEntry, err error) {
	req := &tokenRequest{
		Action: "list",
	}
	res, err := s.client.GenericCall("api/auth", "api_token", req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing token: %v\n", err.Error())
		return nil, err
	}
	var result ListTokenResponse
	err = json.Unmarshal(res.GetBytes(), &result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing list token response: %v\n", err.Error())
		return nil, err
	}
	return result.Entries, nil
}
