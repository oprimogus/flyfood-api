package xerrors_test

import (
	// "fmt"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	// "github.com/Nerzal/gocloak/v13"
	// "github.com/jackc/pgx/v5"
	// "github.com/jackc/pgx/v5/pgconn"
	// "github.com/oprimogus/cardapiogo/internal/config"
	// "github.com/oprimogus/cardapiogo/internal/core/store"
	// "github.com/oprimogus/cardapiogo/internal/core/user"
	"github.com/Nerzal/gocloak/v13"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oprimogus/cardapiogo/internal/config"
	"github.com/oprimogus/cardapiogo/internal/core/store"
	"github.com/oprimogus/cardapiogo/internal/core/user"
	xerrors "github.com/oprimogus/cardapiogo/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type XerrorSuite struct {
	suite.Suite
}

func (s *XerrorSuite) SetupTest() {
	config.GetInstance().Api.Environment = string(config.Production)
}

func TestUnitXerrorSuite(t *testing.T) {
	suite.Run(t, new(XerrorSuite))
}

func (s *XerrorSuite) TestNew() {
	tests := []struct {
		name    string
		TraceID string
		status  int
		message string
		details []interface{}
		want    *xerrors.CustomError
	}{
		{
			name:    "should create error without details",
			TraceID: "123",
			status:  http.StatusBadRequest,
			message: "test error",
			details: nil,
			want: &xerrors.CustomError{
				Status:       http.StatusBadRequest,
				ErrorMessage: "test error",
				TraceID:      "123",
			},
		},
		{
			name:    "should create error with details",
			TraceID: "456",
			status:  http.StatusNotFound,
			message: "not found",
			details: []interface{}{"additional info"},
			want: &xerrors.CustomError{
				Status:       http.StatusNotFound,
				ErrorMessage: "not found",
				TraceID:      "456",
				Details:      "additional info",
			},
		},
	}

	for _, tt := range tests {
		var details []interface{}
		if tt.details != nil {
			details = tt.details
		}
		got := xerrors.New(tt.TraceID, tt.status, tt.message, details...)
		assert.Equal(s.T(), tt.want, got, tt.name)
	}
}

func (s *XerrorSuite) TestHandleDatabaseError() {
	tests := []struct {
		name        string
		err         error
		TraceID     string
		want        *xerrors.CustomError
		environment string
	}{
		{
			name:    "should handle unique violation",
			err:     &pgconn.PgError{Code: "23505", Detail: "(email)=(test@test.com)"},
			TraceID: "123",
			want: xerrors.New("123", http.StatusConflict, xerrors.DuplicatedRecord, xerrors.FieldError{
				Field:   "email",
				Input:   "test@test.com",
				Message: "This value is already in use",
			}),
			environment: string(config.Production),
		},
		{
			name:    "should handle null violation",
			err:     &pgconn.PgError{Code: "23502", Message: "null value in column"},
			TraceID: "123",
			want: xerrors.New("123", http.StatusBadRequest, xerrors.InvalidValues, xerrors.FieldError{
				Message: "null value in column",
			}),
			environment: string(config.Production),
		},
		{
			name:        "should handle pgx.ErrNoRows",
			err:         pgx.ErrNoRows,
			TraceID:     "123",
			want:        xerrors.New("123", http.StatusNotFound, xerrors.NotFoundRecord),
			environment: string(config.Production),
		},
		{
			name:        "should handle pgx.ErrTooManyRows",
			err:         pgx.ErrTooManyRows,
			TraceID:     "123",
			want:        xerrors.New("123", http.StatusInternalServerError, xerrors.TooManyValues),
			environment: string(config.Production),
		},
	}

	for _, tt := range tests {
		config.GetInstance().Api.Environment = tt.environment
		got := xerrors.HandleError(tt.err, tt.TraceID)
		assert.Equal(s.T(), tt.want, got, tt.name)
	}
}

func (s *XerrorSuite) TestHandleCoreError() {
	tests := []struct {
		name    string
		err     error
		TraceID string
		want    *xerrors.CustomError
	}{
		{
			name:    "should handle user phone conflict",
			err:     user.ErrExistUserWithPhone,
			TraceID: "123",
			want: &xerrors.CustomError{
				Status:       http.StatusConflict,
				ErrorMessage: user.ErrExistUserWithPhone.Error(),
				TraceID:      "123",
			},
		},
		{
			name:    "should handle user email conflict",
			err:     user.ErrExistUserWithEmail,
			TraceID: "123",
			want: &xerrors.CustomError{
				Status:       http.StatusConflict,
				ErrorMessage: user.ErrExistUserWithEmail.Error(),
				TraceID:      "123",
			},
		},
		{
			name:    "should handle store closing time error",
			err:     store.ErrClosingTimeBeforeOpeningTime,
			TraceID: "123",
			want: &xerrors.CustomError{
				Status:       http.StatusBadRequest,
				ErrorMessage: store.ErrClosingTimeBeforeOpeningTime.Error(),
				TraceID:      "123",
			},
		},
		{
			name:    "should handle store opening time error",
			err:     store.ErrOpeningTimeAfterClosingTime,
			TraceID: "123",
			want: &xerrors.CustomError{
				Status:       http.StatusBadRequest,
				ErrorMessage: store.ErrOpeningTimeAfterClosingTime.Error(),
				TraceID:      "123",
			},
		},
		{
			name:    "should handle store not resource owner error",
			err:     store.ErrNotOwner,
			TraceID: "123",
			want: &xerrors.CustomError{
				Status:       http.StatusForbidden,
				ErrorMessage: store.ErrNotOwner.Error(),
				TraceID:      "123",
			},
		},
	}

	for _, tt := range tests {
		got := xerrors.HandleError(tt.err, tt.TraceID)
		assert.Equal(s.T(), tt.want, got, tt.name)
	}
}

