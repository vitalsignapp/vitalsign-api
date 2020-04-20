package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vitalsignapp/vitalsign-api/response"
)

// ChangePasswordRequest ChangePasswordRequest
type ChangePasswordRequest struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

// ChangePassword ChangePassword
func ChangePassword(repo func(context.Context, string, string) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p ChangePasswordRequest

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		if p.Password == "" || p.ConfirmPassword == "" {
			response.BadRequest(w, errors.New("password or confirm password is empty"))
			return
		}

		if p.Password != p.ConfirmPassword {
			response.BadRequest(w, errors.New("the password and confirmation password do not match"))
			return
		}

		vars := mux.Vars(r)
		err = repo(context.Background(), vars["userID"], p.Password)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		json.NewEncoder(w).Encode(http.StatusOK)
	}
}
