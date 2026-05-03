package application

import (
	"fmt"
	"mongo-backup/internal/adapters/configurator"
	"mongo-backup/internal/core/ports"
	"mongo-backup/internal/services"
	"os"
	"path/filepath"
)

type Application struct {
	ApplicationConfigPath string
	BackupConfigPath      string

	Configurator   ports.Configurator
	BackupService  ports.BackupService
	RestoreService ports.RestoreService
}

func NewApplication(applicationConfigPath string) (*Application, error) {
	backupConfigPath := filepath.Join(applicationConfigPath, "backup")
	err := os.MkdirAll(backupConfigPath, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed creating backup folder at path %v: %w", backupConfigPath, err)
	}

	backupLogsPath := filepath.Join(applicationConfigPath, "backup_logs")
	err = os.MkdirAll(backupLogsPath, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed creating backup logs folder at path %v: %w", backupLogsPath, err)
	}

	restoreLogsPath := filepath.Join(applicationConfigPath, "restore_logs")
	err = os.MkdirAll(restoreLogsPath, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed creating restore logs folder at path %v: %w", restoreLogsPath, err)
	}

	fileSystemConfigurator, err := configurator.NewFileSystem(backupConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed creating new FileSystem Configurator: %w", err)
	}

	return &Application{
		ApplicationConfigPath: applicationConfigPath,
		BackupConfigPath:      backupConfigPath,
		Configurator:          fileSystemConfigurator,
		BackupService:         services.NewBackup(backupLogsPath),
		RestoreService:        services.NewRestore(),
	}, nil
}
