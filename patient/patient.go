package patient

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

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
