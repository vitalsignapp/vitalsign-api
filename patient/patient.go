package patient

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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
