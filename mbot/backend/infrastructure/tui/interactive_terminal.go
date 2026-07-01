package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hpcloud/tail"
)

type newMessage struct {
	text string
}

type Model struct {
	logs     []string
	tail     *tail.Tail
	viewport viewport.Model
}

func NewModel(tail *tail.Tail) *Model {
	return &Model{
		tail: tail,
	}
}

func (m Model) Init() tea.Cmd {
	return waitForLine(m.tail)
}

func waitForLine(logs *tail.Tail) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-logs.Lines
		if !ok || line == nil {
			return nil
		}
		return newMessage{text: line.Text}
	}
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case newMessage:
		rendered := renderMessage(msg.text)
		m.logs = append(m.logs, rendered)
		m.viewport.SetContent(strings.Join(m.logs, "\n"))
		m.viewport.GotoBottom()
		return m, waitForLine(m.tail)

	case tea.WindowSizeMsg:
		m.viewport = viewport.New(msg.Width, msg.Height)
		m.viewport.SetContent(strings.Join(m.logs, "\n"))
	}

	return m, nil
}

func renderMessage(msg string) string {
	parts := strings.SplitN(msg, " | ", 3)
	if len(parts) == 3 {

		var color lipgloss.Color

		switch parts[0] {
		case "Success":
			color = lipgloss.Color("#22C55E")
		case "Error":
			color = lipgloss.Color("#EF4444")
		case "Failure":
			color = lipgloss.Color("#DC2626")
		case "Info":
			color = lipgloss.Color("#3B82F6")
		}

		style := lipgloss.NewStyle().Foreground(color)
		rendered := style.Render(msg)
		return rendered

	}

	return msg
}

func (m Model) View() string {
	return m.viewport.View()
}
