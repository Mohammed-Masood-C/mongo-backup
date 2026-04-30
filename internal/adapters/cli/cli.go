package cli

import (
	"mongo-backup/internal/adapters/cli/models"
	"mongo-backup/internal/application"

	tea "github.com/charmbracelet/bubbletea"
)

type Cli struct {
	app       *application.Application
	rootModel tea.Model
}

func NewCli(app *application.Application) *Cli {

	backupModel := models.InitialBackupModel(app.Configurator.GetConfigurations(), app.BackupService)
	restoreModel := models.InitialRestoreModel(app.Configurator.GetConfigurations(), app.BackupService, app.RestoreService)
	homeModel := models.InitialHomeModel(backupModel, restoreModel)

	return &Cli{
		app:       app,
		rootModel: homeModel,
	}
}

func (c Cli) Run() error {
	program := tea.NewProgram(c.rootModel)
	if _, err := program.Run(); err != nil {
		return err
	}
	return nil
}
