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
	// Colors
	primaryColor   = lipgloss.AdaptiveColor{Light: "#4205f7", Dark: "#baf705"}
	secondaryColor = lipgloss.AdaptiveColor{Light: "#2a05f7", Dark: "#d2f705"}
	warningColor   = lipgloss.Color("#f72b05")

	// General
	appStyle      = lipgloss.NewStyle().Padding(1, 2)
	appTitleStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			Background(primaryColor).
			Foreground(lipgloss.AdaptiveColor{Light: "#f5f5f5", Dark: "#151515"}).
			Padding(1).
			Align(lipgloss.Center).
			Bold(true)
	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)

	// Query Input
	queryInputTextStyle       = lipgloss.NewStyle().Foreground(primaryColor).Bold(true)
	queryInputBackgroundStyle = lipgloss.NewStyle().Foreground(secondaryColor)
	invalidInputStyle         = lipgloss.NewStyle().Foreground(warningColor)

	// List
	listTitleStyle = lipgloss.NewStyle().
			Foreground(primaryColor)
	listStatusBarStyle = lipgloss.NewStyle().
				Foreground(secondaryColor)
	listItemStyle = lipgloss.
			NewStyle().
			Foreground(primaryColor).
			Padding(0, 0, 0, 2).
			Border(lipgloss.RoundedBorder(), false, false, false, true).BorderForeground(primaryColor)
	listHelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(1).PaddingBottom(0)
)

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type model struct {
	query          textinput.Model
	querySubmitted bool

	// resultLinks  list.Model
	// linkSelected bool

	resultLimits  list.Model
	limitSelected bool

	quit bool
	err  error
}

// type errMsg struct {
// 	err error
// }

// func (e errMsg) Error() string {
// 	return e.err.Error()
// }

func InitialModel() model {
	var (
		delegateKeys = newDelegateKeyMap()
	)

	// Setup limit list
	var limits = []list.Item{item{title: "5", description: "I'm really close to crying"}, item{title: "10", description: "Why is this not rendering?!?"}, item{title: "20", description: "No idea lol."}}
	d := newItemDelegate(delegateKeys)

	l := list.New(limits, d, 75, 20)
	l.Title = "How many results you would like to see?"
	l.SetShowStatusBar(false)
	l.Styles.StatusBar = listStatusBarStyle
	l.Styles.Title = listTitleStyle
	l.Styles.HelpStyle = listHelpStyle

	// Setup query textinput
	t := textinput.New()
	t.TextStyle = queryInputTextStyle
	t.BackgroundStyle = queryInputBackgroundStyle
	t.Placeholder = "C'mon pal ask me anything!"
	t.Validate = func(s string) error {
		if strings.Trim(s, " ") == "" {
			return errors.New("please enter a valid input")
		}
		return nil
	}
	t.Focus()
	t.CharLimit = 250
	t.Width = 30

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

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// This gets called at the initialization and everytime window gets resized.
		appTitleStyle.Width(msg.Width - (appTitleStyle.GetBorderLeftSize() + appTitleStyle.GetBorderRightSize() + appStyle.GetHorizontalPadding()))
		appStyle.Width(msg.Width)
		appStyle.Height(msg.Height)
		if m.querySubmitted && !m.limitSelected {
			m.resultLimits.SetHeight(msg.Height - appTitleStyle.GetHeight() - 10)
			m.resultLimits.SetWidth(msg.Width)
		}
		return m, nil
	case tea.KeyMsg:
		if m.resultLimits.FilterState() == list.Filtering {
			break
		}

		if msg.Type == tea.KeyCtrlC {
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
	title := appTitleStyle.Render("Stacked-In CLI")

	if m.quit {
		return quitTextStyle.Render("\nSee you soon pal!\n")
	}

	if !m.querySubmitted {
		return appStyle.Render(fmt.Sprintf("%s\n\n%s", title, m.query.View()))
	}

	if !m.limitSelected {
		return appStyle.Render(fmt.Sprintf("%s\n\n%s", title, m.resultLimits.View()))
	}

	if m.limitSelected {
		return appStyle.Render(fmt.Sprintf("%s\n\n%s", title, m.resultLimits.SelectedItem().(item).Title()))
	}

	return title
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
	var cmd tea.Cmd

	// TODO: There is something flaky with the transition between query and result limit

	if m.resultLimits.FilterState() == list.Filtering {
		m.resultLimits, cmd = m.resultLimits.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "backspace":
			// TODO: By default q and esc triggers quit cmd, eliminate this behavior to navigate between views.
			m.querySubmitted = false
			m.query.Focus()
			return m, nil
		case "enter", " ":
			// TODO: call the scrapper and show links, once clicks browser should open OR use viewport and display answer (maybe most liked)
			m.limitSelected = true
			m.resultLimits, cmd = m.resultLimits.Update(msg)
			return m, cmd
		}
	}

	m.resultLimits, cmd = m.resultLimits.Update(msg)
	return m, cmd
}
