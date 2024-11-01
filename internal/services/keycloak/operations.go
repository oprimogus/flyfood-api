package keycloak

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/oprimogus/cardapiogo/internal/core/authentication"
	"github.com/oprimogus/cardapiogo/internal/core/user"
)

func entityUserToKeycloakUser(user user.User, userEnabled, userEmailVerified bool) *gocloak.User {
	realmRoles := make([]string, len(user.Roles))
	for i, v := range user.Roles {
		realmRoles[i] = string(v)
	}

	attributes := make(map[string][]string)
	if user.Profile.Document != "" {
		attributes["document"] = []string{user.Profile.Document}
	}
	if user.Profile.Phone != "" {
		attributes["phone"] = []string{user.Profile.Phone}
	}

	gocloakUser := gocloak.User{
		FirstName:     &user.Profile.Name,
		LastName:      &user.Profile.LastName,
		Username:      &user.Email,
		Enabled:       &userEnabled,
		EmailVerified: &userEmailVerified,
		Email:         &user.Email,
		RealmRoles:    &realmRoles,
		Attributes:    &attributes,
	}
	return &gocloakUser
}

func keycloakUserToEntityUser(userGocloak *gocloak.User) user.User {
	var document, phone string
	if userGocloak.Attributes != nil {
		if docValues, ok := (*userGocloak.Attributes)["document"]; ok && len(docValues) > 0 {
			document = docValues[0]
		}
		if phoneValues, ok := (*userGocloak.Attributes)["phone"]; ok && len(phoneValues) > 0 {
			phone = phoneValues[0]
		}
	}

	var userRoles []user.Role
	if userGocloak.RealmRoles != nil {
		userRoles = make([]user.Role, len(*userGocloak.RealmRoles))
		for i, v := range *userGocloak.RealmRoles {
			if user.IsValidRole(v) {
				userRoles[i] = user.Role(v)
			}
		}

	}

	return user.User{
		ID:    *userGocloak.ID,
		Email: *userGocloak.Email,
		Profile: user.Profile{
			Name:     *userGocloak.FirstName,
			LastName: *userGocloak.LastName,
			Document: document,
			Phone:    phone,
		},
		Roles: userRoles,
	}
}

func (k *Service) Create(ctx context.Context, user user.User) error {
	keycloakUser := entityUserToKeycloakUser(user, true, false)

	id, err := k.Client.CreateUser(ctx, k.Token.AccessToken, k.Realm, *keycloakUser)
	if err != nil {
		return err
	}
	errSetPassword := k.Client.SetPassword(ctx, k.Token.AccessToken, id, k.Realm, user.Password, false)
	if errSetPassword != nil {
		return errSetPassword
	}
	return nil
}

func (k *Service) FindByID(ctx context.Context, id string) (user.User, error) {
	userGocloak, err := k.Client.GetUserByID(ctx, k.Token.AccessToken, k.Realm, id)
	if err != nil {
		return user.User{}, err
	}
	return keycloakUserToEntityUser(userGocloak), nil
}

func (k *Service) FindByEmail(ctx context.Context, email string) (user.User, error) {
	maxUsers := 1
	userGocloak, err := k.Client.GetUsers(
		ctx,
		k.Token.AccessToken,
		k.Realm,
		gocloak.GetUsersParams{Email: &email, Max: &maxUsers},
	)
	if err != nil {
		return user.User{}, err
	}
	if len(userGocloak) == 0 {
		return user.User{}, nil
	}
	return keycloakUserToEntityUser(userGocloak[0]), nil
}

func (k *Service) Update(ctx context.Context, user user.User) error {
	actualUser, err := k.Client.GetUserByID(ctx, k.Token.AccessToken, k.Realm, user.ID)
	if err != nil {
		return fmt.Errorf("fail in found user: %w", err)
	}
	actualUser.FirstName = &user.Profile.Name
	actualUser.LastName = &user.Profile.LastName
	if actualUser.Attributes != nil {
		if phoneValues, ok := (*actualUser.Attributes)["phone"]; ok && len(phoneValues) > 0 {
			phoneValues[0] = user.Profile.Phone
		}
	}
	errUpdateUser := k.Client.UpdateUser(ctx, k.Token.AccessToken, k.Realm, *actualUser)
	if errUpdateUser != nil {
		return fmt.Errorf("fail in update user data: %w", errUpdateUser)
	}
	return nil
}

func (k *Service) Delete(ctx context.Context, id string) error {
	return k.Client.DeleteUser(ctx, k.Token.AccessToken, k.Realm, id)
}

func (k *Service) ResetPasswordByEmail(ctx context.Context, id string) error {
	lifespan := 30 * 60
	actions := []string{"UPDATE_PASSWORD"}
	return k.Client.ExecuteActionsEmail(
		ctx,
		k.Token.AccessToken,
		k.Realm,
		gocloak.ExecuteActionsEmail{
			Lifespan: &lifespan,
			ClientID: &k.ClientID,
			Actions:  &actions,
		},
	)
}

