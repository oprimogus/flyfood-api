package converters

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
)

func FileHeaderToBytes(file *multipart.FileHeader) ([]byte, error) {
	fileConverted, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("fail on converted fileHeader: %w", err)
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			slog.Error("fail on close file", "error", err)
		}
	}(fileConverted)

	buff := bytes.NewBuffer(nil)
	if _, err := io.Copy(buff, fileConverted); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func FileToBytes(file multipart.File) ([]byte, error) {
	buff := bytes.NewBuffer(nil)
	if _, err := io.Copy(buff, file); err != nil {
		return nil, fmt.Errorf("fail on copy file: %w", err)
	}
	err := file.Close()
	if err != nil {
		return nil, fmt.Errorf("fail on close file: %w", err)
	}

	return buff.Bytes(), nil

}
