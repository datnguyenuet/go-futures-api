package middleware

import (
	"github.com/gin-gonic/gin"
	httpErrors "go-futures-api/pkg/http_errors"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func (m *Manager) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerHeader := c.Request.Header.Get("Authorization")

		m.logger.Infof("auth middleware bearerHeader %s", bearerHeader)
		if bearerHeader != "" {
			headerParts := strings.Split(bearerHeader, " ")
			if len(headerParts) != 2 {
				m.logger.Error("auth middleware", zap.String("headerParts", "len(headerParts) != 2"))
				m.respondWithError(c, http.StatusUnauthorized, httpErrors.ErrUnauthorized)
			}
			tokenString := headerParts[1]
			if err := m.validateJWTToken(tokenString); err != nil {
				m.logger.Error("middleware validateJWTToken", zap.String("headerJWT", err.Error()))
				m.respondWithError(c, http.StatusUnauthorized, httpErrors.ErrUnauthorized)
			}
			c.Next()
			return
		}

		cookie, err := c.Cookie("jwt-token")
		if err != nil {
			m.logger.Error("c.Cookie", err.Error())
			m.respondWithError(c, http.StatusUnauthorized, httpErrors.ErrUnauthorized)
			return
		}

		if err = m.validateJWTToken(cookie); err != nil {
			m.logger.Error("validateJWTToken", err.Error())
			m.respondWithError(c, http.StatusUnauthorized, httpErrors.ErrUnauthorized)
			return
		}
		c.Next()
	}
}

func (m *Manager) validateJWTToken(token string) error {
	return nil
}

func (m *Manager) respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"description": message})
}

func (m *Manager) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		userId := c.Request.Header.Get("x-bce-uid")
		if userId == "" {
			m.logger.Error("auth middleware", zap.String("headerParts", "Missing x-bce-uid"))
			m.respondWithError(c, http.StatusUnauthorized, httpErrors.ErrUnauthorized)
		}

		user, err := m.authUseCase.FindOne(ctx, userId)
		if err != nil {
			m.logger.Error("auth middleware", err.Error())
			m.respondWithError(c, http.StatusBadRequest, httpErrors.ErrBadRequest)
			return
		}

		if user == nil {
			m.logger.Error("auth middleware", zap.String("headerParts", "NotFound user"))
			m.respondWithError(c, http.StatusUnauthorized, httpErrors.ErrUnauthorized)
			return
		}

		if user.AccountLv < 2 {
			m.logger.Error("auth middleware", zap.String("headerJWT", "Not enough permissions"))
			m.respondWithError(c, http.StatusForbidden, httpErrors.ErrForbidden)
			return
		}

		if user.AuthenticatorVerifyStatus != 1 {
			m.logger.Error("auth middleware", zap.String("headerJWT", "User not verified"))
			m.respondWithError(c, http.StatusUnauthorized, httpErrors.ErrUnauthorized)
			return
		}
		c.Next()
	}
}
