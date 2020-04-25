package hospital

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vitalsignapp/vitalsign-api/auth"
)

func TestUpdateHospitalConfig(t *testing.T) {

	t.Run("it should return httpCode 200 when call POST /hospital", func(t *testing.T) {
		body := `{
			"vitalSignsArr": [
		    {
		        "status": true,
		        "sym": "xxxx1"
		    },
		    {
		        "status": true,
		        "sym": "xxxx2"
				}
			]
		}`

		req, err := http.NewRequest(http.MethodPost, "/hospital", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(UpdateHospitalConfig(mockParseToken, mockRepository))
		handler.ServeHTTP(resp, req)
		if status := resp.Code; status != http.StatusOK {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("it should return httpCode 400 when call POST /hospital and parse token is invalid", func(t *testing.T) {
		body := `{
			"vitalSignsArr": [
		    {
		        "status": true,
		        "sym": "xxxx1"
		    },
		    {
		        "status": true,
		        "sym": "xxxx2"
				}
			]
		}`

		req, err := http.NewRequest(http.MethodPost, "/hospital", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(UpdateHospitalConfig(mockParseTokenError, mockRepository))
		handler.ServeHTTP(resp, req)
		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("wrong code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("it should return httpCode 400 when call POST /hospital and body is incorrect", func(t *testing.T) {
		body := `{
			"vitalSignsArr": {}
		}`

		req, err := http.NewRequest(http.MethodPost, "/hospital", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(UpdateHospitalConfig(mockParseToken, mockRepository))
		handler.ServeHTTP(resp, req)
		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("wrong code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("it should return httpCode 400 when call POST /hospital and update value is failed", func(t *testing.T) {
		body := `{
			"vitalSignsArr": [
		    {
		        "status": true,
		        "sym": "xxxx1"
		    },
		    {
		        "status": true,
		        "sym": "xxxx2"
				}
			]
		}`

		req, err := http.NewRequest(http.MethodPost, "/hospital", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(UpdateHospitalConfig(mockParseToken, mockRepositoryError))
		handler.ServeHTTP(resp, req)
		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("wrong code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func mockParseTokenError(w http.ResponseWriter, r *http.Request) (auth.TokenParseValue, error) {
	return auth.TokenParseValue{}, errors.New("Error")
}

func mockParseToken(w http.ResponseWriter, r *http.Request) (auth.TokenParseValue, error) {
	return auth.TokenParseValue{
		Email:       "email@email.com",
		HospitalKey: "123",
	}, nil
}

func mockRepository(context.Context, string, []HospitalConfig) error {
	return nil
}

func mockRepositoryError(context.Context, string, []HospitalConfig) error {
	return errors.New("Error")
}
