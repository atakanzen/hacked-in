package cli

import (
	"errors"
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item int

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %d", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	query          textinput.Model
	querySubmitted bool

	resultLinks  list.Model
	selectedLink string

	resultLimits  list.Model
	selectedLimit int
	limitChosen   bool

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
	var limits = []list.Item{item(5), item(10), item(20)}

	l := list.New(limits, itemDelegate{}, 20, 10)
	l.Title = "How many results you would like to see?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	t := textinput.New()
	t.Placeholder = "C'mon pal ask me anything!"
	t.TextStyle.BorderBottom(true)
	t.Validate = func(s string) error {
		if string(s) == "" {
			return errors.New("not valid")
		}
		return nil
	}
	t.Focus()
	t.CharLimit = 250
	t.Width = 20

	return model{
		resultLimits:   l,
		limitChosen:    false,
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
		if key == "q" && m.querySubmitted || key == "ctrl+c" || key == "esc" {
			m.quit = true
			return m, tea.Quit
		}
	}

	if !m.querySubmitted {
		return m.updateQuery(msg)
	}

	if !m.limitChosen {
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
		return fmt.Sprintf(s, m.query.View())
	}

	if !m.limitChosen {
		return fmt.Sprintf(s, m.resultLimits.View())
	}

	return s
}

func (m model) updateQuery(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.query.Err != nil {
				m.querySubmitted = true
				return m, nil
			}
		}

	}

	m.query, cmd = m.query.Update(msg)
	return m, cmd
}

func (m model) updateLimits(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(tea.KeyMsg); msg.String() {
	case "q", "esc":
		m.querySubmitted = false
		return m, nil
	case "enter", " ":
		i, ok := m.resultLimits.SelectedItem().(item)
		if ok {
			m.limitChosen = true
			m.selectedLimit = int(i)
		}
		return m, nil
	}

	return m, nil
}
