package main

import (
	"log"
	"mongo-backup/internal/adapters/cli"
	"mongo-backup/internal/application"
	"os"
	"path/filepath"
)

func main() {
	applicationConfigPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("failed retrieving UserConfigDir: %v", err)
		return
	}
	applicationConfigPath = filepath.Join(applicationConfigPath, "mongo-backup")

	app, err := application.NewApplication(applicationConfigPath)
	if err != nil {
		log.Fatalf("failed creating new application: %v", err)
	}

	cliApp := cli.NewCli(app)

	if err := cliApp.Run(); err != nil {
		log.Fatalf("failed running cli application: %v", err)
	}
}
