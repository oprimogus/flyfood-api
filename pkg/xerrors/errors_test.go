package xerrors

//import (
//	"net/http"
//	"testing"
//
//	"github.com/oprimogus/flyfood-api/internal/config"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/suite"
//)
//
//type XerrorSuite struct {
//	suite.Suite
//}
//
//func (s *XerrorSuite) SetupTest() {
//	config.GetInstance().Api.Environment = string(config.Production)
//}
//
//func TestUnitXerrorSuite(t *testing.T) {
//	suite.Run(t, new(XerrorSuite))
//}
//
//func (s *XerrorSuite) TestNew() {
//	tests := []struct {
//		name    string
//		TraceID string
//		status  int
//		message string
//		details []any
//		want    *CustomError
//	}{
//		{
//			name:    "should create error without details",
//			TraceID: "123",
//			status:  http.StatusBadRequest,
//			message: "test error",
//			details: nil,
//			want: &CustomError{
//				Status:       http.StatusBadRequest,
//				ErrorMessage: "test error",
//				TraceID:      "123",
//			},
//		},
//		{
//			name:    "should create error with details",
//			TraceID: "456",
//			status:  http.StatusNotFound,
//			message: "not found",
//			details: []any{"additional info"},
//			want: &CustomError{
//				Status:       http.StatusNotFound,
//				ErrorMessage: "not found",
//				TraceID:      "456",
//				Details:      "additional info",
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		var details []any
//		if tt.details != nil {
//			details = tt.details
//		}
//		got := New(tt.TraceID, tt.status, tt.message, details...)
//		assert.Equal(s.T(), tt.want, got, tt.name)
//	}
//}

//func (s *XerrorSuite) TestHandleDatabaseError() {
//	tests := []struct {
//		name        string
//		err         error
//		TraceID     string
//		want        *CustomError
//		environment string
//	}{
//		{
//			name:    "should handle unique violation",
//			err:     &pgconn.PgError{Code: "23505", Detail: "(email)=(test@test.com)"},
//			TraceID: "123",
//			want: New("123", http.StatusConflict, ErrDuplicatedRecord, FieldError{
//				Field:   "email",
//				Input:   "test@test.com",
//				Message: "This value is already in use",
//			}),
//			environment: string(config.Production),
//		},
//		{
//			name:    "should handle null violation",
//			err:     &pgconn.PgError{Code: "23502", Message: "null value in column"},
//			TraceID: "123",
//			want: New("123", http.StatusBadRequest, ErrInvalidValues, FieldError{
//				Message: "null value in column",
//			}),
//			environment: string(config.Production),
//		},
//		{
//			name:        "should handle pgx.ErrNoRows",
//			err:         pgx.ErrNoRows,
//			TraceID:     "123",
//			want:        New("123", http.StatusNotFound, ),
//			environment: string(config.Production),
//		},
//		{
//			name:        "should handle pgx.ErrTooManyRows",
//			err:         pgx.ErrTooManyRows,
//			TraceID:     "123",
//			want:        New("123", http.StatusInternalServerError, ErrTooManyValues),
//			environment: string(config.Production),
//		},
//	}
//
//	for _, tt := range tests {
//		config.GetInstance().Api.Environment = tt.environment
//		got := HandleError(tt.err, tt.TraceID)
//		assert.Equal(s.T(), tt.want, got, tt.name)
//	}
//}
//
//func (s *XerrorSuite) TestInvalidInput() {
//	tests := []struct {
//		name    string
//		TraceID string
//		errs    map[string]string
//		want    *CustomError
//	}{
//		{
//			name:    "should create invalid input error",
//			TraceID: "123",
//			errs: map[string]string{
//				"email": "invalid email format",
//				"name":  "name is required",
//			},
//			want: &CustomError{
//				Status:       http.StatusBadRequest,
//				ErrorMessage: InvalidInputMessage,
//				TraceID:      "123",
//				Details: []InvalidField{
//					{Field: "email", Error: "invalid email format"},
//					{Field: "name", Error: "name is required"},
//				},
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		got := InvalidInput(tt.TraceID, tt.errs)
//		assert.Equal(s.T(), tt.want.Status, got.Status, tt.name)
//		assert.Equal(s.T(), tt.want.ErrorMessage, got.ErrorMessage, tt.name)
//		assert.Equal(s.T(), tt.want.TraceID, got.TraceID, tt.name)
//
//		gotDetails, ok := got.Details.([]InvalidField)
//		assert.True(s.T(), ok, tt.name)
//		wantDetails := tt.want.Details.([]InvalidField)
//		assert.ElementsMatch(s.T(), wantDetails, gotDetails, tt.name)
//	}
//}

//func (s *XerrorSuite) TestHandleJSONError() {
//	tests := []struct {
//		name    string
//		err     error
//		TraceID string
//		want    *CustomError
//	}{
//		{
//			name: "should handle SyntaxError",
//			err: &json.SyntaxError{
//				Offset: 42,
//			},
//			TraceID: "123",
//			want: &CustomError{
//				Status:       http.StatusBadRequest,
//				ErrorMessage: fmt.Sprintf("Invalid JSON: %s", &json.SyntaxError{Offset: 42}),
//				Details:      int64(42),
//				TraceID:      "123",
//			},
//		},
//	}
//
//	for _, tt := range tests {
//		got := HandleError(tt.err, tt.TraceID)
//		assert.Equal(s.T(), tt.want, got, tt.name)
//	}
//}
//
//func (s *XerrorSuite) TestHandleGocloakError() {
//	tests := []struct {
//		name        string
//		err         error
//		TraceID     string
//		want        *CustomError
//		environment string
//	}{
//		{
//			name: "should handle single message in production",
//			err: &gocloak.APIError{
//				Code:    http.StatusUnauthorized,
//				Message: "Invalid credentials",
//			},
//			TraceID: "123",
//			want: &CustomError{
//				Status:       http.StatusUnauthorized,
//				ErrorMessage: "An authentication error occurred while processing your request.",
//				TraceID:      "123",
//			},
//			environment: string(config.Production),
//		},
//		{
//			name: "should handle single message in development",
//			err: &gocloak.APIError{
//				Code:    http.StatusUnauthorized,
//				Message: "Invalid credentials",
//			},
//			TraceID: "123",
//			want: &CustomError{
//				Status:       http.StatusUnauthorized,
//				ErrorMessage: "Invalid credentials",
//				TraceID:      "123",
//				Debug: &gocloak.APIError{
//					Code:    http.StatusUnauthorized,
//					Message: "Invalid credentials",
//				},
//			},
//			environment: string(config.Staging),
//		},
//		{
//			name: "should handle split message in development",
//			err: &gocloak.APIError{
//				Code:    http.StatusBadRequest,
//				Message: "validation error: invalid token format",
//			},
//			TraceID: "123",
//			want: &CustomError{
//				Status:       http.StatusBadRequest,
//				ErrorMessage: "invalid token format",
//				TraceID:      "123",
//				Debug:        &gocloak.APIError{Code: http.StatusBadRequest, Message: "validation error: invalid token format"},
//			},
//			environment: string(config.Staging),
//		},
//	}
//
//	for _, tt := range tests {
//		config.GetInstance().Api.Environment = tt.environment
//		got := HandleError(tt.err, tt.TraceID)
//		assert.Equal(s.T(), tt.want, got, tt.name)
//	}
//}
