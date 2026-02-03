package main

import (
	"log"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

type model struct {
	isWorking     bool
	shortRestTime int
	longeRestTime int
	workTime      int
	timer         timer.Model
	width         int
	height        int
}

func (m model) Init() tea.Cmd {
	return m.timer.Stop()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		exec.Command("notify-send", "Finished").Run()
		m.isWorking = !m.isWorking
		m.timer = timer.NewWithInterval(time.Duration(time.Second*3), time.Second)
		return m, m.timer.Stop()

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
		case "enter":
			return m, m.timer.Toggle()

		}
	}
	return m, nil
}

func (m model) View() string {
	backgroundColor := lipgloss.Color("#AA2B1D")
	if !m.isWorking {
		backgroundColor = lipgloss.Color("#088395")
	}
	return lipgloss.NewStyle().Background(backgroundColor).Render(
		lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(
				lipgloss.Center,
				Title(),
				m.Timer(),
			),
		),
	)
}

func (m model) Timer() string {
	timer := figure.NewFigure(m.timer.View(), "larry3d", true).String()
	return lipgloss.NewStyle().AlignVertical(lipgloss.Center).Render(timer)
}

func Title() string {
	title := figure.NewFigure("POMODOGO", "larry3d", true).String()
	return lipgloss.NewStyle().AlignVertical(lipgloss.Center).Render(title)
}

func main() {
	m := model{
		isWorking:     true,
		shortRestTime: int(time.Minute * 5),
		longeRestTime: int(time.Minute * 15),
		workTime:      int(time.Minute * 25),
	}
	m.timer = timer.NewWithInterval(time.Duration(time.Second*5), time.Second)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
