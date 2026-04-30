package ports

import "mongo-backup/internal/core/models"

type Configurator interface {
	GetConfigurations() []models.Configuration
}
