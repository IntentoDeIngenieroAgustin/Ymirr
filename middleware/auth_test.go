package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockValidator struct{}

func (mockValidator) ValidateSessionID(_ context.Context, s string) (string, error) {
	if s == "valid" {
		return "42", nil
	}
	return "", errors.New("invalid")
}
func TestRequireAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequireAuth(mockValidator{}))
	r.GET("/", func(c *gin.Context) { c.String(200, c.GetString("user_id")) })
	for _, tt := range []struct {
		cookie string
		status int
		body   string
	}{{"", 401, ""}, {"bad", 401, ""}, {"valid", 200, "42"}} {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if tt.cookie != "" {
			req.AddCookie(&http.Cookie{Name: "session_id", Value: tt.cookie})
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tt.status {
			t.Errorf("cookie %q: got %d want %d", tt.cookie, w.Code, tt.status)
		}
		if tt.body != "" && w.Body.String() != tt.body {
			t.Errorf("body: %q", w.Body.String())
		}
	}
}
