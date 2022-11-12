package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle   = lipgloss.NewStyle().Padding(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"})
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	invalidInputStyle = lipgloss.NewStyle().Background(lipgloss.Color("#FF400050")).Foreground(lipgloss.Color("#FFFDF5"))
)

type item struct {
	title string
}

func (i item) FilterValue() string { return "" }
func (i item) Title() string       { return i.title }

type model struct {
	query          textinput.Model
	querySubmitted bool

	resultLinks  list.Model
	selectedLink string
	linkSelected bool

	resultLimits  list.Model
	selectedLimit string
	limitSelected bool

	quit bool
	err  error
}

type errMsg struct {
	err error
}

func (e errMsg) Error() string {
	return e.err.Error()
}

func InitialModel() model {
	var (
		delegateKeys = newDelegateKeyMap()
	)

	// Setup limit list
	var limits = []list.Item{item{title: "5"}, item{title: "10"}, item{title: "20"}}
	delegate := newItemDelegate(delegateKeys)
	l := list.New(limits, delegate, 0, 0)
	l.Title = "How many results you would like to see?"
	l.SetShowStatusBar(true)
	l.Styles.StatusBar = statusMessageStyle
	l.SetFilteringEnabled(false)
	l.SetShowPagination(false)
	l.Styles.Title = titleStyle
	l.Styles.HelpStyle = helpStyle

	// Setup query textinput
	t := textinput.New()
	t.Placeholder = "C'mon pal ask me anything!"
	t.Validate = func(s string) error {
		if strings.Trim(s, " ") == "" {
			return errors.New("please enter a valid input")
		}
		return nil
	}
	t.Focus()
	t.CharLimit = 250
	t.Width = 20

	return model{
		resultLimits:   l,
		limitSelected:  false,
		query:          t,
		querySubmitted: false,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		key := msg.String()
		if key == "ctrl+c" {
			m.quit = true
			return m, tea.Quit
		}
	}

	if m.err != nil {
		m.query.PlaceholderStyle = invalidInputStyle
		m.query.Placeholder = "Invalid input!!!"
	}

	if !m.querySubmitted {
		return m.updateQuery(msg)
	}

	if !m.limitSelected {
		return m.updateLimits(msg)
	}

	return m, nil
}

func (m model) View() string {
	s := "\nWelcome to Stacked-In CLI!\n\n\n%s"

	if m.quit {
		return quitTextStyle.Render("\nSee you soon pal!\n")
	}

	if !m.querySubmitted {
		return appStyle.Render(fmt.Sprintf(s, m.query.View()))
	}

	if !m.limitSelected {
		return appStyle.Render(fmt.Sprintf(s, m.resultLimits.View()))
	}

	return s
}

func (m model) updateQuery(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			err := m.query.Validate(m.query.Value())
			if err != nil {
				m.err = err
				return m, nil
			}
			m.querySubmitted = true
			m.query.Blur()
			return m, nil
		}

	}

	m.query, cmd = m.query.Update(msg)
	return m, cmd
}

func (m model) updateLimits(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "backspace":
			m.querySubmitted = false
			m.query.Focus()
			return m, nil
		case "enter", " ":
			i, ok := m.resultLimits.SelectedItem().(item)
			if ok {
				m.limitSelected = true
				m.selectedLimit = string(i.title)
			}
			return m, nil
		}

	}

	newListModel, cmd := m.resultLimits.Update(msg)
	m.resultLimits = newListModel

	return m, cmd
}
