package ward

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func Rooms(repo func(context.Context, string) []map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		data := repo(context.Background(), vars["hospitalKey"])

		json.NewEncoder(w).Encode(&data)
	}
}
