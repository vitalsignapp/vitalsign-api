package auth

import (
	"encoding/json"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/dgrijalva/jwt-go"
	"github.com/vitalsignapp/vitalsign-api/response"
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
		p := LoginResponse{
			ID:               "AjGMr39iDFclLYPJrTGb",
			DateCreated:      "02/04/2020",
			Email:            "test1@gmail.com",
			HospitalKey:      "7yfcpkXkME2OrvbYNAq1",
			MicrotimeCreated: 1585839201000,
			Name:             "สมหมาย",
			Surname:          "ขายฝัน",
			UserID:           "999",
		}
		json.NewEncoder(w).Encode(&p)
		// var credential Credential
		// err := json.NewDecoder(r.Body).Decode(&credential)
		// if err != nil {
		// 	response.BadRequest(w,err)
		// 	return
		// }

		// ctx := context.Background()
		// iter := fs.Collection("userData").Where("email", "==", credential.Email).
		// 	Where("password", "==", credential.Password).
		// 	Where("hospitalKey", "==", credential.HospitalKey).
		// 	Documents(ctx)

		// defer iter.Stop()

		// for {
		// 	doc, err := iter.Next()
		// 	if err == iterator.Done {
		// 		break
		// 	}
		// 	if err != nil {
		// 		continue
		// 	}

		// 	p := LoginResponse{}
		// 	err = doc.DataTo(&p)
		// 	if err != nil {
		// 		continue
		// 	}

		// 	p.ID = doc.Ref.ID

		// 	json.NewEncoder(w).Encode(&p)
		// 	return
		// }

		// response.Unauthorized(w, nil)
		// return
	}
}
