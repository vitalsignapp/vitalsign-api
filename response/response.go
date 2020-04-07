package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Msg)
}

func NewErrorResponse(code int, msg string) error {
	return &ErrorResponse{
		Code: code,
		Msg:  msg,
	}
}

func InternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(NewErrorResponse(http.StatusInternalServerError, err.Error()))
}

func Unauthorized(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(NewErrorResponse(http.StatusUnauthorized, err.Error()))
}

func BadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(NewErrorResponse(http.StatusBadRequest, err.Error()))
}

func Forbidden(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(NewErrorResponse(http.StatusForbidden, err.Error()))
}
