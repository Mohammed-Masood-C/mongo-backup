package configurator

import (
	"encoding/json"
	"fmt"
	"mongo-backup/internal/core/models"
	"os"
	"path/filepath"
)

type FileSystem struct {
	Path           string
	Configurations []models.Configuration
}

func NewFileSystem(path string) (*FileSystem, error) {

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed creating backup folder at path %v: %w", path, err)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed reading backup folder at path %v: %w", path, err)
	}

	var allConfigs []models.Configuration
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			jsonFilePath := filepath.Join(path, file.Name())
			data, err := os.ReadFile(jsonFilePath)
			if err != nil {
				return nil, fmt.Errorf("failed reading json from file %v: %w", jsonFilePath, err)
			}

			var configData models.Configuration
			if err := json.Unmarshal(data, &configData); err != nil {
				return nil, fmt.Errorf("failed parsing config from file %v: %w", jsonFilePath, err)
			}
			if configData.DatabaseName == "" {
				return nil, fmt.Errorf("'databaseName' cannot be empty or missing. found in file %v", jsonFilePath)
			}
			if configData.DatabaseName == "all" {
				return nil, fmt.Errorf("databaseName 'all' is not allowed. found in file %v", jsonFilePath)
			}
			if configData.Uri == "" {
				return nil, fmt.Errorf("'uri' cannot be empty or missing. found in file %v", jsonFilePath)
			}
			if len(configData.Uri) > 0 && configData.Uri[len(configData.Uri)-1] == '/' {
				configData.Uri = configData.Uri[:len(configData.Uri)-1]
			}
			if configData.BackupFolderPath == "" {
				return nil, fmt.Errorf("'backupFolderPath' cannot be empty or missing. found in file %v", jsonFilePath)
			}

			if _, err := os.Stat(configData.BackupFolderPath); err != nil {
				return nil, fmt.Errorf("failed finding backup folder at path %v: %w", configData.BackupFolderPath, err)
			}

			allConfigs = append(allConfigs, configData)
		}
	}

	return &FileSystem{
		Path:           path,
		Configurations: allConfigs,
	}, nil
}

func (c *FileSystem) GetConfigurations() []models.Configuration {
	return c.Configurations
}
