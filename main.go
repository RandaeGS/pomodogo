package main

import (
	"log"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	isRunning     bool
	shortRestTime int
	longeRestTime int
	workTime      int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case " ":
			err := exec.Command("notify-send", "Hello").Run()
			if err != nil {
				log.Fatal(err)
			}
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() string {
	return "Hello World!"
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := tea.NewProgram(model{}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
