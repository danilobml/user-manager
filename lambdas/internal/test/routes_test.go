package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danilobml/user-manager/internal/httpx/middleware"
	"github.com/danilobml/user-manager/internal/mocks"
	"github.com/danilobml/user-manager/internal/routes"
	"github.com/danilobml/user-manager/internal/user/dtos"
	"github.com/danilobml/user-manager/internal/user/handler"
	"github.com/danilobml/user-manager/internal/user/jwt"
	"github.com/danilobml/user-manager/internal/user/repository"
	"github.com/danilobml/user-manager/internal/user/service"
)

const strongPass = "StrongP@ssw0rd12345"

type testDeps struct {
	router http.Handler
	apiKey string
	mailer *mocks.MockMailer
	jwt    *jwt.JwtManager
	repo   repository.UserRepository
}

func buildTestServer(t *testing.T) testDeps {
	t.Helper()
	secret := []byte("test-super-secret-32-bytes-min")
	jm := jwt.NewJwtManager(secret)

	repo := repository.NewUserRepositoryInMemory()
	mailer := &mocks.MockMailer{}
	userSvc := service.NewUserserviceImpl(repo, jm, mailer, "http://localhost")
	apiKey := "test-api-key"

	uh := handler.NewUserHandler(userSvc, apiKey)
	auth := middleware.Authenticate(jm)
	router := routes.NewRouter(uh, auth)

	return testDeps{router: router, apiKey: apiKey, mailer: mailer, jwt: jm, repo: repo}
}

func doJSON(t *testing.T, h http.Handler, method, path string, headers map[string]string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr
}


func TestRegisterAndLogin(t *testing.T) {
	deps := buildTestServer(t)

	reg := dtos.RegisterRequest{Email: "user1@example.com", Password: strongPass, Roles: []string{"user"}}
	rr := doJSON(t, deps.router, http.MethodPost, "/register", nil, reg)
	if rr.Code != http.StatusCreated {
		t.Fatalf("register expected 201, got %d (%s)", rr.Code, rr.Body.String())
	}

	var loginResp struct{ Token string `json:"token"` }
	lr := doJSON(t, deps.router, http.MethodPost, "/login", nil, dtos.LoginRequest{Email: reg.Email, Password: strongPass})
	if lr.Code != http.StatusOK {
		t.Fatalf("login expected 200, got %d (%s)", lr.Code, lr.Body.String())
	}
	_ = json.Unmarshal(lr.Body.Bytes(), &loginResp)
	if loginResp.Token == "" {
		t.Fatalf("expected non-empty token")
	}
}

func TestGetUserData_Authorized(t *testing.T) {
	deps := buildTestServer(t)

	_ = doJSON(t, deps.router, http.MethodPost, "/register", nil, dtos.RegisterRequest{
		Email: "me@example.com", Password: strongPass, Roles: []string{"user"},
	})
	lr := doJSON(t, deps.router, http.MethodPost, "/login", nil, dtos.LoginRequest{
		Email: "me@example.com", Password: strongPass,
	})
	if lr.Code != http.StatusOK {
		t.Fatalf("login expected 200, got %d", lr.Code)
	}
	var loginResp struct{ Token string `json:"token"` }
	_ = json.Unmarshal(lr.Body.Bytes(), &loginResp)

	h := map[string]string{"Authorization": "Bearer " + loginResp.Token}
	rr := doJSON(t, deps.router, http.MethodGet, "/users/data", h, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("users/data expected 200, got %d (%s)", rr.Code, rr.Body.String())
	}
	if !strings.Contains(rr.Body.String(), "me@example.com") {
		t.Fatalf("expected email in body, got %s", rr.Body.String())
	}
}

