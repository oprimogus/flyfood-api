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

func (x *XerrorSuite) TestNew() {
	tests := []struct {
		name          string
		transactionID string
		status        int
		message       string
		details       []interface{}
		want          *xerrors.CustomError
	}{
		{
			name:          "should create error without details",
			transactionID: "123",
			status:        http.StatusBadRequest,
			message:       "test error",
			details:       nil,
			want: &xerrors.CustomError{
				Status:        http.StatusBadRequest,
				ErrorMessage:  "test error",
				TransactionID: "123",
			},
		},
		{
			name:          "should create error with details",
			transactionID: "456",
			status:        http.StatusNotFound,
			message:       "not found",
			details:       []interface{}{"additional info"},
			want: &xerrors.CustomError{
				Status:        http.StatusNotFound,
				ErrorMessage:  "not found",
				TransactionID: "456",
				Details:       "additional info",
			},
		},
	}

	for _, tt := range tests {
		var details []interface{}
		if tt.details != nil {
			details = tt.details
		}
		got := xerrors.New(tt.transactionID, tt.status, tt.message, details...)
		assert.Equal(x.T(), tt.want, got, tt.name)
	}
}

func (x *XerrorSuite) TestHandleDatabaseError() {
	tests := []struct {
		name          string
		err           error
		transactionID string
		want          *xerrors.CustomError
		environment   string
	}{
		{
			name:          "should handle unique violation",
			err:           &pgconn.PgError{Code: "23505", Detail: "(email)=(test@test.com)"},
			transactionID: "123",
			want: xerrors.New("123", http.StatusConflict, xerrors.DUPLICATED_RECORD, xerrors.FieldError{
				Field:   "email",
				Input:   "test@test.com",
				Message: "This value is already in use",
			}),
			environment: string(config.Production),
		},
		{
			name:          "should handle null violation",
			err:           &pgconn.PgError{Code: "23502", Message: "null value in column"},
			transactionID: "123",
			want: xerrors.New("123", http.StatusBadRequest, xerrors.INVALID_VALUES, xerrors.FieldError{
				Message: "null value in column",
			}),
			environment: string(config.Production),
		},
		{
			name:          "should handle pgx.ErrNoRows",
			err:           pgx.ErrNoRows,
			transactionID: "123",
			want:          xerrors.New("123", http.StatusNotFound, xerrors.NOT_FOUND_RECORD),
			environment:   string(config.Production),
		},
		{
			name:          "should handle pgx.ErrTooManyRows",
			err:           pgx.ErrTooManyRows,
			transactionID: "123",
			want:          xerrors.New("123", http.StatusInternalServerError, xerrors.TOO_MANY_VALUES),
			environment:   string(config.Production),
		},
	}

	for _, tt := range tests {
		config.GetInstance().Api.Environment = tt.environment
		got := xerrors.HandleError(tt.err, tt.transactionID)
		assert.Equal(x.T(), tt.want, got, tt.name)
	}
}

func (x *XerrorSuite) TestHandleCoreError() {
	tests := []struct {
		name          string
		err           error
		transactionID string
		want          *xerrors.CustomError
	}{
		{
			name:          "should handle user phone conflict",
			err:           user.ErrExistUserWithPhone,
			transactionID: "123",
			want: &xerrors.CustomError{
				Status:        http.StatusConflict,
				ErrorMessage:  user.ErrExistUserWithPhone.Error(),
				TransactionID: "123",
			},
		},
		{
			name:          "should handle user email conflict",
			err:           user.ErrExistUserWithEmail,
			transactionID: "123",
			want: &xerrors.CustomError{
				Status:        http.StatusConflict,
				ErrorMessage:  user.ErrExistUserWithEmail.Error(),
				TransactionID: "123",
			},
		},
		{
			name:          "should handle store closing time error",
			err:           store.ErrClosingTimeBeforeOpeningTime,
			transactionID: "123",
			want: &xerrors.CustomError{
				Status:        http.StatusBadRequest,
				ErrorMessage:  store.ErrClosingTimeBeforeOpeningTime.Error(),
				TransactionID: "123",
			},
		},
		{
			name:          "should handle store opening time error",
			err:           store.ErrOpeningTimeAfterClosingTime,
			transactionID: "123",
			want: &xerrors.CustomError{
				Status:        http.StatusBadRequest,
				ErrorMessage:  store.ErrOpeningTimeAfterClosingTime.Error(),
				TransactionID: "123",
			},
		},
		{
			name:          "should handle store not resource owner error",
			err:           store.ErrNotOwner,
			transactionID: "123",
			want: &xerrors.CustomError{
				Status:        http.StatusForbidden,
				ErrorMessage:  store.ErrNotOwner.Error(),
				TransactionID: "123",
			},
		},
	}

	for _, tt := range tests {
		got := xerrors.HandleError(tt.err, tt.transactionID)
		assert.Equal(x.T(), tt.want, got, tt.name)
	}
}

