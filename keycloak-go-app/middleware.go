package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "No Authorization header provided"},
			)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "" {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				gin.H{"error": "No token provided"},
			)
			return
		}

		validToken, err := keycloakRetrospectToken(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		if !validToken {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}
	}
}

func TokenRevoke() gin.HandlerFunc {
	return func(c *gin.Context) {

		auth := c.Request.Header.Get("Authorization")
		token := strings.TrimPrefix(auth, "Bearer ")

		err := keycloakClientTokenRevoke(token)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

	}
}
