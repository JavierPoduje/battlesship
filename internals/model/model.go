package model

import (
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
)

// tea.WithAltScreen) on a session by session basis.
func TeaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()
	renderer := bubbletea.MakeRenderer(s)

	bg := "light"
	if renderer.HasDarkBackground() {
		bg = "dark"
	}

	m := NewModel(pty, bg)
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

type Model struct {
	bg           string
	pty          ssh.Pty
	clientHeight int
	clientWidth  int
}

func NewModel(pty ssh.Pty, bg string) Model {
	return Model{
		bg:           bg,
		pty:          pty,
		clientHeight: pty.Window.Height,
		clientWidth:  pty.Window.Width,
	}
}

func (m *Model) HandleWindowResize(msg tea.WindowSizeMsg) {
	m.clientWidth, m.clientHeight = msg.Width, msg.Height
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.HandleWindowResize(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	return gloss.Place(
		m.clientWidth, m.clientHeight,
		gloss.Center, gloss.Center,
		gloss.NewStyle().
			BorderStyle(gloss.RoundedBorder()).
			BorderForeground(gloss.Color("#8BD5CA")).
			Align(gloss.Center, gloss.Center).
			Render("Battleship!"))
}
