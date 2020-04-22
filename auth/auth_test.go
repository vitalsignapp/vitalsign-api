package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Run("it should return httpCode 200 and set cookie access-token when call /login", func(t *testing.T) {
		b := `
		{
			"email": "mock@mock.com",
			"password": "supersecret",
			"hospitalKey": "kcom"
		}
		`
		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(b))
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(Login(mockUserData))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}

		if cookie := resp.Header().Get("Set-Cookie"); len(cookie) < 20 {
			t.Errorf("wrong code: got %v want %v", cookie, "access-token=; Path=/")
		}
	})

	t.Run("it should return httpCode 401 when call /login with wrong email", func(t *testing.T) {
		b := `
		{
			"email": "mock1@mock.com",
			"password": "supersecret",
			"hospitalKey": "kcom"
		}
		`
		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(b))
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(Login(mockEmptyUserData))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusUnauthorized {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("it should return httpCode 500 when call /login with empty body", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(""))
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(Login(mockEmptyUserData))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}
	})
}

func mockUserData(context.Context, string, string, string) *LoginResponse {
	u := LoginResponse{
		ID:               "mock",
		DateCreated:      "02/04/2020",
		Email:            "mock@mock.com",
		HospitalKey:      "kcom",
		MicrotimeCreated: 1585839201000,
		Name:             "Mr.Mock",
		Surname:          "Mock",
		UserID:           "11211",
	}
	return &u
}

func mockEmptyUserData(context.Context, string, string, string) *LoginResponse {
	return nil
}

func TestLogout(t *testing.T) {
	t.Run("it should return httpCode 200 and set cookie accessToken is empty when call /logout", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/logout", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(Logout())
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}

		if cookie := resp.Header().Get("Set-Cookie"); cookie != "access-token=; Path=/" {
			t.Errorf("wrong code: got %v want %v", cookie, "access-token=; Path=/")
		}
	})
}