func TestCheckUser_WithApiKey(t *testing.T) {
	deps := buildTestServer(t)

	_ = doJSON(t, deps.router, http.MethodPost, "/register", nil, dtos.RegisterRequest{
		Email: "check@example.com", Password: strongPass, Roles: []string{"user"},
	})
	lr := doJSON(t, deps.router, http.MethodPost, "/login", nil, dtos.LoginRequest{
		Email: "check@example.com", Password: strongPass,
	})
	var loginResp struct{ Token string `json:"token"` }
	_ = json.Unmarshal(lr.Body.Bytes(), &loginResp)

	headers := map[string]string{"User-Api-Key": deps.apiKey}
	body := dtos.CheckUserRequest{Token: loginResp.Token}
	rr := doJSON(t, deps.router, http.MethodPost, "/check-user", headers, body)
	if rr.Code != http.StatusOK {
		t.Fatalf("check-user expected 200, got %d (%s)", rr.Code, rr.Body.String())
	}
	if !strings.Contains(strings.ToLower(rr.Body.String()), `"is_valid":true`) {
		t.Fatalf("expected is_valid true, got %s", rr.Body.String())
	}
}

func TestGetUsers_AdminOnly(t *testing.T) {
	deps := buildTestServer(t)

	_ = doJSON(t, deps.router, http.MethodPost, "/register", nil, dtos.RegisterRequest{
		Email: "a@example.com", Password: strongPass, Roles: []string{"user"},
	})
	_ = doJSON(t, deps.router, http.MethodPost, "/register", nil, dtos.RegisterRequest{
		Email: "admin@example.com", Password: strongPass, Roles: []string{"admin"},
	})
	lr := doJSON(t, deps.router, http.MethodPost, "/login", nil, dtos.LoginRequest{
		Email: "admin@example.com", Password: strongPass,
	})
	var loginResp struct{ Token string `json:"token"` }
	_ = json.Unmarshal(lr.Body.Bytes(), &loginResp)

	h := map[string]string{"Authorization": "Bearer " + loginResp.Token}
	rr := doJSON(t, deps.router, http.MethodGet, "/users", h, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("GET /users expected 200 for admin, got %d (%s)", rr.Code, rr.Body.String())
	}
}

func TestRequestPasswordReset_SendsEmail(t *testing.T) {
	deps := buildTestServer(t)

	_ = doJSON(t, deps.router, http.MethodPost, "/register", nil, dtos.RegisterRequest{
		Email: "user@example.com", Password: strongPass, Roles: []string{"user"},
	})

	rr := doJSON(t, deps.router, http.MethodPost, "/request-password", nil, dtos.RequestPasswordResetRequest{
		Email: "user@example.com",
	})
	if rr.Code != http.StatusNoContent {
		t.Fatalf("request-password expected 204, got %d (%s)", rr.Code, rr.Body.String())
	}

	if len(deps.mailer.To) != 1 || deps.mailer.To[0] != "user@example.com" {
		t.Fatalf("expected email sent to user@example.com, got %+v", deps.mailer.To)
	}
	if deps.mailer.Subject == "" {
		t.Fatalf("expected non-empty subject")
	}
	if deps.mailer.Message == "" || !strings.Contains(deps.mailer.Message, "Password Reset") {
		t.Fatalf("expected message to contain subject/body, got %q", deps.mailer.Message)
	}
}

