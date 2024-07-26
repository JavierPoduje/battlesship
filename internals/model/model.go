package model

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/javierpoduje/battlesshiplib"
)

// tea.WithAltScreen) on a session by session basis.
func TeaHandler(redisConn *battlesshiplib.Redis) func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	return func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
		pty, _, _ := s.Pty()
		renderer := bubbletea.MakeRenderer(s)

		bg := "light"
		if renderer.HasDarkBackground() {
			bg = "dark"
		}

		m := NewModel(pty, bg, redisConn)
		return m, []tea.ProgramOption{tea.WithAltScreen()}
	}
}

type Model struct {
	bg           string
	pty          ssh.Pty
	redisConn    *battlesshiplib.Redis
	clientHeight int
	clientWidth  int
}

func NewModel(pty ssh.Pty, bg string, redisConn *battlesshiplib.Redis) Model {
	return Model{
		bg:           bg,
		pty:          pty,
		redisConn:    redisConn,
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
		case "a":
			a, err := m.redisConn.Get("a")
			if err != nil {
				log.Error("Could not get key", "error", err)
				return m, tea.Quit
			}

			aInt, err := strconv.Atoi(a)
			if err != nil {
				log.Error("Could not convert key to int", "error", err)
				return m, tea.Quit
			}

			aInt++

			err = m.redisConn.Set("a", strconv.Itoa(aInt), 0)
			if err != nil {
				log.Error("Could not set key", "error", err)
				return m, tea.Quit
			}

			return m, nil
		case "b":
			b, err := m.redisConn.Get("b")
			if err != nil {
				log.Error("Could not get key", "error", err)
				return m, tea.Quit
			}

			bInt, err := strconv.Atoi(b)
			if err != nil {
				log.Error("Could not convert key to int", "error", err)
				return m, tea.Quit
			}

			bInt++

			err = m.redisConn.Set("b", strconv.Itoa(bInt), 0)
			if err != nil {
				log.Error("Could not set key", "error", err)
				return m, tea.Quit
			}

			return m, nil
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	a, err := m.redisConn.Get("a")
	if err != nil {
		log.Error("View - couldn't get a", "error", err)
	}
	b, err := m.redisConn.Get("b")
	if err != nil {
		log.Error("View - couldn't get b", "error", err)
	}

	return gloss.Place(
		m.clientWidth, m.clientHeight,
		gloss.Center, gloss.Center,
		gloss.JoinVertical(
			gloss.Center,
			gloss.JoinHorizontal(
				gloss.Center,
				gloss.NewStyle().
					BorderStyle(gloss.RoundedBorder()).
					BorderForeground(gloss.Color("#8BD5CA")).
					Align(gloss.Center, gloss.Center).
					Render("A"),
				a,
				gloss.NewStyle().
					BorderStyle(gloss.RoundedBorder()).
					BorderForeground(gloss.Color("#8BD5CA")).
					Align(gloss.Center, gloss.Center).
					Render("B"),
				b,
			),
		),
	)
}
