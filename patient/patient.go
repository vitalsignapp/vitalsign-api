package patient

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func Patients(repo func(context.Context, string) []Patient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		data := repo(context.Background(), vars["patientRoomKey"])

		json.NewEncoder(w).Encode(&data)
	}
}

func PatientByID(repo func(context.Context, string) *Patient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		data := repo(context.Background(), vars["patientID"])

		json.NewEncoder(w).Encode(&data)
	}
}
