package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/vitalsignapp/vitalsign-api/pkg/applog"
	"github.com/vitalsignapp/vitalsign-api/response"
)

const (
	authorization = "Authorization"
)

var (
	hmacSecret = "drowssap"
)

func Authorization(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authorization)

		if header == "" {
			applog.Error.Log(r.Context(), "authorization", "authorization token required")
			response.Unauthorized(w, errors.New("authorization token required"))
			return
		}

		splitedHeader := strings.Split(header, "Bearer ")

		if len(splitedHeader) != 2 {
			applog.Error.Log(r.Context(), "authorization", "authorization token required")
			response.Unauthorized(w, errors.New("bearer"))
			return
		}

		_, err := jwt.Parse(header[7:], func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(hmacSecret), nil
		})
		if err != nil {
			applog.Error.Log(r.Context(), "authorization", "authorization token required")
			response.Unauthorized(w, errors.New("token is not valid"))
			return
		}

		handler.ServeHTTP(w, r)
	})
}
