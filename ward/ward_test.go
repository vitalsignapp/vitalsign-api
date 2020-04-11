package ward

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cloud.google.com/go/firestore"
)

type FirestoreClient struct {
	client *firestore.Client
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
