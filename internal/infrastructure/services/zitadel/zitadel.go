package zitadel

import (
	"context"
	"fmt"
	"github.com/oprimogus/flyfood-api/internal/config"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/client"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
	"slices"
)

var Instance *ServiceZitadel

type Role string

const (
	Owner Role = "owner"
)

type ServiceZitadel struct {
	Auth       *authorization.Authorizer[*oauth.IntrospectionContext]
	Client     *client.Client
	Middleware *middleware.Interceptor[*oauth.IntrospectionContext]
}

func NewZitadel() (*ServiceZitadel, error) {
	ctx := context.Background()
	conf := config.GetInstance().Zitadel

	zitadelConf := zitadel.New(conf.Domain, zitadel.WithInsecure(conf.Port))

	authZ, err := authorization.New(ctx, zitadelConf, oauth.DefaultAuthorization(conf.KeyPath))
	if err != nil {
		return nil, fmt.Errorf("could not create auth zitadel: %w", err)
	}
	mw := middleware.New(authZ)

	c, err := client.New(ctx, zitadelConf,
		client.WithAuth(client.DefaultServiceUserAuthentication(conf.ServiceAccountKeyPath, oidc.ScopeOpenID, client.ScopeZitadelAPI())),
	)
	if err != nil {
		return nil, fmt.Errorf("zitadel sdk could not initialize client: %w", err)
	}
	return &ServiceZitadel{
		Auth:       authZ,
		Client:     c,
		Middleware: mw}, nil
}

func GetInstance() *ServiceZitadel {
	if Instance == nil {
		instance, err := NewZitadel()
		if err != nil {
			panic(err)
		}
		Instance = instance
	}
	return Instance
}

func (s *ServiceZitadel) GetContext(ctx context.Context) *oauth.IntrospectionContext {
	return s.Middleware.Context(ctx)
}

func (s *ServiceZitadel) SetRole(ctx context.Context, userID string, role Role) error {
	zt := s.Client
	conf := config.GetInstance().Zitadel

	listUserGrantsArg := management.ListUserGrantRequest{}

	userGrants, err := zt.ManagementService().ListUserGrants(ctx, &listUserGrantsArg)
	if err != nil {
		return fmt.Errorf("could not list user grants: %w", err)
	}

	var grantID string
	for _, grant := range userGrants.Result {
		if grant.ProjectId == conf.ProjectID && grant.UserId == userID {
			grantID = grant.Id
			if slices.Contains(grant.RoleKeys, string(role)) {
				return nil
			}
			break
		}
	}

	if grantID == "" {
		addGrants := management.AddUserGrantRequest{
			UserId:    userID,
			ProjectId: conf.ProjectID,
			RoleKeys:  []string{string(role)},
		}
		_, err = zt.ManagementService().AddUserGrant(ctx, &addGrants)
		if err != nil {
			return fmt.Errorf("could not add user grant: %w", err)
		}
	} else {
		updateUserGrants := management.UpdateUserGrantRequest{
			UserId:   userID,
			RoleKeys: []string{string(role)},
			GrantId:  grantID,
		}
		_, err = zt.ManagementService().UpdateUserGrant(ctx, &updateUserGrants)
		if err != nil {
			return fmt.Errorf("could not update user grant: %w", err)
		}
	}
	return nil
}