func (x *XerrorSuite) TestInvalidInput() {
	tests := []struct {
		name          string
		transactionID string
		errs          map[string]string
		want          *xerrors.CustomError
	}{
		{
			name:          "should create invalid input error",
			transactionID: "123",
			errs: map[string]string{
				"email": "invalid email format",
				"name":  "name is required",
			},
			want: &xerrors.CustomError{
				Status:        http.StatusBadRequest,
				ErrorMessage:  xerrors.INVALID_INPUT_MESSAGE,
				TransactionID: "123",
				Details: []xerrors.InvalidField{
					{Field: "email", Error: "invalid email format"},
					{Field: "name", Error: "name is required"},
				},
			},
		},
	}

	for _, tt := range tests {
		got := xerrors.InvalidInput(tt.transactionID, tt.errs)
		assert.Equal(x.T(), tt.want.Status, got.Status, tt.name)
		assert.Equal(x.T(), tt.want.ErrorMessage, got.ErrorMessage, tt.name)
		assert.Equal(x.T(), tt.want.TransactionID, got.TransactionID, tt.name)

		gotDetails, ok := got.Details.([]xerrors.InvalidField)
		assert.True(x.T(), ok, tt.name)
		wantDetails := tt.want.Details.([]xerrors.InvalidField)
		assert.ElementsMatch(x.T(), wantDetails, gotDetails, tt.name)
	}
}

func (x *XerrorSuite) TestHandleJSONError() {
	tests := []struct {
		name          string
		err           error
		transactionID string
		want          *xerrors.CustomError
	}{
		{
			name: "should handle SyntaxError",
			err: &json.SyntaxError{
				Offset: 42,
			},
			transactionID: "123",
			want: &xerrors.CustomError{
				Status:        http.StatusBadRequest,
				ErrorMessage:  fmt.Sprintf("Invalid JSON: %s", &json.SyntaxError{Offset: 42}),
				Details:       int64(42),
				TransactionID: "123",
			},
		},
	}

	for _, tt := range tests {
		got := xerrors.HandleError(tt.err, tt.transactionID)
		assert.Equal(x.T(), tt.want, got, tt.name)
	}
}

func (x *XerrorSuite) TestHandleGocloakError() {
	tests := []struct {
		name           string
		err            error
		transactionID  string
		want           *xerrors.CustomError
		environment    string
	}{
		{
			name: "should handle single message in production",
			err: &gocloak.APIError{
				Code:    http.StatusUnauthorized,
				Message: "Invalid credentials",
			},
			transactionID: "123",
			want: &xerrors.CustomError{
				Status:        http.StatusUnauthorized,
				ErrorMessage:  "An authentication error occurred while processing your request.",
				TransactionID: "123",
			},
			environment: string(config.Production),
		},
		{
			name: "should handle single message in development",
			err: &gocloak.APIError{
				Code:    http.StatusUnauthorized,
				Message: "Invalid credentials",
			},
			transactionID: "123",
			want: &xerrors.CustomError{
				Status:        http.StatusUnauthorized,
				ErrorMessage:  "Invalid credentials",
				TransactionID: "123",
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
			transactionID: "123",
			want: &xerrors.CustomError{
				Status:        http.StatusBadRequest,
				ErrorMessage:  "invalid token format",
				TransactionID: "123",
				Debug:         &gocloak.APIError{Code: http.StatusBadRequest, Message: "validation error: invalid token format"},
			},
			environment: string(config.Staging),
		},
	}

	for _, tt := range tests {
		config.GetInstance().Api.Environment = tt.environment
		got := xerrors.HandleError(tt.err, tt.transactionID)
		assert.Equal(x.T(), tt.want, got, tt.name)
	}
}
