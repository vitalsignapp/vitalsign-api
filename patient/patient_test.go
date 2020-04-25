package patient

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vitalsignapp/vitalsign-api/auth"
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

func TestPatientLogByIDHandler(t *testing.T) {
	t.Run("it should return httpCode 200 when call /patient/{patientID}/log", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/patient/Patient1/log", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(LogByIDHandler(mockPatientLogs))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}

	})

	t.Run("it should return 2 logs when patientID has data", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/patient/Patient1/log", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(LogByIDHandler(mockPatientLogs))
		handler.ServeHTTP(resp, req)

		var res []PatientLog
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if res == nil {
			t.Errorf("it should has data but got nil")
		}
	})

	t.Run("it should return 0 logs when Patient 2 has 0 log", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/patient/Patient2/log", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(LogByIDHandler(mockEmptyPatientLogs))
		handler.ServeHTTP(resp, req)

		var res []PatientLog
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if len(res) != 0 {
			t.Errorf("Length of res isn't 0 but got %d", len(res))
		}
	})
}

func TestPatientsByHospital(t *testing.T) {
	t.Run("it should return httpCode 200 when call /patient/hospital/{HospitalID}", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/patient/hospital/mockHospitalID1", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(ByHospital(mockPatients))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}

	})

	t.Run("it should return 4 patients when this hospital has 4 patients", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/patient/hospital/mockHospitalID1", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(ByHospital(mockPatients))
		handler.ServeHTTP(resp, req)

		var res []Patient
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		if len(res) != 2 {
			t.Errorf("Length of res isn't 4 but got %d", len(res))
		}
	})

	t.Run("it should return 0 patients when this hospital has 0 patients", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/patient/hospital/mockHospitalID2", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(ByHospital(mockEmptyPatientsByHospital))
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

func TestUpdatePatientStatus(t *testing.T) {
	t.Run("should return status 200 when call POST /patient/{patientID}/status", func(t *testing.T) {
		body := `{
			"isRead": true,
			"isNotify": false
		}`
		req, err := http.NewRequest(http.MethodPatch, "/patient/123/status", strings.NewReader((body)))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(UpdatePatientStatus(mockParseToken, mockUpdateStatusRequest))
		handler.ServeHTTP(resp, req)
		if status := resp.Code; status != http.StatusOK {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("should return status 400 when call POST /patient/{patientID}/status and body does not set field isRead and isNotify", func(t *testing.T) {
		body := `{}`
		req, err := http.NewRequest(http.MethodPatch, "/patient/123/status", strings.NewReader((body)))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(UpdatePatientStatus(mockParseToken, mockUpdateStatusRequest))
		handler.ServeHTTP(resp, req)
		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("wrong code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("should return status 400 when call POST /patient/{patientID}/status and body is incorrect", func(t *testing.T) {
		body := `[]`
		req, err := http.NewRequest(http.MethodPatch, "/patient/123/status", strings.NewReader((body)))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(UpdatePatientStatus(mockParseToken, mockUpdateStatusRequest))
		handler.ServeHTTP(resp, req)
		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("wrong code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("should return status 400 when call POST /patient/{patientID}/status and body is incorrect", func(t *testing.T) {
		body := `{
			"isRead": true,
			"isNotify": false
		}`
		req, err := http.NewRequest(http.MethodPatch, "/patient/123/status", strings.NewReader((body)))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(UpdatePatientStatus(mockParseToken, mockUpdateStatusRequestError))
		handler.ServeHTTP(resp, req)
		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("wrong code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}

func TestAddNewPatient(t *testing.T) {
	t.Run("it should return httpCode 200 when add new patient success", func(t *testing.T) {
		body := `
		{
			"dateOfAdmit": "25/04/2020",
			"dateOfBirth": "25/04/2000",
			"diagnosis": "diagnosis test",
			"hospitalKey": "7yfcpkXkME2OrvbYNAq1",
			"isRead": true,
			"isShowNotify": true,
			"name": "ผู้ป่วย",
			"patientRoomKey": "xgsqZBm6SYERhlgYgoLa",
			"sex": "male",
			"surname": "นามสมมุต",
			"username": "900100"
			}
		`
		req, err := http.NewRequest(http.MethodPost, "/patient", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(Create(mockAddNewRepositorySuccess))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}
	})

	t.Run("it should return httpCode 500 when create patient fail", func(t *testing.T) {
		body := `
		{
			"dateOfAdmit": "25/04/2020",
			"dateOfBirth": "25/04/2000",
			"diagnosis": "diagnosis test",
			"hospitalKey": "7yfcpkXkME2OrvbYNAq1",
			"isRead": true,
			"isShowNotify": true,
			"name": "ผู้ป่วย",
			"patientRoomKey": "xgsqZBm6SYERhlgYgoLa",
			"sex": "male",
			"surname": "นามสมมุต",
			"username": "900100"
			}
		`
		req, err := http.NewRequest(http.MethodPost, "/patient", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(Create(mockAddNewRepositoryFail))
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusInternalServerError {
			t.Errorf("wrong code: got %v want %v", status, http.StatusInternalServerError)
		}
	})

	t.Run("it should return httpCode 400 when body request is empty", func(t *testing.T) {
		body := ``
		req, err := http.NewRequest(http.MethodPost, "/patient", strings.NewReader(body))
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(Create(mockAddNewRepositoryFail))
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

func mockUpdateStatusRequest(context.Context, string, string, PatientStatusRequest) error {
	return nil
}

func mockUpdateStatusRequestError(context.Context, string, string, PatientStatusRequest) error {
	return errors.New("Error")
}

func mockPatients(context.Context, string) []Patient {
	mockPatients := []Patient{
		{
			ID:             "Patient1",
			Username:       "Mock001",
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
		{
			ID:             "Patient2",
			Username:       "Mock001",
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
		Username:       "Mock001",
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

func mockPatientLogs(context.Context, string) []PatientLog {
	mockPatientLogs := []PatientLog{
		{
			ID:             "mockId1",
			BloodPressure:  "120/70",
			HeartRate:      "99",
			HospitalKey:    "MockHospitalKey",
			InputDate:      "10/04/2563",
			InputRound:     6,
			Microtime:      1586480268000,
			OtherSymptoms:  "อยากดื่มกาแฟ",
			Oxygen:         "96",
			PatientKey:     "Patient1",
			PatientRoomKey: "MockPatientRoomKey",
			SymptomsCheck: []Symptom{
				{
					Status: false,
					Sym:    "ไข้",
				},
				{
					Status: true,
					Sym:    "ไอ",
				},
			},
			Temperature: "35.0",
		},
		{
			ID:             "mockId2",
			BloodPressure:  "110/60",
			HeartRate:      "82",
			HospitalKey:    "MockHospitalKey",
			InputDate:      "10/04/2563",
			InputRound:     6,
			Microtime:      1586480268000,
			OtherSymptoms:  "อยากดื่มกาแฟเอสเปรสโซ่",
			Oxygen:         "96",
			PatientKey:     "Patient1",
			PatientRoomKey: "MockPatientRoomKey",
			SymptomsCheck: []Symptom{
				{
					Status: false,
					Sym:    "ไข้",
				},
				{
					Status: true,
					Sym:    "ไอ",
				},
			},
			Temperature: "36.0",
		},
	}
	return mockPatientLogs
}

func mockEmptyPatientLogs(context.Context, string) []PatientLog {
	mockPatientLogs := []PatientLog{}
	return mockPatientLogs
}

func mockPatientsByHospital(context.Context, string) []Patient {
	mockPatients := []Patient{
		{
			ID:             "Patient1",
			Username:       "Mock001",
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
		{
			ID:             "Patient2",
			Username:       "Mock001",
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
		{
			ID:             "Patient3",
			Username:       "Mock003",
			DateOfAdmit:    "31/01/2020",
			DateOfBirth:    "14/05/1998",
			Diagnosis:      "mock diagnosis",
			IsRead:         true,
			IsShowNotify:   true,
			Name:           "Johny",
			Sex:            "male",
			Surname:        "Smith",
			PatientRoomKey: "MockRoom2",
		},
		{
			ID:             "Patient4",
			Username:       "Mock004",
			DateOfAdmit:    "15/01/2020",
			DateOfBirth:    "23/04/1994",
			Diagnosis:      "mock diagnosis",
			IsRead:         true,
			IsShowNotify:   true,
			Name:           "Alice",
			Sex:            "female",
			Surname:        "Ever",
			PatientRoomKey: "MockRoom2",
		},
	}
	return mockPatients
}

func mockEmptyPatientsByHospital(context.Context, string) []Patient {
	mockPatients := []Patient{}
	return mockPatients
}

func mockAddNewRepositorySuccess(context.Context, PatientRequest) error {
	return nil
}

func mockAddNewRepositoryFail(context.Context, PatientRequest) error {
	return errors.New("something error")
}
