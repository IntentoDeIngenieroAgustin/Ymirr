package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CSRFValidator must validate a token bound to the authenticated session.
// The consuming application must issue unpredictable tokens securely.
type CSRFValidator interface {
	ValidateCSRFToken(ctx context.Context, sessionID, token string) bool
}

type CSRFConfig struct {
	CookieName string
	HeaderName string
}

func RequireCSRF(validator CSRFValidator) gin.HandlerFunc {
	return RequireCSRFWithConfig(validator, CSRFConfig{})
}

func RequireCSRFWithConfig(validator CSRFValidator, cfg CSRFConfig) gin.HandlerFunc {
	if validator == nil {
		panic("ymirr/middleware: nil CSRF validator")
	}
	if cfg.CookieName == "" {
		cfg.CookieName = "session_id"
	}
	if cfg.HeaderName == "" {
		cfg.HeaderName = "X-CSRF-Token"
	}
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
			return
		}
		sessionID, err := c.Cookie(cfg.CookieName)
		token := c.GetHeader(cfg.HeaderName)
		if err != nil || sessionID == "" || token == "" || !validator.ValidateCSRFToken(c.Request.Context(), sessionID, token) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid csrf token"})
			return
		}
		c.Next()
	}
}
