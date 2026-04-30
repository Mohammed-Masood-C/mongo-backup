package models

import (
	"fmt"
	coreModels "mongo-backup/internal/core/models"
	"mongo-backup/internal/core/ports"

	tea "github.com/charmbracelet/bubbletea"
)

type BackupModel struct {
	cursor        int
	options       []coreModels.Configuration
	isBackingUp   bool
	backupService ports.BackupService
}

type backupFinishedMsg struct {
	err error
}

func InitialBackupModel(options []coreModels.Configuration, backupService ports.BackupService) BackupModel {
	modifiedOptions := append([]coreModels.Configuration{{DatabaseName: "all"}}, options...)

	return BackupModel{
		cursor:        0,
		options:       modifiedOptions,
		isBackingUp:   false,
		backupService: backupService,
	}
}

func (m BackupModel) Init() tea.Cmd {
	return nil
}

func (m BackupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			m.isBackingUp = true
			return m, m.runBackupCmd()
		}
	case backupFinishedMsg:
		m.isBackingUp = false
		return m, tea.Quit
	}
	return m, nil
}

func (m BackupModel) View() string {
	s := "\n\n"

	if m.isBackingUp {
		s += fmt.Sprintf("Backing up %s in progress, please wait...", m.options[m.cursor].DatabaseName)
		return s
	}

	s += "[Available Databases] :\n\n"
	for i, option := range m.options {
		cursor := "  "
		if i == m.cursor {
			cursor = "(>"
		}
		s += fmt.Sprintf("%s %s\n", cursor, option.DatabaseName)
	}

	return s
}

func (m BackupModel) runBackupCmd() tea.Cmd {
	return func() tea.Msg {
		for _, option := range m.options {
			selectedDatabaseName := m.options[m.cursor].DatabaseName
			if selectedDatabaseName == "all" || selectedDatabaseName == option.DatabaseName {
				if option.DatabaseName == "all" {
					continue
				}
				m.backupService.CreateBackup(option)
			}
		}

		return backupFinishedMsg{err: nil}
	}
}
