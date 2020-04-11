package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/dgrijalva/jwt-go"
	"github.com/vitalsignapp/vitalsign-api/response"
	"google.golang.org/api/iterator"
)

type Credential struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	HospitalKey string `json:"hospitalKey"`
}

type LoginResponse struct {
	ID               string `json:"id"`
	DateCreated      string `json:"dateCreated"`
	Email            string `json:"email"`
	HospitalKey      string `json:"hospitalKey"`
	MicrotimeCreated int64  `json:"microtimeCreated"`
	Name             string `json:"name"`
	Surname          string `json:"surname"`
	UserID           string `json:"userId"`
}

func Authen(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenString, err := token.SignedString([]byte(hmacSecret))
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func Login(fs *firestore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credential Credential
		err := json.NewDecoder(r.Body).Decode(&credential)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		ctx := context.Background()
		iter := fs.Collection("userData").Where("email", "==", credential.Email).
			Where("password", "==", credential.Password).
			Documents(ctx)

		defer iter.Stop()

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			
			if err != nil {
				break
			}

			p := LoginResponse{}
			p.ID = doc.Ref.ID
			err = doc.DataTo(&p)
			if err != nil {
				break
			}

			json.NewEncoder(w).Encode(&p)
			return
		}

		response.Unauthorized(w, nil)
		return
	}
}
