package controller

import (
	"encoding/json"
	"fmt"
	"github.com/oprimogus/cardapiogo/internal/application/api/middleware"
	"github.com/oprimogus/cardapiogo/internal/xerrors"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type mimeType string

const (
	JPEG mimeType = "image/jpeg"
	PNG  mimeType = "image/png"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	data := middleware.GetRequestData(r.Context())
	xerror := xerrors.HandleError(err, data.TraceID)
	w.WriteHeader(xerror.Status)
	_ = json.NewEncoder(w).Encode(xerror)
}

func JSONResponse(w http.ResponseWriter, status int, response interface{}) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func GetFileFormData(w http.ResponseWriter, r *http.Request, maxSize int64, key string, types []mimeType) (multipart.File, *multipart.FileHeader, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxSize<<20)

	if err := r.ParseMultipartForm(maxSize << 20); err != nil {
		message := fmt.Sprintf("arquivo excede o tamanho máximo: %d MB", maxSize)
		xerror := xerrors.BadRequest("", message)
		return nil, nil, xerror
	}

	file, handler, err := r.FormFile(key)
	if err != nil {
		xerror := xerrors.BadRequest("", "falha ao recuperar o arquivo enviado")
		return nil, nil, xerror
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			slog.ErrorContext(r.Context(), "erro ao fechar o arquivo", "error", cerr.Error())
		}
	}()

	// Validate file extension
	ext := filepath.Ext(handler.Filename)
	allowedExtensions := make(map[string]bool)
	for _, t := range types {
		allowedExtensions[string(t)] = true
	}

	if !allowedExtensions[ext] {
		message := fmt.Sprintf("tipo de arquivo não suportado. Os tipos permitidos são: %v", types)
		xerror := xerrors.BadRequest("", message)
		return nil, nil, xerror
	}

	// Validate MIME type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		xerror := xerrors.InternalServer("", "falha ao ler o arquivo")
		return nil, nil, xerror
	}

	if _, err := file.Seek(0, 0); err != nil {
		xerror := xerrors.InternalServer("", "falha ao reposicionar o cursor do arquivo")
		return nil, nil, xerror
	}

	mimeTypeDetected := mimeType(http.DetectContentType(buffer))
	mimeAllowed := false
	for _, t := range types {
		if t == mimeTypeDetected {
			mimeAllowed = true
			break
		}
	}

	if !mimeAllowed {
		message := fmt.Sprintf("formato de arquivo inválido. O tipo %s não é permitido", mimeTypeDetected)
		xerror := xerrors.BadRequest("", message)
		return nil, nil, xerror
	}

	return file, handler, nil
}
