package services

import (
	"bytes"
	"fmt"
	"io"
	"mongo-backup/internal/core/models"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Restore struct {
	restoreLogFolderPath string
}

func NewRestore(restoreLogPath string) *Restore {
	return &Restore{
		restoreLogFolderPath: restoreLogPath,
	}
}

func (service *Restore) RestoreBackup(config models.Configuration, fileName string, mongoConnectionUri string) (string, error) {

	localTime := time.Now().Local()
	dateStr := localTime.Format("02-01-2006_15-04-05")
	restoreLogFileName := fmt.Sprintf("logs_%s-%s.txt", config.DatabaseName, dateStr)

	restoreLogFolderPath := filepath.Join(service.restoreLogFolderPath, config.DatabaseName)
	logPath := filepath.Join(restoreLogFolderPath, restoreLogFileName)

	err := os.MkdirAll(restoreLogFolderPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed creating restore log folder at path %v: %w", restoreLogFolderPath, err)
	}

	filePath := filepath.Join(config.BackupFolderPath, fileName)
	if _, err := os.Stat(filePath); err != nil {
		return "", fmt.Errorf("restore failed for file at path %v: %w", fileName, err)
	}

	cmd := exec.Command("mongorestore",
		"--uri="+mongoConnectionUri,
		"--archive="+filePath,
		"--gzip",
	)

	logFile, err := os.Create(logPath)
	if err != nil {
		return "", fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFile.Close()

	var logBuffer bytes.Buffer

	w := io.MultiWriter(logFile, &logBuffer)

	header := fmt.Sprintf("Restore started at: %s\nTarget Backup: %s\n---\n",
		localTime,
		fileName,
	)

	if _, err := fmt.Fprint(w, header); err != nil {
		return "", fmt.Errorf("failed to write log header: %w", err)
	}

	cmd.Stdout = w
	cmd.Stderr = w

	if err := cmd.Run(); err != nil {
		return logBuffer.String(), fmt.Errorf("restore failed to execute for file at path %v: %w", fileName, err)
	}

	return logBuffer.String(), nil
}
