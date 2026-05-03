package ports

import "mongo-backup/internal/core/models"

type BackupService interface {
	CreateBackup(config models.Configuration) (string, error)
	FetchAllBackups(config models.Configuration) ([]string, error)
}

type RestoreService interface {
	RestoreBackup(config models.Configuration, fileName string, mongoConnectionUri string) (string, error)
}