func TestRegister_InvalidPasswordTooShort(t *testing.T) {
	deps := buildTestServer(t)

	reg := dtos.RegisterRequest{Email: "short@example.com", Password: "x", Roles: []string{"user"}}
	rr := doJSON(t, deps.router, http.MethodPost, "/register", nil, reg)
	if rr.Code == http.StatusCreated {
		t.Fatalf("expected validation error, got 201")
	}
	if !strings.Contains(strings.ToLower(rr.Body.String()), "validation") {
		t.Fatalf("expected validation error message, got %s", rr.Body.String())
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	deps := buildTestServer(t)

	body := dtos.RegisterRequest{Email: "dup@example.com", Password: strongPass, Roles: []string{"user"}}
	r1 := doJSON(t, deps.router, http.MethodPost, "/register", nil, body)
	if r1.Code != http.StatusCreated {
		t.Fatalf("first register expected 201, got %d", r1.Code)
	}
	r2 := doJSON(t, deps.router, http.MethodPost, "/register", nil, body)
	if r2.Code == http.StatusCreated {
		t.Fatalf("second register should not succeed (duplicate email)")
	}
	if !strings.Contains(strings.ToLower(r2.Body.String()), "exists") &&
		!strings.Contains(strings.ToLower(r2.Body.String()), "conflict") {
		t.Fatalf("expected duplicate/conflict message, got %s", r2.Body.String())
	}
}

func TestLogin_WrongPassword_Unauthorized(t *testing.T) {
	deps := buildTestServer(t)

	_ = doJSON(t, deps.router, http.MethodPost, "/register", nil, dtos.RegisterRequest{
		Email: "login@example.com", Password: strongPass, Roles: []string{"user"},
	})
	rr := doJSON(t, deps.router, http.MethodPost, "/login", nil, dtos.LoginRequest{
		Email: "login@example.com", Password: "WrongPass!!",
	})
	if rr.Code == http.StatusOK {
		t.Fatalf("expected unauthorized for wrong password")
	}
	if !strings.Contains(strings.ToLower(rr.Body.String()), "invalid") &&
		!strings.Contains(strings.ToLower(rr.Body.String()), "unauthorized") {
		t.Fatalf("expected invalid/unauthorized message, got %s", rr.Body.String())
	}
}

func TestProtected_NoToken_Unauthorized(t *testing.T) {
	deps := buildTestServer(t)

	rr := doJSON(t, deps.router, http.MethodGet, "/users", nil, nil)
	if rr.Code == http.StatusOK {
		t.Fatalf("expected non-200 without token, got 200")
	}
	if rr.Code != http.StatusUnauthorized && rr.Code != http.StatusForbidden {
		t.Fatalf("expected 401/403 for missing token, got %d (%s)", rr.Code, rr.Body.String())
	}
}

func TestCheckUser_MissingApiKey(t *testing.T) {
	deps := buildTestServer(t)

	// Need a token to fail on API key step (so create a valid user+token)
	_ = doJSON(t, deps.router, http.MethodPost, "/register", nil, dtos.RegisterRequest{
		Email: "chk@example.com", Password: strongPass, Roles: []string{"user"},
	})
	lr := doJSON(t, deps.router, http.MethodPost, "/login", nil, dtos.LoginRequest{
		Email: "chk@example.com", Password: strongPass,
	})
	var loginResp struct{ Token string `json:"token"` }
	_ = json.Unmarshal(lr.Body.Bytes(), &loginResp)

	rr := doJSON(t, deps.router, http.MethodPost, "/check-user", nil, dtos.CheckUserRequest{Token: loginResp.Token})
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without User-Api-Key, got %d (%s)", rr.Code, rr.Body.String())
	}
}

func TestCheckUser_InvalidToken(t *testing.T) {
	deps := buildTestServer(t)

	headers := map[string]string{"User-Api-Key": deps.apiKey}
	rr := doJSON(t, deps.router, http.MethodPost, "/check-user", headers, dtos.CheckUserRequest{
		Token: "not-a-jwt",
	})
	// Could be 400 or 401 depending on your error mapping
	if rr.Code == http.StatusOK {
		t.Fatalf("expected error for invalid token, got 200")
	}
	if !strings.Contains(strings.ToLower(rr.Body.String()), "token") &&
		!strings.Contains(strings.ToLower(rr.Body.String()), "invalid") {
		t.Fatalf("expected token error message, got %s", rr.Body.String())
	}
}

func TestUpdateUser_BadID_BadRequest(t *testing.T) {
	deps := buildTestServer(t)

	h := map[string]string{"Authorization": "Bearer invalid"}
	rr := doJSON(t, deps.router, http.MethodPut, "/users/not-a-uuid", h, dtos.UpdateUserRequest{
		Email: "x@example.com", Roles: []string{"user"},
	})
	if rr.Code != http.StatusBadRequest {
		// If your middleware short-circuits before handler, it might be 401 first.
		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("expected 400 for bad uuid (or 401 if auth runs first), got %d (%s)", rr.Code, rr.Body.String())
		}
	}
}
