package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func SetWorkingDirToProjectRoot() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current directory: %w", err)
	}

	// Caminha para o diretório pai até encontrar o arquivo go.mod (indicando a raiz do projeto)
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			break
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return fmt.Errorf("could not find project root (go.mod not found)")
		}

		currentDir = parentDir
	}

	// Define o diretório de trabalho como a raiz do projeto
	if err := os.Chdir(currentDir); err != nil {
		return fmt.Errorf("could not change to project root: %w", err)
	}

	return nil
}
