package middleware

import (
	"net/http"

	"github.com/oprimogus/flyfood-api/internal/infra/services/zitadel"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
)

func Authentication(next http.Handler) http.Handler {
	return zitadel.GetInstance().Middleware.RequireAuthorization()(next)
}

func Authorization(authorities ...string) func(http.Handler) http.Handler {
	zt := zitadel.GetInstance()

	return func(next http.Handler) http.Handler {
		if len(authorities) == 0 {
			panic("Authorization middleware requires at least one authority")
		}

		options := make([]authorization.CheckOption, len(authorities))
		for i, authority := range authorities {
			options[i] = authorization.WithRole(authority)
		}

		return zt.Middleware.RequireAuthorization(options...)(next)
	}
}
