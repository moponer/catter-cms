package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"catter-cms-gin/config"
	"catter-cms-gin/internal/handlers/api"
	"github.com/gin-gonic/gin"
)

const (
	csrfCookieName = "catter_csrf"
	csrfHeaderName = "X-CSRF-Token"
	csrfTokenLen   = 32
)

func generateCSRFToken() (string, error) {
	b := make([]byte, csrfTokenLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// CSRFProtection implements Double Submit Cookie pattern for CSRF defense.
func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.Get()

		if c.Request.Method == http.MethodGet ||
			c.Request.Method == http.MethodHead ||
			c.Request.Method == http.MethodOptions {

			if _, err := c.Cookie(csrfCookieName); err != nil {
				token, genErr := generateCSRFToken()
				if genErr != nil {
					api.InternalError(c, "CSRF Token generation failed")
					c.Abort()
					return
				}
				setCSRFCookie(c, token, cfg)
			}
			c.Next()
			return
		}

		cookieToken, err := c.Cookie(csrfCookieName)
		if err != nil || cookieToken == "" {
			api.Forbidden(c, "Missing CSRF Token, please refresh the page")
			c.Abort()
			return
		}

		headerToken := c.GetHeader(csrfHeaderName)
		if headerToken == "" {
			api.Forbidden(c, "Missing X-CSRF-Token header")
			c.Abort()
			return
		}

		if cookieToken != headerToken {
			api.Forbidden(c, "CSRF Token mismatch, please refresh the page")
			c.Abort()
			return
		}

		newToken, genErr := generateCSRFToken()
		if genErr == nil {
			setCSRFCookie(c, newToken, cfg)
			c.Header(csrfHeaderName, newToken)
		}

		c.Next()
	}
}

func setCSRFCookie(c *gin.Context, token string, cfg *config.Config) {
	cookie := &http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   86400,
		Secure:   cfg.IsProduction(),
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(c.Writer, cookie)
}
