package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func newItemDelegate(keys *delegateKeyMap) list.DefaultDelegate {
	d := list.NewDefaultDelegate()

	d.ShowDescription = false

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keys.choose):
				return m.NewStatusMessage(fmt.Sprintf("Let's go with %s!", title))
			case key.Matches(msg, keys.cancel):
				return m.NewStatusMessage("I need to go back to input")
			}
		}

		return nil
	}

	help := []key.Binding{keys.choose, keys.cancel}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	return d
}

type delegateKeyMap struct {
	choose key.Binding
	cancel key.Binding
}

func (d delegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.choose,
		d.cancel,
	}
}

func newDelegateKeyMap() *delegateKeyMap {
	return &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose")),
		cancel: key.NewBinding(
			key.WithKeys("esc", "q", "backspace"),
			key.WithHelp("esc/q/backspace", "go back"),
		),
	}
}
