package ward

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWardByHospitalKey(t *testing.T) {
	t.Run("it should return httpCode 200 when call /ward/{hospitalKey}", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/ward/hospitalKeyMock", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(Rooms(mockRepository))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}

	})

	t.Run("it should return 2 wards when hospitalKeyMock has 2 ward", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/ward/hospitalKeyMock", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(Rooms(mockRepository))
		handler.ServeHTTP(resp, req)

		var res []Ward
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if len(res) != 2 {
			t.Errorf("Length of res isn't 2 but got %d", len(res))
		}
	})

	t.Run("it should return 0 wards when hospitalKeyMock has 0 ward or hasn't hospitalKeyMock in database", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/ward/hospitalKeyMock", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(Rooms(mockEmptyRepository))
		handler.ServeHTTP(resp, req)

		var res []Ward
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if len(res) != 0 {
			t.Errorf("Length of res isn't 0 but got %d", len(res))
		}
	})
}

func TestNewRoom(t *testing.T) {
	t.Run("it should return httpCode 200 when call POST/ward", func(t *testing.T) {
		body := `
		{
			"name": "room name",
			"date": {
				"date": "24/04/2020",
				"microtime": 1587662257000,
				"week": 17
			},
			"hospitalKey": "7yfcpkXkME2",
			"addTime": 1587662257000
		}
		`
		req, err := http.NewRequest(http.MethodPost, "/ward", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(NewRoom(mockNewRoomRepository))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("it should return httpCode 500 when have some error in ward repository", func(t *testing.T) {
		body := `
		{
			"name": "room name",
			"date": {
				"date": "24/04/2020",
				"microtime": 1587662257000,
				"week": 17
			},
			"hospitalKey": "7yfcpkXkME2",
			"addTime": 1587662257000
		}
		`
		req, err := http.NewRequest(http.MethodPost, "/ward", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(NewRoom(mockNewRoomRepositoryError))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusInternalServerError {
			t.Errorf("wrong code: got %v want %v", status, http.StatusInternalServerError)
		}
	})

	t.Run("it should return httpCode 400 when request body is empty", func(t *testing.T) {
		body := ``
		req, err := http.NewRequest(http.MethodPost, "/ward", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(NewRoom(mockNewRoomRepository))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("wrong code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func mockRepository(context.Context, string) []Ward {
	mockWards := []Ward{
		Ward{
			ID:          "1",
			Name:        "Test",
			CreatedDate: 1,
		},
		Ward{
			ID:          "2",
			Name:        "Test2",
			CreatedDate: 1,
		},
	}
	return mockWards
}

func mockEmptyRepository(context.Context, string) []Ward {
	mockWards := []Ward{}
	return mockWards
}

func mockNewRoomRepository(context.Context, RoomRequest) error {
	return nil
}

func mockNewRoomRepositoryError(context.Context, RoomRequest) error {
	return errors.New("something error")
}
