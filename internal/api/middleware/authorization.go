package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oprimogus/cardapiogo/internal/core/user"
	xerrors "github.com/oprimogus/cardapiogo/internal/errors"
	logger "github.com/oprimogus/cardapiogo/pkg/log"
)

func AuthorizationMiddleware(allowedRoles []user.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		transactionID := c.GetString(string(logger.TransactionIDKey))
		userRolesContext, exists := c.Get("userRoles")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, xerrors.Forbidden(transactionID, ""))
			return
		}

		userRoles, ok := userRolesContext.([]user.Role)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, xerrors.Forbidden(transactionID, ""))
			return
		}
		if len(userRoles) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, xerrors.Forbidden(transactionID, ""))
			return
		}

		isAllowed := false
		for _, userRole := range userRoles {
			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					isAllowed = true
				}
			}
		}
		if !isAllowed && len(allowedRoles) != 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, xerrors.Forbidden(transactionID, ""))
			return
		}
		c.Next()

	}
}
