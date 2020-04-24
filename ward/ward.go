package ward

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/vitalsignapp/vitalsign-api/response"

	"github.com/gorilla/mux"
)

type (
	RoomRequest struct {
		Name        string      `json:"name"`
		HospitalKey string      `json:"hospitalKey"`
		AddTime     int         `json:"addTime"`
		Date        DateRequest `json:"date"`
	}
	DateRequest struct {
		Date      string `json:"date"`
		Microtime int    `json:"microtime"`
		Week      int    `json:"week"`
	}
)

func Rooms(repo func(context.Context, string) []Ward) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		data := repo(context.Background(), vars["hospitalKey"])

		json.NewEncoder(w).Encode(&data)
	}
}

func NewRoom(repo func(context.Context, RoomRequest) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rr RoomRequest

		err := json.NewDecoder(r.Body).Decode(&rr)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		err = repo(context.Background(), rr)
		if err != nil {
			response.InternalServerError(w, err)
			return
		}

		json.NewEncoder(w).Encode(http.StatusOK)
	}
}
