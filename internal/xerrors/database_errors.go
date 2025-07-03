package xerrors

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oprimogus/flyfood-api/internal/config"
)

var (
	NotFoundRecord   = errors.New("Registro não encontrado")
	DuplicatedRecord = errors.New("Já existe um registro com os dados fornecidos")
	TooManyValues    = errors.New("Existe mais de um registro com o mesmo identificador")
	InvalidValues    = errors.New("Alguns campos estão com valores inválidos")
	UnknownError     = errors.New("Erro interno desconhecido")
)

func handleDatabaseError(ctx context.Context, err error) *CustomError {
	if !isDatabaseError(err) {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return New(ctx, http.StatusNotFound, NotFoundRecord)
	}

	if errors.Is(err, pgx.ErrTooManyRows) {
		return New(ctx, http.StatusInternalServerError, TooManyValues)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		slog.Error(pgErr.Error())
		switch pgErr.Code {
		case "23505":
			return handleUniqueViolation(ctx, pgErr)
		case "23502", "22001", "22P02":
			return handleColumnViolation(ctx, pgErr)
		default:
			if config.GetInstance().Api.Environment != string(config.Production) {
				return New(ctx, http.StatusInternalServerError, UnknownError, pgErr)
			}
			return New(ctx, http.StatusInternalServerError, UnknownError)
		}
	}

	return nil
}

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

func handleUniqueViolation(ctx context.Context, pgErr *pgconn.PgError) *CustomError {
	startField := strings.Index(pgErr.Detail, "(") + 1
	endField := strings.Index(pgErr.Detail, ")=")
	field := snakeToCamelCase(pgErr.Detail[startField:endField])

	startValue := strings.Index(pgErr.Detail, "=(") + 2
	endValue := strings.LastIndex(pgErr.Detail, ")")
	value := pgErr.Detail[startValue:endValue]

	fieldErr := FieldError{
		Field:   field,
		Input:   value,
		Message: "Existe um registro com este valor",
	}

	if config.GetInstance().Api.Environment != string(config.Production) {
		fieldErr.Debug = pgErr
	}

	return New(ctx, http.StatusConflict, DuplicatedRecord, fieldErr)
}

func handleColumnViolation(ctx context.Context, pgErr *pgconn.PgError) *CustomError {
	fieldErr := FieldError{
		Field:   "",
		Input:   "",
		Message: pgErr.Message,
	}

	if config.GetInstance().Api.Environment != string(config.Production) {
		fieldErr.Debug = pgErr
	}

	return New(ctx, http.StatusBadRequest, InvalidValues, fieldErr)
}