func (s *XerrorSuite) TestInvalidInput() {
	tests := []struct {
		name    string
		TraceID string
		errs    map[string]string
		want    *xerrors.CustomError
	}{
		{
			name:    "should create invalid input error",
			TraceID: "123",
			errs: map[string]string{
				"email": "invalid email format",
				"name":  "name is required",
			},
			want: &xerrors.CustomError{
				Status:       http.StatusBadRequest,
				ErrorMessage: xerrors.InvalidInputMessage,
				TraceID:      "123",
				Details: []xerrors.InvalidField{
					{Field: "email", Error: "invalid email format"},
					{Field: "name", Error: "name is required"},
				},
			},
		},
	}

	for _, tt := range tests {
		got := xerrors.InvalidInput(tt.TraceID, tt.errs)
		assert.Equal(s.T(), tt.want.Status, got.Status, tt.name)
		assert.Equal(s.T(), tt.want.ErrorMessage, got.ErrorMessage, tt.name)
		assert.Equal(s.T(), tt.want.TraceID, got.TraceID, tt.name)

		gotDetails, ok := got.Details.([]xerrors.InvalidField)
		assert.True(s.T(), ok, tt.name)
		wantDetails := tt.want.Details.([]xerrors.InvalidField)
		assert.ElementsMatch(s.T(), wantDetails, gotDetails, tt.name)
	}
}

func (s *XerrorSuite) TestHandleJSONError() {
	tests := []struct {
		name    string
		err     error
		TraceID string
		want    *xerrors.CustomError
	}{
		{
			name: "should handle SyntaxError",
			err: &json.SyntaxError{
				Offset: 42,
			},
			TraceID: "123",
			want: &xerrors.CustomError{
				Status:       http.StatusBadRequest,
				ErrorMessage: fmt.Sprintf("Invalid JSON: %s", &json.SyntaxError{Offset: 42}),
				Details:      int64(42),
				TraceID:      "123",
			},
		},
	}

	for _, tt := range tests {
		got := xerrors.HandleError(tt.err, tt.TraceID)
		assert.Equal(s.T(), tt.want, got, tt.name)
	}
}

func (s *XerrorSuite) TestHandleGocloakError() {
	tests := []struct {
		name        string
		err         error
		TraceID     string
		want        *xerrors.CustomError
		environment string
	}{
		{
			name: "should handle single message in production",
			err: &gocloak.APIError{
				Code:    http.StatusUnauthorized,
				Message: "Invalid credentials",
			},
			TraceID: "123",
			want: &xerrors.CustomError{
				Status:       http.StatusUnauthorized,
				ErrorMessage: "An authentication error occurred while processing your request.",
				TraceID:      "123",
			},
			environment: string(config.Production),
		},
		{
			name: "should handle single message in development",
			err: &gocloak.APIError{
				Code:    http.StatusUnauthorized,
				Message: "Invalid credentials",
			},
			TraceID: "123",
			want: &xerrors.CustomError{
				Status:       http.StatusUnauthorized,
				ErrorMessage: "Invalid credentials",
				TraceID:      "123",
				Debug: &gocloak.APIError{
					Code:    http.StatusUnauthorized,
					Message: "Invalid credentials",
				},
			},
			environment: string(config.Staging),
		},
		{
			name: "should handle split message in development",
			err: &gocloak.APIError{
				Code:    http.StatusBadRequest,
				Message: "validation error: invalid token format",
			},
			TraceID: "123",
			want: &xerrors.CustomError{
				Status:       http.StatusBadRequest,
				ErrorMessage: "invalid token format",
				TraceID:      "123",
				Debug:        &gocloak.APIError{Code: http.StatusBadRequest, Message: "validation error: invalid token format"},
			},
			environment: string(config.Staging),
		},
	}

	for _, tt := range tests {
		config.GetInstance().Api.Environment = tt.environment
		got := xerrors.HandleError(tt.err, tt.TraceID)
		assert.Equal(s.T(), tt.want, got, tt.name)
	}
}
