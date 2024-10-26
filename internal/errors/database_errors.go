package xerrors

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oprimogus/cardapiogo/internal/config"
)

const (
	NOT_FOUND_RECORD      = "Record not found"
	DUPLICATED_RECORD     = "There is a record with this data"
	FOREIGN_KEY_VIOLATION = "Foreign key violation"
	NULL_VIOLATION        = "Null value not allowed for column"
	VALUE_TOO_LONG        = "Input value too long for column"
	INTERNAL_SERVER_ERROR = "Internal Server Error"
	TOO_MANY_VALUES       = "There is more than one record"
	INVALID_VALUES        = "Invalid values for few fields"
	UNKNOWN_ERROR         = "Unknown error"
)

type fieldError struct {
	Field   string      `json:"field"`
	Input   string      `json:"input"`
	Message string      `json:"message"`
	Debug   interface{} `json:"debug,omitempty"`
}

func isDatabaseError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return true
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return true
	}
	if errors.Is(err, pgx.ErrTooManyRows) {
		return true
	}
	return false
}

func handleDatabaseErrors(err error, transactionID string) *CustomError {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return New(transactionID, http.StatusNotFound, NOT_FOUND_RECORD)
	}

	if errors.Is(err, pgx.ErrTooManyRows) {
		return New(transactionID, http.StatusInternalServerError, TOO_MANY_VALUES)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		slog.Error(pgErr.Error())
		switch pgErr.Code {
		case "23505":
			return handleUniqueViolation(transactionID, pgErr)
		case "23502", "22001", "22P02":
			return handleColumnViolation(transactionID, pgErr)
		default:
			if environment != string(config.Production) {
				return New(transactionID, http.StatusInternalServerError, UNKNOWN_ERROR, pgErr)
			}
			return New(transactionID, http.StatusInternalServerError, UNKNOWN_ERROR)
		}
	}
	return nil
}

func snakeToCamelCase(s string) string {
	words := strings.Split(s, "_")
	for i := 1; i < len(words); i++ {
		firstChar := []rune(words[i])[0]
		words[i] = string(unicode.ToUpper(firstChar)) + words[i][1:]
	}
	return strings.Join(words, "")

}

func handleUniqueViolation(transactionID string, pgErr *pgconn.PgError) *CustomError {
	startField := strings.Index(pgErr.Detail, "(") + 1
	endField := strings.Index(pgErr.Detail, ")=")
	field := snakeToCamelCase(pgErr.Detail[startField:endField])

	startValue := strings.Index(pgErr.Detail, "=(") + 2
	endValue := strings.LastIndex(pgErr.Detail, ")")
	value := pgErr.Detail[startValue:endValue]

	description := "This value is already in use"

	fieldError := fieldError{
		Field:   field,
		Input:   value,
		Message: description,
	}
	if environment != string(config.Production) {
		fieldError.Debug = pgErr
	}
	return New(transactionID, http.StatusConflict, DUPLICATED_RECORD, fieldError)
}

func handleColumnViolation(transactionID string, pgErr *pgconn.PgError) *CustomError {
	fieldError := fieldError{
		Field:   "",
		Input:   "",
		Message: pgErr.Message,
	}
	if environment != string(config.Production) {
		fieldError.Debug = pgErr
	}
	return New(transactionID, http.StatusBadRequest, INVALID_VALUES, fieldError)
}
