package user

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChangePassword(t *testing.T) {

	t.Run("it should return httpCode 200 when call /userData/reset/{userID}", func(t *testing.T) {
		body := `{
			"password": "1212312121hey!",
			"confirmPassword": "1212312121hey!"
		}`
		// b, _ := json.Marshal(body)
		req, err := http.NewRequest(http.MethodPut, "/userData/reset/wtf", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(ChangePassword(mockRepository))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("it should throw error when password is empty", func(t *testing.T) {
		body := `{
			"password": "",
			"confirmPassword": "1212312121hey!"
		}`
		// b, _ := json.Marshal(body)
		req, err := http.NewRequest(http.MethodPut, "/userData/reset/wtf", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(ChangePassword(mockRepository))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("wrong code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("it should throw error when password and confirm password does not match", func(t *testing.T) {
		body := `{
			"password": "1212111",
			"confirmPassword": "1212312121hey!"
		}`
		// b, _ := json.Marshal(body)
		req, err := http.NewRequest(http.MethodPut, "/userData/reset/wtf", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(ChangePassword(mockRepository))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("wrong code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("it should bad request when body is empty", func(t *testing.T) {
		body := ``
		// b, _ := json.Marshal(body)
		req, err := http.NewRequest(http.MethodPut, "/userData/reset/wtf", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(ChangePassword(mockRepository))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("wrong code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("it should bad request when repo return error", func(t *testing.T) {
		body := `{
			"password": "1212312121hey!",
			"confirmPassword": "1212312121hey!"
		}`
		// b, _ := json.Marshal(body)
		req, err := http.NewRequest(http.MethodPut, "/userData/reset/wtf", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(ChangePassword(mockRepositoryError))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("wrong code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func mockRepository(context.Context, string, string) error {
	return nil
}

func mockRepositoryError(context.Context, string, string) error {
	return errors.New("Error")
}
