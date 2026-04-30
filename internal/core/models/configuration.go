package models

type Configuration struct {
	DatabaseName     string `json:"databaseName"`
	Uri              string `json:"uri"`
	BackupFolderPath string `json:"backupFolderPath"`
}
