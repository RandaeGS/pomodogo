package main

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	dialogBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 3).
			Width(50)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	fieldStyle = lipgloss.NewStyle().
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)

type SaveTimesMsg struct {
	Work      int
	ShortRest int
	LongRest  int
}

type QuitDialogMsg struct{}

type configDialog struct {
	work      textinput.Model
	shortRest textinput.Model
	longRest  textinput.Model
	focused   int
	width     int
	height    int
}

func NewConfigDialog(workSec, shortSec, longSec int, w, h int) configDialog {
	workMin := workSec / 60
	shortMin := shortSec / 60
	longMin := longSec / 60

	ti := func(placeholder string, value int) textinput.Model {
		t := textinput.New()
		t.Placeholder = placeholder
		t.SetValue(fmt.Sprintf("%d", value))
		t.CharLimit = 4
		t.Width = 3
		return t
	}

	d := configDialog{
		work:      ti("Work (minutes)", workMin),
		shortRest: ti("Short rest (minutes)", shortMin),
		longRest:  ti("Long rest (minutes)", longMin),
		focused:   0,
		width:     w,
		height:    h,
	}

	d.work.Focus()
	return d
}

// Init, Update, View remain exactly the same as your last version
// (only the save part is slightly cleaner)

func (d configDialog) Init() tea.Cmd {
	return textinput.Blink
}

func (d configDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
			d.focused = (d.focused + 1) % 3
			d.work.Blur()
			d.shortRest.Blur()
			d.longRest.Blur()
			switch d.focused {
			case 0:
				d.work.Focus()
			case 1:
				d.shortRest.Focus()
			case 2:
				d.longRest.Focus()
			}
			return d, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
			d.focused = (d.focused - 1 + 3) % 3
			d.work.Blur()
			d.shortRest.Blur()
			d.longRest.Blur()
			switch d.focused {
			case 0:
				d.work.Focus()
			case 1:
				d.shortRest.Focus()
			case 2:
				d.longRest.Focus()
			}
			return d, nil

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			w, _ := strconv.Atoi(d.work.Value())
			s, _ := strconv.Atoi(d.shortRest.Value())
			l, _ := strconv.Atoi(d.longRest.Value())

			if w < 1 {
				w = 1
			}
			if s < 1 {
				s = 1
			}
			if l < 1 {
				l = 1
			}

			// Convert minutes back to seconds before saving
			return d, tea.Batch(
				func() tea.Msg {
					return SaveTimesMsg{
						Work:      w * 60,
						ShortRest: s * 60,
						LongRest:  l * 60,
					}
				},
				func() tea.Msg { return QuitDialogMsg{} },
			)

		case key.Matches(msg, key.NewBinding(key.WithKeys("esc", "ctrl+c"))):
			return d, func() tea.Msg { return QuitDialogMsg{} }
		}
	}

	// Forward to focused input
	switch d.focused {
	case 0:
		d.work, cmd = d.work.Update(msg)
	case 1:
		d.shortRest, cmd = d.shortRest.Update(msg)
	case 2:
		d.longRest, cmd = d.longRest.Update(msg)
	}

	return d, cmd
}

func (d configDialog) View() string {
	title := titleStyle.Render("Configure Pomodoro Times")

	fields := []string{
		fieldStyle.Render(fmt.Sprintf("  Work:       %s min", d.work.View())),
		fieldStyle.Render(fmt.Sprintf("  Short rest: %s min", d.shortRest.View())),
		fieldStyle.Render(fmt.Sprintf("  Long rest:  %s min", d.longRest.View())),
	}

	content := lipgloss.JoinVertical(lipgloss.Left, fields...)
	helpText := helpStyle.Render("\nTab / Shift+Tab = move • Enter = save • Esc = cancel")

	box := dialogBorder.Render(lipgloss.JoinVertical(lipgloss.Left, title, content, helpText))

	return lipgloss.Place(
		d.width, d.height,
		lipgloss.Center, lipgloss.Center,
		box,
		lipgloss.WithWhitespaceChars("·"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("235")),
	)
}
