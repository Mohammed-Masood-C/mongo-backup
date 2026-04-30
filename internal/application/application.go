package application

import (
	"fmt"
	"mongo-backup/internal/adapters/configurator"
	"mongo-backup/internal/core/ports"
	"mongo-backup/internal/services"
)

type Application struct {
	Configurator   ports.Configurator
	BackupService  ports.BackupService
	RestoreService ports.RestoreService
}

func NewApplication(configurationPath string) (*Application, error) {
	fileSystemConfigurator, err := configurator.NewFileSystem(configurationPath)

	if err != nil {
		return nil, fmt.Errorf("failed creating new FileSystem Configurator: %w", err)
	}

	return &Application{
		Configurator:   fileSystemConfigurator,
		BackupService:  services.NewBackup(),
		RestoreService: services.NewRestore(),
	}, nil
}
