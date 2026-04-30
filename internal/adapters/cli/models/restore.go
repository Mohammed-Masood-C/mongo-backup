package models

import (
	"fmt"
	"log"
	coreModels "mongo-backup/internal/core/models"
	"mongo-backup/internal/core/ports"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type RestoreModel struct {
	inConfigScreen bool
	cursorConfig   int
	configOptions  []coreModels.Configuration

	inFileSelectScreen bool
	fileSelectCursor   int
	fileOptions        []string

	inUriScreen bool
	uriInput    textinput.Model

	isRestoring bool

	backupService  ports.BackupService
	restoreService ports.RestoreService
}

type restoreFinishedMsg struct {
	err error
}

func InitialRestoreModel(configOptions []coreModels.Configuration, backupService ports.BackupService, restoreService ports.RestoreService) RestoreModel {
	uriInput := textinput.New()
	uriInput.Placeholder = "mongodb://localhost:27017"
	uriInput.Width = 100

	return RestoreModel{
		inConfigScreen: true,
		cursorConfig:   0,
		configOptions:  configOptions,

		inFileSelectScreen: false,
		fileSelectCursor:   0,
		fileOptions:        []string{},

		inUriScreen: false,
		uriInput:    uriInput,

		isRestoring: false,

		backupService:  backupService,
		restoreService: restoreService,
	}
}

func (m RestoreModel) Init() tea.Cmd {
	return nil
}

func (m RestoreModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	if m.inConfigScreen {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up":
				if m.cursorConfig > 0 {
					m.cursorConfig--
				}
				return m, nil
			case "down":
				if m.cursorConfig < len(m.configOptions)-1 {
					m.cursorConfig++
				}
				return m, nil
			case "enter":
				m.inConfigScreen = false
				m.inFileSelectScreen = true
				fileNames, err := m.backupService.FetchAllBackups(m.configOptions[m.cursorConfig])
				if err != nil {
					log.Fatal(err)
				}
				m.fileOptions = fileNames
				return m, nil
			}
		}
	}

	if m.inFileSelectScreen {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up":
				if m.fileSelectCursor > 0 {
					m.fileSelectCursor--
				}
				return m, nil
			case "down":
				if m.fileSelectCursor < len(m.fileOptions)-1 {
					m.fileSelectCursor++
				}
				return m, nil
			case "enter":
				m.inFileSelectScreen = false
				m.inUriScreen = true

				m.uriInput.Focus()
				return m, textinput.Blink
			}
		}
	}

	if m.inUriScreen {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.inUriScreen = false
				m.isRestoring = true
				return m, m.runRestoreCmd()
			}
		}

		var cmd tea.Cmd
		m.uriInput, cmd = m.uriInput.Update(msg)
		return m, cmd
	}

	if m.isRestoring {
		switch msg.(type) {
		case restoreFinishedMsg:
			m.isRestoring = false
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m RestoreModel) View() string {
	s := "\n\n"

	if m.inConfigScreen {
		s += "[Available Databases] :\n\n"
		for i, option := range m.configOptions {
			cursor := "  "
			if i == m.cursorConfig {
				cursor = "(>"
			}
			s += fmt.Sprintf("%s %s\n", cursor, option.DatabaseName)
		}
	}

	if m.inFileSelectScreen {
		s += "[Available Backups] :\n\n"
		for i, fileName := range m.fileOptions {
			cursor := "  "
			if i == m.fileSelectCursor {
				cursor = "(>"
			}
			s += fmt.Sprintf("%s %s\n", cursor, fileName)
		}
	}

	if m.inUriScreen {
		s += "Enter MongoDB URI to restore to:\n\n"
		s += m.uriInput.View()
		s += "\n\n(press Enter to confirm)"
	}

	if m.isRestoring {
		s += fmt.Sprintf("[Restore in progress]\n\n")
		s += fmt.Sprintf("%v -> %v\n\n", m.fileOptions[m.fileSelectCursor], m.uriInput.Value())
		s += "please wait..."
	}

	return s
}

func (m RestoreModel) runRestoreCmd() tea.Cmd {
	return func() tea.Msg {
		m.restoreService.RestoreBackup(m.configOptions[m.cursorConfig], m.fileOptions[m.fileSelectCursor], m.uriInput.Value())
		return restoreFinishedMsg{err: nil}
	}
}
