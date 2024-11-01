package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oprimogus/cardapiogo/internal/api/controller"
	validatorutils "github.com/oprimogus/cardapiogo/internal/api/validator"
	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/oprimogus/cardapiogo/internal/core/authentication"
	"github.com/oprimogus/cardapiogo/internal/core/user"
	postgresDB "github.com/oprimogus/cardapiogo/internal/database/postgres"
	"github.com/oprimogus/cardapiogo/internal/persistence"
	"github.com/oprimogus/cardapiogo/internal/services/adapter"
	"github.com/oprimogus/cardapiogo/test/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	signInEndpoint = "/v1/auth/sign-in"
	signUpEndpoint = "/v1/auth/sign-up"
)

type AuthControllerSuite struct {
	suite.Suite
	ctx               context.Context
	authController    *controller.AuthController
	userController    *controller.UserController
	pgContainer       *integration.Container
	keycloakContainer *integration.Container
}

func createUserMocked(router *gin.Engine) error {
	input := map[string]interface{}{
		"email":    "johndoe@example.com",
		"password": "teste123",
		"profile": map[string]interface{}{
			"lastName": "Doe",
			"name":     "John Mock",
			"phone":    "+5513995590865",
		},
	}

	requestBody, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", signUpEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	fmt.Println(w.Body.String())

	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		return fmt.Errorf("failed to create user, status code: %d, body: %s",
			w.Code, w.Body.String())
	}
	return nil
}

func (s *AuthControllerSuite) SetupSuite() {
	ctx := context.Background()
	s.ctx = ctx
	postgres, err := integration.MakePostgres(s.ctx)
	if err != nil {
		panic(err)
	}
	s.pgContainer = postgres

	keycloak, err := integration.MakeKeycloak(s.ctx)
	if err != nil {
		panic(err)
	}
	s.keycloakContainer = keycloak

	validator, err := validatorutils.NewValidator("pt")
	if err != nil && validator == nil {
		s.T().Fatal("fail on make validator: %w", err)
	}
	db := postgresDB.GetInstance()
	config.GetInstance().Api.Environment = string(config.Production)

	serviceFactory := adapter.NewServiceFactory()
	factory := persistence.NewDataBaseRepositoryFactory(db, serviceFactory)
	s.authController = controller.NewAuthController(validator, factory.NewAuthenticationRepository())
	s.userController = controller.NewUserController(validator, factory.NewUserRepository())
}

func (s *AuthControllerSuite) TearDownSuite() {
	s.keycloakContainer.Kill(s.ctx)
	s.pgContainer.Kill(s.ctx)
}

func TestIntegrationAuthControllerSuite(t *testing.T) {
	suite.Run(t, new(AuthControllerSuite))
}

func (s *AuthControllerSuite) TestSignUp() {
	tests := []struct {
		name               string
		requestBody        map[string]interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "Should create a new user with sucess",
			requestBody: map[string]interface{}{
				"email":    "newuser@example.com",
				"password": "teste123",
				"profile": map[string]interface{}{
					"lastName": "Doe",
					"name":     "John 1",
					"phone":    "+5513995590868",
				},
			},
			expectedStatusCode: 201,
			expectedResponse:   nil,
		},
		{
			name: "Should return error when exist user with same email",
			requestBody: map[string]interface{}{
				"email":    "johndoe@example.com",
				"password": "teste123",
				"profile": map[string]interface{}{
					"lastName": "Doe",
					"name":     "John 2",
					"phone":    "+5513995590865",
				},
			},
			expectedStatusCode: 409,
			expectedResponse: map[string]interface{}{
				"error":   user.ErrExistUserWithEmail.Error(),
				"traceID": "",
			},
		},
		{
			name: "Should return error when exist user with same phone",
			requestBody: map[string]interface{}{
				"email":    "johndoe2@example.com",
				"password": "teste123",
				"profile": map[string]interface{}{
					"lastName": "Doe",
					"name":     "John 3",
					"phone":    "+5513995590865",
				},
			},
			expectedStatusCode: 409,
			expectedResponse: map[string]interface{}{
				"error":   user.ErrExistUserWithPhone.Error(),
				"traceID": "",
			},
		},
	}

	for _, test := range tests {
		router := gin.New()
		gin.SetMode(gin.TestMode)
		router.POST(signUpEndpoint, s.userController.CreateUser)
		createUserMocked(router)

		requestBody, _ := json.Marshal(test.requestBody)

		req, _ := http.NewRequest("POST", signUpEndpoint, bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(s.T(), test.expectedStatusCode, w.Code, test.name)

		if test.expectedResponse == nil {
			assert.Equal(s.T(), 0, len(w.Body.Bytes()), test.name)
		} else {
			var responseBody interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(s.T(), err, test.name)
			assert.Equal(s.T(), test.expectedResponse, responseBody, test.name)
		}
	}
}

func (s *AuthControllerSuite) TestSignIn() {
	tests := []struct {
		name               string
		requestBody        map[string]interface{}
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "Should sign-in with success",

			requestBody: map[string]interface{}{
				"email":    "johndoe@example.com",
				"password": "teste123",
			},
			expectedStatusCode: 200,
			expectedResponse:   authentication.JWT{},
		},
		{
			name: "Should return credentials error",
			requestBody: map[string]interface{}{
				"email":    "johndoe@example.com",
				"password": "teste1234",
			},
			expectedStatusCode: 401,
			expectedResponse: map[string]interface{}{
				"error":   "Invalid user credentials",
				"details": "invalid_grant",
				"traceID": "",
			},
		},
	}

	for _, test := range tests {
		router := gin.New()
		gin.SetMode(gin.TestMode)
		router.POST(signUpEndpoint, s.userController.CreateUser)
		router.POST(signInEndpoint, s.authController.SignIn)
		createUserMocked(router)

		requestBody, _ := json.Marshal(test.requestBody)

		req, _ := http.NewRequest("POST", signInEndpoint, bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(s.T(), test.expectedStatusCode, w.Code, test.name)

		var responseBody map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(s.T(), err)

		if w.Code == http.StatusOK {
			assert.Contains(s.T(), responseBody, "accessToken", test.name)
			assert.Contains(s.T(), responseBody, "idToken", test.name)
			assert.Contains(s.T(), responseBody, "expiresIn", test.name)
			assert.Contains(s.T(), responseBody, "refreshToken", test.name)
			assert.Contains(s.T(), responseBody, "tokenType", test.name)
			assert.Contains(s.T(), responseBody, "scope", test.name)

			assert.IsType(s.T(), "", responseBody["accessToken"], test.name)
			assert.IsType(s.T(), "", responseBody["idToken"], test.name)
			assert.IsType(s.T(), float64(0), responseBody["expiresIn"], test.name)
			assert.IsType(s.T(), "", responseBody["refreshToken"], test.name)
			assert.IsType(s.T(), "", responseBody["tokenType"], test.name)
			assert.IsType(s.T(), "", responseBody["scope"], test.name)
		} else {
			assert.Equal(s.T(), test.expectedResponse, responseBody, test.name)
		}

	}
}
