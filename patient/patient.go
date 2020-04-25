package patient

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vitalsignapp/vitalsign-api/auth"
	"github.com/vitalsignapp/vitalsign-api/response"
)

type PatientRequest struct {
	ID             string `json:"id"`
	DateOfAdmit    string `json:"dateOfAdmit"`
	DateOfBirth    string `json:"dateOfBirth"`
	Diagnosis      string `json:"diagnosis"`
	IsRead         bool   `json:"isRead"`
	IsShowNotify   bool   `json:"isShowNotify"`
	Name           string `json:"name"`
	PatientRoomKey string `json:"patientRoomKey"`
	Sex            string `json:"sex"`
	Surname        string `json:"surname"`
	Username       string `json:"username"`
}

type PatientStatusRequest struct {
	IsRead   *bool `json:"isRead"`
	IsNotify *bool `json:"isNotify"`
}

func ByRoomKeyHandler(repo func(context.Context, string) []Patient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		data := repo(context.Background(), vars["patientRoomKey"])

		json.NewEncoder(w).Encode(&data)
	}
}

func ByIDHandler(repo func(context.Context, string) *Patient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		data := repo(context.Background(), vars["patientID"])

		json.NewEncoder(w).Encode(&data)
	}
}

func ByHospital(repo func(context.Context, string) []Patient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		data := repo(context.Background(), vars["hospitalID"])

		json.NewEncoder(w).Encode(&data)
	}
}

func Update(repo func(context.Context, string, PatientRequest) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p PatientRequest
		vars := mux.Vars(r)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		err = repo(context.Background(), vars["patientID"], p)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		json.NewEncoder(w).Encode(http.StatusOK)
	}
}

func LogByIDHandler(repo func(context.Context, string) []PatientLog) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		data := repo(context.Background(), vars["patientID"])

		json.NewEncoder(w).Encode(&data)
	}
}

// UpdatePatientStatus UpdatePatientStatus
func UpdatePatientStatus(parseToken func(http.ResponseWriter, *http.Request) (auth.TokenParseValue, error), repo func(context.Context, string, string, PatientStatusRequest) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenParse, err := parseToken(w, r)

		if err != nil {
			response.BadRequest(w, err)
			return
		}
		hospitalKey := tokenParse.HospitalKey

		var p PatientStatusRequest
		err = json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		if p.IsRead == nil {
			response.BadRequest(w, errors.New("isRead is required"))
			return
		}

		if p.IsNotify == nil {
			response.BadRequest(w, errors.New("isNotify is required"))
			return
		}

		vars := mux.Vars(r)
		err = repo(context.Background(), hospitalKey, vars["patientID"], p)

		if err != nil {
			response.BadRequest(w, err)
			return
		}

		json.NewEncoder(w).Encode(http.StatusOK)
	}
}
