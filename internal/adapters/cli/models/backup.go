package models

import (
	"fmt"
	coreModels "mongo-backup/internal/core/models"
	"mongo-backup/internal/core/ports"

	tea "github.com/charmbracelet/bubbletea"
)

type BackupModel struct {
	cursor  int
	options []coreModels.Configuration

	isBackingUp bool
	isCompleted bool

	backupLog   string
	backupError error

	backupService ports.BackupService
}

type backupFinishedMsg struct {
	backupLog string
	err       error
}

func InitialBackupModel(options []coreModels.Configuration, backupService ports.BackupService) BackupModel {
	return BackupModel{
		cursor:  0,
		options: options,

		isBackingUp: false,
		isCompleted: false,

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
		m.isCompleted = true

		m.backupLog = msg.backupLog
		m.backupError = msg.err

		return m, tea.Quit
	}
	return m, nil
}

func (m BackupModel) View() string {
	s := "\n\n"

	if m.isCompleted {
		s += fmt.Sprintf("[Backup Output]\n\n")
		if m.backupError != nil {
			s += fmt.Sprintf("Error:\n%v\n\n", m.backupError)
			s += fmt.Sprintf("Mongodump :\n%v\n\n", m.backupLog)
		} else {
			s += fmt.Sprintf("Mongodump :\n%v\n\n", m.backupLog)
		}
	} else if m.isBackingUp {
		s += fmt.Sprintf("Backing up %s in progress, please wait...", m.options[m.cursor].DatabaseName)
		return s
	} else {
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

	return s
}

func (m BackupModel) runBackupCmd() tea.Cmd {
	return func() tea.Msg {
		option := m.options[m.cursor]
		backupLog, err := m.backupService.CreateBackup(option)
		return backupFinishedMsg{backupLog: backupLog, err: err}
	}
}
