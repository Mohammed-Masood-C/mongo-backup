package services

import (
	"bytes"
	"fmt"
	"mongo-backup/internal/core/models"
	"os"
	"os/exec"
	"path/filepath"
)

type Restore struct {
}

func NewRestore() *Restore {
	return &Restore{}
}

func (service *Restore) RestoreBackup(config models.Configuration, fileName string, mongoConnectionUri string) (string, error) {

	filePath := filepath.Join(config.BackupFolderPath, fileName)
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("restore failed for file at path %v: %w", fileName, err)
	}

	cmd := exec.Command("mongorestore",
		"--uri="+mongoConnectionUri,
		"--archive="+filePath,
		"--gzip",
	)

	var outBuffer bytes.Buffer
	var errBuffer bytes.Buffer

	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	if err := cmd.Run(); err != nil {
		return outBuffer.String() + errBuffer.String(), fmt.Errorf("restore failed to execute for file at path %v: %w", fileName, err)
	}

	return outBuffer.String() + errBuffer.String(), nil
}
