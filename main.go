package main

import (
	"log"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

type model struct {
	isWorking      bool
	shortRestTime  int
	shortRestCount int
	longRestTime   int
	longRestCount  int
	workTime       int
	workCount      int
	timer          timer.Model
	help           help.Model
	progress       progress.Model
	activeDialog   tea.Model
	width          int
	height         int
}

func (m model) Init() tea.Cmd {
	return m.timer.Stop()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SaveTimesMsg:
		m.workTime = msg.Work
		m.shortRestTime = msg.ShortRest
		m.longRestTime = msg.LongRest
		m.isWorking = true

		m.timer = timer.NewWithInterval(time.Duration(m.workTime)*time.Second, time.Second)

		m.activeDialog = nil
		return m, m.timer.Init()

	case QuitDialogMsg:
		m.activeDialog = nil
		return m, nil
	}

	if m.activeDialog != nil {
		var cmd tea.Cmd
		m.activeDialog, cmd = m.activeDialog.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = m.width / 2

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m.UpdateCounts()
		message := ""
		if m.isWorking {
			message = "Time to work!"
			m.timer = timer.NewWithInterval(time.Duration(time.Second*time.Duration(m.workTime)), time.Second)
		} else if !m.isWorking && m.workCount%4 != 0 || m.workCount == 0 {
			message = "Time to take a short rest!"
			m.timer = timer.NewWithInterval(time.Duration(time.Second*time.Duration(m.shortRestTime)), time.Second)
		} else if !m.isWorking && m.workCount%4 == 0 && m.workCount != 0 {
			message = "Time to take a long rest!"
			m.timer = timer.NewWithInterval(time.Duration(time.Second*time.Duration(m.longRestTime)), time.Second)
		}

		exec.Command("notify-send", message).Run()
		return m, m.timer.Stop()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Exit):
			return m, tea.Quit
		case key.Matches(msg, DefaultKeyMap.Toggle):
			return m, m.timer.Toggle()

		case key.Matches(msg, DefaultKeyMap.Settings):
			if m.activeDialog == nil {
				m.activeDialog = NewConfigDialog(
					m.workTime,
					m.shortRestTime,
					m.longRestTime,
					m.width,
					m.height,
				)
				return m, m.activeDialog.Init()
			}
		}
	}
	return m, nil
}

func (m *model) UpdateCounts() {
	if m.isWorking {
		m.workCount++
	} else if !m.isWorking && m.workCount%4 != 0 || m.workCount == 0 {
		m.shortRestCount++
	} else if !m.isWorking && m.workCount%4 == 0 && m.workCount != 0 {
		m.longRestCount++
	}
	m.isWorking = !m.isWorking
}

func (m model) View() string {
	if m.activeDialog != nil {
		return m.activeDialog.View()
	}

	bg := lipgloss.Color("#AA2B1D")
	if !m.isWorking {
		bg = lipgloss.Color("#088395")
	}
	if m.timer.Running() {
		bg = lipgloss.Color("#060c0d")
	}

	// Main content: title + timer, centered horizontally and vertically
	main := lipgloss.JoinVertical(
		lipgloss.Center,
		Title(),
		m.Timer(),
	)

	mainCentered := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height / 2).
		Background(bg).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Bottom).
		Render(main)

	// Compose: main content + spacer + help at bottom
	return lipgloss.JoinVertical(
		lipgloss.Left,
		mainCentered,
		m.ProgressBar(bg),
		m.helpView(),
	)
}

func (m model) Timer() string {
	timer := figure.NewFigure(m.timer.View(), "colossal", true).String()
	return lipgloss.NewStyle().AlignVertical(lipgloss.Center).Render(timer)
}

func Title() string {
	title := figure.NewFigure("POMODOGO", "colossal", true).String()
	return lipgloss.NewStyle().AlignVertical(lipgloss.Center).Padding(3).Render(title)
}

func (m model) ProgressBar(bg lipgloss.Color) string {
	percentage := 0.00
	if m.isWorking {
		percentage = m.timer.Timeout.Seconds() / float64(m.workTime)
	} else if !m.isWorking && m.workCount%4 != 0 || m.workCount == 0 {
		percentage = m.timer.Timeout.Seconds() / float64(m.shortRestTime)
	} else if !m.isWorking && m.workCount%4 == 0 && m.workCount != 0 {
		percentage = m.timer.Timeout.Seconds() / float64(m.longRestTime)
	}

	m.progress.Width = m.width / 2
	progress := m.progress.ViewAs(percentage)
	return lipgloss.NewStyle().
		Width(m.width).
		Height((m.height - 1) / 2).
		PaddingTop(3).
		Background(bg).
		AlignHorizontal(lipgloss.Center).
		Render(progress)
}

func (m model) helpView() string {
	base := lipgloss.NewStyle()
	m.help.Styles = help.Styles{
		Ellipsis:       base.Foreground(lipgloss.Color("240")),
		ShortKey:       base.Foreground(lipgloss.Color("228")),
		ShortDesc:      base.Foreground(lipgloss.Color("252")),
		ShortSeparator: base.Foreground(lipgloss.Color("240")),
		FullKey:        base.Foreground(lipgloss.Color("228")),
		FullDesc:       base.Foreground(lipgloss.Color("252")),
		FullSeparator:  base.Foreground(lipgloss.Color("240")),
	}

	bgStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#AA2B1D")).
		Height(1).
		AlignHorizontal(lipgloss.Center)

	helpContent := m.help.View(DefaultKeyMap)

	return bgStyle.Render("") + lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Center,
		helpContent,
	)
}

func main() {
	m := model{
		isWorking:      true,
		shortRestTime:  5 * 60,
		shortRestCount: 0,
		longRestTime:   15 * 60,
		longRestCount:  0,
		workTime:       25 * 60,
		workCount:      0,
		help:           help.New(),
		timer:          timer.NewWithInterval(time.Duration(time.Minute*25), time.Second),
		progress:       progress.New(progress.WithDefaultGradient()),
	}
	m.progress.ShowPercentage = false

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
