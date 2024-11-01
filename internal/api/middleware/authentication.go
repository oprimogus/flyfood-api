package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/oprimogus/cardapiogo/internal/core/authentication"
	"github.com/oprimogus/cardapiogo/internal/core/user"
	xerrors "github.com/oprimogus/cardapiogo/internal/errors"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

func AuthenticationMiddleware(repository authentication.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetString(string(logger.TraceIDKey))
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				xerrors.Unauthorized(traceID, ""))
			return
		}
		token = strings.Replace(token, "Bearer ", "", -1)
		isValidToken, err := repository.IsValidToken(c, token)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				xerrors.Unauthorized(traceID, err.Error()))
			return
		}

		if !isValidToken {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				xerrors.Unauthorized(traceID, "Invalid access token"))
			return
		}

		claims, err := repository.DecodeAccessToken(c, token)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				xerrors.New(traceID, http.StatusUnauthorized, err.Error()))
			return
		}

		realmAccess, ok := claims["realm_access"].(map[string]interface{})
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				xerrors.New(traceID, http.StatusUnauthorized, "Invalid token"))
			return
		}

		roles, ok := realmAccess["roles"].([]interface{})
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				xerrors.New(traceID, http.StatusUnauthorized, "Invalid token"))
			return
		}

		var userRoles []user.Role
		for _, role := range roles {
			if roleStr, ok := role.(string); ok {
				if user.IsValidRole(roleStr) {
					userRoles = append(userRoles, user.Role(roleStr))
				}

			}
		}
		c.Set("user_id", claims["sub"])
		c.Set("user_roles", userRoles)
		c.Next()
	}
}
