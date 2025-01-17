package controller

import (
	"encoding/json"
	"fmt"
	//"github.com/oprimogus/cardapiogo/internal/api"
	"github.com/oprimogus/cardapiogo/internal/api/middleware"
	"github.com/oprimogus/cardapiogo/internal/xerrors"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	data := middleware.GetRequestData(r.Context())
	xerror := xerrors.HandleError(r.Context(), err, data.TraceID)
	w.WriteHeader(xerror.Status)
	_ = json.NewEncoder(w).Encode(xerror)
}

func JSONResponse(w http.ResponseWriter, status int, response interface{}) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func IsValidBool(queryParam string) (bool, error) {
	return strconv.ParseBool(queryParam)
}

func isValidInt(queryParam string) (int, error) {
	return strconv.Atoi(queryParam)
}

func GetFileFormData(w http.ResponseWriter, r *http.Request, maxSize int64, key string, types []string) (multipart.File, string, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxSize<<20)

	if err := r.ParseMultipartForm(maxSize << 20); err != nil {
		message := fmt.Sprintf("arquivo excede o tamanho máximo: %d MB", maxSize)
		err = xerrors.BadRequest("", message)
		return nil, "", err
	}

	file, handler, err := r.FormFile(key)
	if err != nil {
		err = xerrors.BadRequest("", "falha ao recuperar o arquivo enviado")
		return nil, "", err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			slog.ErrorContext(r.Context(), "erro ao fechar o arquivo", "error", cerr.Error())
		}
	}()

	ext := filepath.Ext(handler.Filename)

	// Validate MIME type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		err = xerrors.InternalServer("", "falha ao ler o arquivo")
		return nil, "", err
	}

	if _, err := file.Seek(0, 0); err != nil {
		err = xerrors.InternalServer("", "falha ao reposicionar o cursor do arquivo")
		return nil, "", err
	}

	mimeTypeDetected := http.DetectContentType(buffer)
	mimeAllowed := false
	for _, t := range types {
		if t == mimeTypeDetected {
			mimeAllowed = true
			break
		}
	}

	if !mimeAllowed {
		message := fmt.Sprintf("formato de arquivo inválido. O tipo %s não é permitido", mimeTypeDetected)
		err = xerrors.BadRequest("", message)
		return nil, "", err
	}

	return file, ext, nil
}
