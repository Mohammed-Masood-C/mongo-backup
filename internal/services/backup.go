package services

import (
	"fmt"
	"mongo-backup/internal/core/models"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

type Backup struct {
}

func NewBackup() *Backup {
	return &Backup{}
}

func (service *Backup) CreateBackup(config models.Configuration) error {
	localTime := time.Now().Local()
	dateStr := localTime.Format("02-01-2006_15-04-05")

	fileName := fmt.Sprintf("%s-%s", config.DatabaseName, dateStr)
	backupPath := filepath.Join(config.BackupFolderPath, fmt.Sprintf("%s.gz", fileName))
	logPath := filepath.Join(config.BackupFolderPath, fmt.Sprintf("logs_%s.txt", fileName))

	if _, err := os.Stat(backupPath); err == nil {
		return fmt.Errorf("backup skipped: file '%s' already exists", fileName)
	}

	logFile, err := os.Create(logPath)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer logFile.Close()

	cmd := exec.Command("mongodump",
		"--uri="+fmt.Sprintf("%s/%s", config.Uri, config.DatabaseName),
		"--archive="+backupPath,
		"--gzip",
	)

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed backing up %v", fileName)
	}

	return nil
}

func (service *Backup) FetchAllBackups(config models.Configuration) ([]string, error) {
	entries, err := os.ReadDir(config.BackupFolderPath)
	if err != nil {
		return nil, fmt.Errorf("failed reading backupFolderPath at path %v: %w",config.BackupFolderPath,err)
	}

	type fileInfo struct {
		name    string
		modTime int64
	}
	var filteredFiles []fileInfo

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".gz" {
			info, err := entry.Info()
			if err != nil {
				continue // Skip files where info cannot be retrieved
			}
			filteredFiles = append(filteredFiles, fileInfo{
				name:    entry.Name(),
				modTime: info.ModTime().Unix(),
			})
		}
	}

	sort.Slice(filteredFiles, func(i, j int) bool {
		return filteredFiles[i].modTime > filteredFiles[j].modTime
	})

	result := make([]string, len(filteredFiles))
	for i, f := range filteredFiles {
		result[i] = f.name
	}

	return result, nil
}
