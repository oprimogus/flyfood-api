package converters

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
)

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
