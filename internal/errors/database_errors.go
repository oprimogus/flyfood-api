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

func handleDatabaseError(err error, transactionID string) *CustomError {
	if !isDatabaseError(err) {
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
			if config.GetInstance().Api.Environment != string(config.Production) {
				return New(transactionID, http.StatusInternalServerError, UNKNOWN_ERROR, pgErr)
			}
			return New(transactionID, http.StatusInternalServerError, UNKNOWN_ERROR)
		}
	}

	return nil
}

// Utils
func isDatabaseError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) ||
		errors.Is(err, pgx.ErrNoRows) ||
		errors.Is(err, pgx.ErrTooManyRows)
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

	fieldErr := FieldError{
		Field:   field,
		Input:   value,
		Message: "This value is already in use",
	}

	if config.GetInstance().Api.Environment != string(config.Production) {
		fieldErr.Debug = pgErr
	}

	return New(transactionID, http.StatusConflict, DUPLICATED_RECORD, fieldErr)
}

func handleColumnViolation(transactionID string, pgErr *pgconn.PgError) *CustomError {
	fieldErr := FieldError{
		Field:   "",
		Input:   "",
		Message: pgErr.Message,
	}

	if config.GetInstance().Api.Environment != string(config.Production) {
		fieldErr.Debug = pgErr
	}

	return New(transactionID, http.StatusBadRequest, INVALID_VALUES, fieldErr)
}
