package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type HomeModel struct {
	cursor       int
	options      []string
	backupModel  tea.Model
	restoreModel tea.Model
}

func InitialHomeModel(backupModel tea.Model, restoreModel tea.Model) HomeModel {
	return HomeModel{
		cursor:       0,
		options:      []string{"Backup", "Restore"},
		backupModel:  backupModel,
		restoreModel: restoreModel,
	}
}

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.options[m.cursor] == "Backup" {
				return m.backupModel, nil
			} else if m.options[m.cursor] == "Restore" {
				return m.restoreModel, nil
			}
		}
	}

	return m, nil
}

func (m HomeModel) View() string {
	s := "\n\n"

	for i, option := range m.options {
		cursor := "  "
		if i == m.cursor {
			cursor = "(>"
		}
		s += fmt.Sprintf("%s %s\n", cursor, option)
	}

	return s
}
