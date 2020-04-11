package patient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPatientByRoomKeyHandler(t *testing.T) {
	t.Run("it should return httpCode 200 when call /ward/{patientRoomKey}/patients", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/ward/MockRoom1/patients", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(ByRoomKeyHandler(mockPatients))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}

	})

	t.Run("it should return 2 patients when ward 1 has 2 patients", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/ward/MockRoom1/patients", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(ByRoomKeyHandler(mockPatients))
		handler.ServeHTTP(resp, req)

		var res []Patient
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if len(res) != 2 {
			t.Errorf("Length of res isn't 2 but got %d", len(res))
		}
	})

	t.Run("it should return 0 patients when ward 2 has 0 patients", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/ward/MockRoom2/patients", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(ByRoomKeyHandler(mockEmptyPatients))
		handler.ServeHTTP(resp, req)

		var res []Patient
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if len(res) != 0 {
			t.Errorf("Length of res isn't 0 but got %d", len(res))
		}
	})
}

func TestPatientByIDHandler(t *testing.T) {
	t.Run("it should return httpCode 200 when call /patient/{patientID}", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/patient/Patient1", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(ByIDHandler(mockPatient))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}

	})

	t.Run("it should return patients when patientID has data", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/patient/Patient1", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(ByIDHandler(mockPatient))
		handler.ServeHTTP(resp, req)

		var res *Patient
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if res == nil {
			t.Errorf("it should has data but got nil")
		}
	})

	t.Run("it should return nil when patientID hasn't data", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/patient/Patient2", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(ByIDHandler(mockEmptyPatient))
		handler.ServeHTTP(resp, req)

		var res *Patient
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if res != nil {
			t.Errorf("it should nil but got %v", res)
		}
	})
}

func mockPatients(context.Context, string) []Patient {
	mockPatients := []Patient{
		Patient{
			ID:             "Patient1",
			AccountID:      "Mock001",
			DateOfAdmit:    "31/01/2020",
			DateOfBirth:    "14/05/1998",
			Diagnosis:      "mock diagnosis",
			IsRead:         true,
			IsShowNotify:   true,
			Name:           "John",
			Sex:            "male",
			Surname:        "Smith",
			PatientRoomKey: "MockRoom1",
		},
		Patient{
			ID:             "Patient2",
			AccountID:      "Mock001",
			DateOfAdmit:    "15/01/2020",
			DateOfBirth:    "23/04/1994",
			Diagnosis:      "mock diagnosis",
			IsRead:         true,
			IsShowNotify:   true,
			Name:           "Alice",
			Sex:            "female",
			Surname:        "Eve",
			PatientRoomKey: "MockRoom1",
		},
	}
	return mockPatients
}

func mockEmptyPatients(context.Context, string) []Patient {
	mockPatients := []Patient{}
	return mockPatients
}

func mockPatient(context.Context, string) *Patient {
	return &Patient{
		ID:             "Patient1",
		AccountID:      "Mock001",
		DateOfAdmit:    "2020/01/01",
		DateOfBirth:    "1998/05/01",
		Diagnosis:      "mock diagnosis",
		IsRead:         true,
		IsShowNotify:   true,
		Name:           "John",
		Sex:            "male",
		Surname:        "Smith",
		PatientRoomKey: "MockRoom1",
	}
}

func mockEmptyPatient(context.Context, string) *Patient {
	return nil
}
