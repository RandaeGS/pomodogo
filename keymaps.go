package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Toggle    key.Binding
	Configure key.Binding
	Settings  key.Binding
	Exit      key.Binding
}

var DefaultKeyMap = KeyMap{
	Toggle: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("Enter", "Start/Stop timer"),
	),
	Settings: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Settings"),
	),
	Configure: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "Configure"),
	),
	Exit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("Ctrl+c", "Exit"),
	),
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Toggle, k.Settings, k.Exit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Toggle, k.Settings, k.Exit},
	}
}
