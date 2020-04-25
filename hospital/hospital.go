package hospital

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/vitalsignapp/vitalsign-api/auth"
	"github.com/vitalsignapp/vitalsign-api/response"
)

// UpdateHospitalConfigRequest UpdateHospitalConfigRequest
type UpdateHospitalConfigRequest struct {
	VitalSignsArr []HospitalConfig `json:"vitalSignsArr"`
}

// UpdateHospitalConfig UpdateHospitalConfig
func UpdateHospitalConfig(parseToken func(http.ResponseWriter, *http.Request) (auth.TokenParseValue, error), repo func(context.Context, string, []HospitalConfig) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenParse, err := parseToken(w, r)

		if err != nil {
			response.BadRequest(w, err)
			return
		}
		hospitalKey := tokenParse.HospitalKey

		var p UpdateHospitalConfigRequest

		err = json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		err = repo(context.Background(), hospitalKey, p.VitalSignsArr)
		if err != nil {
			response.BadRequest(w, err)
			return
		}
		json.NewEncoder(w).Encode(http.StatusOK)
	}
}
