package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SecurityDo/fluency_api/client"
	fsb "github.com/SecurityDo/fluency_api/fsb"
	"github.com/SecurityDo/fluency_api/model"
)

func newAuthServiceForTest(t *testing.T, handler http.HandlerFunc) *AuthService {
	t.Helper()
	ts := httptest.NewServer(handler)
	t.Cleanup(ts.Close)

	flClient := client.NewFluencyClient(ts.URL, "", false, nil)
	return NewAuthService(flClient)
}

func TestAuthService_AddUser(t *testing.T) {
	svc := newAuthServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/auth/userAdd" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}

		var req fsb.CallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Kargs == nil {
			t.Fatalf("missing kargs in call request")
		}
		var payload AddUserRequest
		if err := json.Unmarshal(req.Kargs.GetBytes(), &payload); err != nil {
			t.Fatalf("failed to decode kargs: %v", err)
		}
		if payload.User == nil || payload.User.Username != "alice" {
			t.Fatalf("unexpected user payload: %+v", payload.User)
		}

		respPayload := fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte([]byte(`{}`))}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(respPayload); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	err := svc.AddUser(&AddUserRequest{User: &model.UserEntry{Username: "alice"}})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestAuthService_ListUser(t *testing.T) {
	svc := newAuthServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/auth/userList" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}

		response := ListUserResponse{Users: []*model.UserEntry{{Username: "alice"}}}
		body, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		respPayload := fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte(body)}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(respPayload); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	users, err := svc.ListUser()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 1 || users[0].Username != "alice" {
		t.Fatalf("unexpected user list: %+v", users)
	}
}

func TestAuthService_GetUser(t *testing.T) {
	svc := newAuthServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/auth/getUser" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}

		response := GetUserResponse{User: &model.UserEntry{Username: "alice"}}
		body, err := json.Marshal(response)
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		respPayload := fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte(body)}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(respPayload); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	user, err := svc.GetUser("alice")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil || user.Username != "alice" {
		t.Fatalf("unexpected user: %+v", user)
	}
}

func TestAuthService_DeleteUser(t *testing.T) {
	svc := newAuthServiceForTest(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/auth/userDelete" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}

		var req fsb.CallRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Kargs == nil {
			t.Fatalf("missing kargs in call request")
		}
		var payload DeleteUserRequest
		if err := json.Unmarshal(req.Kargs.GetBytes(), &payload); err != nil {
			t.Fatalf("failed to decode kargs: %v", err)
		}
		if payload.Username != "alice" {
			t.Fatalf("unexpected username %s", payload.Username)
		}

		respPayload := fsb.CallResponse{Verdict: "OK", Response: fsb.NewJNodeByte([]byte(`{}`))}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(respPayload); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	if err := svc.DeleteUser("alice"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
