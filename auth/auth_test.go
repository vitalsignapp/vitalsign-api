package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogout(t *testing.T) {
	t.Run("it should return httpCode 200 and set cookie accessToken is empty when call /logout", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/logout", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(Logout())
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("wrong code: got %v want %v", status, http.StatusOK)
		}

		if cookie := resp.Header().Get("Set-Cookie"); cookie != "access-token=; Path=/" {
			t.Errorf("wrong code: got %v want %v", cookie, "access-token=; Path=/")
		}
	})
}
