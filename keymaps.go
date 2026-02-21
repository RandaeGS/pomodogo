package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Toggle key.Binding
	Exit   key.Binding
}

var DefaultKeyMap = KeyMap{
	Toggle: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Start/Stop timer"),
	),
	Exit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("Ctrl+C", "Exit"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Toggle, k.Exit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Toggle, k.Exit},
	}
}
