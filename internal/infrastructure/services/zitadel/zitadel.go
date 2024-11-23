package zitadel

import (
	"context"
	"fmt"
	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

type ServiceZitadel struct {
	Auth       *authorization.Authorizer[*oauth.IntrospectionContext]
	Middleware *middleware.Interceptor[*oauth.IntrospectionContext]
}

func NewZitadel() (*ServiceZitadel, error) {
	ctx := context.Background()
	conf := config.GetInstance().Zitadel

	zitadelConf := zitadel.New(conf.Domain)

	authZ, err := authorization.New(ctx, zitadelConf, oauth.DefaultAuthorization(conf.Key))
	if err != nil {
		return nil, fmt.Errorf("could not create auth zitadel: %w", err)
	}
	mw := middleware.New(authZ)
	return &ServiceZitadel{Auth: authZ, Middleware: mw}, nil
}
