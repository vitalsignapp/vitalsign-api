package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

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

// ParseToken ParseToken
type TokenParseValue struct {
	Email       string `json:"email"`
	HospitalKey string `json:"hospitalKey"`
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

func GenerateToken(payload map[string]interface{}, expireAt time.Time) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "vitalsign"
	claims["exp"] = expireAt.Unix()

	for k, v := range payload {
		claims[k] = v
	}

	return token.SignedString([]byte(hmacSecret))
}

func ParseToken(w http.ResponseWriter, r *http.Request) (TokenParseValue, error) {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	reqToken = strings.TrimSpace(splitToken[1])
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(hmacSecret), nil
	})
	if err != nil {
		return TokenParseValue{}, err
	}

	return TokenParseValue{
		Email:       claims["email"].(string),
		HospitalKey: claims["hospitalKey"].(string),
	}, nil
}

func Login(repo func(context.Context, string, string, string) *LoginResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credential Credential
		err := json.NewDecoder(r.Body).Decode(&credential)
		if err != nil {
			response.BadRequest(w, err)
			return
		}

		user := repo(context.Background(), credential.Email, credential.Password, credential.HospitalKey)
		if user == nil || user.ID == "" {
			response.Unauthorized(w, errors.New("unauthorized"))
			return
		}

		//Public Claims
		payload := map[string]interface{}{
			"email":       user.Email,
			"hospitalKey": user.HospitalKey,
		}

		token, err := GenerateToken(payload, time.Now().Add(time.Hour*8))
		if err != nil {
			response.InternalServerError(w, err)
			return
		}

		cookie := &http.Cookie{
			Name:  "access-token",
			Value: token,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data":  user,
			"token": token,
		})
		return
	}
}

func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:  "access-token",
			Value: "",
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}
