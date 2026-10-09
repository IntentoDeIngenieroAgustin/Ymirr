package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SessionValidator is implemented by the consuming application's auth adapter.
type SessionValidator interface {
	ValidateSessionID(context.Context, string) (userID string, err error)
}

type AuthConfig struct {
	CookieName     string
	ContextKey     string
	OnUnauthorized gin.HandlerFunc
}

func RequireAuth(validator SessionValidator) gin.HandlerFunc {
	return RequireAuthWithConfig(validator, AuthConfig{})
}

func RequireAuthWithConfig(validator SessionValidator, cfg AuthConfig) gin.HandlerFunc {
	if validator == nil {
		panic("ymirr/middleware: nil session validator")
	}
	if cfg.CookieName == "" {
		cfg.CookieName = "session_id"
	}
	if cfg.ContextKey == "" {
		cfg.ContextKey = "user_id"
	}
	if cfg.OnUnauthorized == nil {
		cfg.OnUnauthorized = func(c *gin.Context) { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"}) }
	}
	return func(c *gin.Context) {
		sessionID, err := c.Cookie(cfg.CookieName)
		if err != nil || sessionID == "" {
			cfg.OnUnauthorized(c)
			c.Abort()
			return
		}
		userID, err := validator.ValidateSessionID(c.Request.Context(), sessionID)
		if err != nil || userID == "" {
			cfg.OnUnauthorized(c)
			c.Abort()
			return
		}
		c.Set(cfg.ContextKey, userID)
		c.Next()
	}
}
