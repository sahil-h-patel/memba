package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	message string
	count   int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) View() tea.View {
	v := tea.NewView("Hello, full screen!")
	v.AltScreen = true
	return v
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func RunPretty() {
	p := tea.NewProgram(model{message: "This is the internal state: ", count: 0})
	if _, teaErr := p.Run(); teaErr != nil {
		log.Panic(teaErr)
	}
}