func (k *Service) AddRoles(ctx context.Context, userID string, roles []user.Role) error {
	userKeycloakRoles, err := k.Client.GetRealmRolesByUserID(ctx, k.Token.AccessToken, k.Realm, userID)
	if err != nil {
		return fmt.Errorf("fail in set new roles: %w", err)
	}

	allKeycloakRoles, err := k.Client.GetRealmRoles(ctx, k.Token.AccessToken, k.Realm, gocloak.GetRoleParams{})
	if err != nil {
		return fmt.Errorf("fail in get all roles of realm: %w", err)
	}

	userKeycloakRolesMap := make(map[string]bool)
	for _, v := range userKeycloakRoles {
		userKeycloakRolesMap[*v.Name] = true
	}

	allKeycloakRolesMap := make(map[string]*gocloak.Role)
	for _, v := range allKeycloakRoles {
		allKeycloakRolesMap[*v.Name] = v
	}

	var rolesToAdd []gocloak.Role
	for _, role := range roles {
		if !userKeycloakRolesMap[string(role)] {
			if keycloakRole, exists := allKeycloakRolesMap[string(role)]; exists {
				rolesToAdd = append(rolesToAdd, *keycloakRole)
			}
		}
	}

	if len(rolesToAdd) == 0 {
		return nil
	}

	errSetRoles := k.Client.AddRealmRoleToUser(ctx, k.Token.AccessToken, k.Realm, userID, rolesToAdd)
	if errSetRoles != nil {
		return fmt.Errorf("fail in set new roles for user: %w", errSetRoles)
	}
	return nil
}

func (k *Service) GetUsersWithUniqueParams(ctx context.Context, params user.GetUsersWithUniqueParams) (*[]user.User, error) {
	query := fmt.Sprintf("phone:%s", params.Phone)
	maxUsers := 1
	exact := true
	usersWithPhone, err := k.Client.GetUsers(ctx, k.Token.AccessToken, k.Realm, gocloak.GetUsersParams{
		Max: &maxUsers,
		Q:   &query,
	})
	if err != nil {
		return nil, err
	}

	usersWithEmail, err := k.Client.GetUsers(ctx, k.Token.AccessToken, k.Realm, gocloak.GetUsersParams{
		Max:   &maxUsers,
		Email: &params.Email,
		Exact: &exact,
	})
	if err != nil {
		return nil, err
	}

	allUsers := make([]gocloak.User, len(usersWithEmail)+len(usersWithPhone))
	copyIndex := 0

	for _, v := range usersWithEmail {
		allUsers[copyIndex] = *v
		copyIndex++
	}
	for _, v := range usersWithPhone {
		allUsers[copyIndex] = *v
		copyIndex++
	}

	userList := make([]user.User, len(allUsers))
	for i, v := range allUsers {
		userList[i] = keycloakUserToEntityUser(&v)
	}
	return &userList, nil
}

func (k *Service) SignIn(ctx context.Context, email, password string) (authentication.JWT, error) {
	jwtInstance, err := k.Client.Login(ctx, k.ClientID, k.ClientSecret, k.Realm, email, password)
	if err != nil {
		return authentication.JWT{}, err
	}
	return authentication.JWT{
		AccessToken:      jwtInstance.AccessToken,
		IDToken:          jwtInstance.IDToken,
		ExpiresIn:        jwtInstance.ExpiresIn,
		RefreshExpiresIn: jwtInstance.RefreshExpiresIn,
		RefreshToken:     jwtInstance.RefreshToken,
		TokenType:        jwtInstance.TokenType,
		NotBeforePolicy:  jwtInstance.NotBeforePolicy,
		SessionState:     jwtInstance.SessionState,
		Scope:            jwtInstance.Scope,
	}, nil
}

func (k *Service) RefreshToken(ctx context.Context, refreshToken string) (authentication.JWT, error) {
	jwtInstance, err := k.Client.RefreshToken(ctx, refreshToken, k.ClientID, k.ClientSecret, k.Realm)
	if err != nil {
		return authentication.JWT{}, err
	}
	return authentication.JWT{
		AccessToken:      jwtInstance.AccessToken,
		IDToken:          jwtInstance.IDToken,
		ExpiresIn:        jwtInstance.ExpiresIn,
		RefreshExpiresIn: jwtInstance.RefreshExpiresIn,
		RefreshToken:     jwtInstance.RefreshToken,
		TokenType:        jwtInstance.TokenType,
		NotBeforePolicy:  jwtInstance.NotBeforePolicy,
		SessionState:     jwtInstance.SessionState,
		Scope:            jwtInstance.Scope,
	}, nil
}

func (k *Service) IsValidToken(ctx context.Context, token string) (bool, error) {
	r, err := k.Client.RetrospectToken(ctx, token, k.ClientID, k.ClientSecret, k.Realm)
	if err != nil {
		return false, fmt.Errorf("ocurred an error while validate your access token: %w", err)
	}
	return *r.Active, nil
}

func (k *Service) DecodeAccessToken(ctx context.Context, accessToken string) (map[string]interface{}, error) {
	token, mapClaims, err := k.Client.DecodeAccessToken(ctx, accessToken, k.Realm)
	if err != nil {
		return nil, fmt.Errorf("unable to decode access token: %w", err)
	}
	if mapClaims == nil {
		return nil, fmt.Errorf("unable to decode access token and get claims")
	}
	if token == nil {
		return nil, fmt.Errorf("unable to decode access token and get metadata")
	}

	return *mapClaims, nil
}
