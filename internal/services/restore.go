package services

import (
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

func (service *Restore) RestoreBackup(config models.Configuration, fileName string, mongoConnectionUri string) error {

	filePath := filepath.Join(config.BackupFolderPath, fileName)
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("restore failed for file at path %v: %w", fileName, err)
	}

	cmd := exec.Command("mongorestore",
		"--uri="+mongoConnectionUri,
		"--archive="+filePath,
		"--gzip",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("restore failed to execute for file at path %v: %w", fileName, err)
	}

	return nil
}
