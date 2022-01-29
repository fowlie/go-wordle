package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	correct = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#538D4E")).PaddingTop(0).PaddingLeft(1).PaddingRight(1)
	present = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#B59F3B")).PaddingTop(0).PaddingLeft(1).PaddingRight(1)
	absent = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#3A3A3C")).PaddingTop(0).PaddingLeft(1).PaddingRight(1)
	idle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).
		PaddingTop(0).PaddingLeft(1).PaddingRight(1)
	gameOver = false
	win      = false
)

type board struct {
	rows       []string
	currentRow int
}

func (b board) AddLetter(l string) {
	b.rows[b.currentRow] += l
}

type model struct {
	b     board
	fasit string
}

func initialModel() model {
	return model{
		b: board{
			rows:       make([]string, 6),
			currentRow: 0,
		},
		fasit: "SKOLE",
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func Update(m model) (string, int) {
	state := "playing"
	row := m.b.currentRow

	if m.b.rows[row] == m.fasit {
		state = "win"
	}

	if row == 5 {
		state = "loose"
	}

	if m.b.currentRow < 5 {
		row++
	}

	return state, row
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			state, row := Update(m)
			switch state {
			case "win":
				win = true
				return m, tea.Quit
			case "loose":
				gameOver = true
				return m, tea.Quit
			case "playing":
				m.b.currentRow = row
			}

		case "backspace":
			m.b.rows[m.b.currentRow] = m.b.rows[m.b.currentRow][:len(m.b.rows[m.b.currentRow])-1]

		default:
			if (win || gameOver) && msg.String() == "q" {
				return m, tea.Quit
			}
			if len(msg.String()) == 1 && len(m.b.rows[m.b.currentRow]) < 5 {
				l := strings.ToUpper(msg.String())
				m.b.AddLetter(l)
			}
		}

	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {

	var s string = "Welcome to wordle game.\nTry to guess the right word!\n\n"

	// Render the rows
	for i := 0; i <= m.b.currentRow; i++ {
		for j := 0; j < len(m.b.rows[i]); j++ {
			l := string(m.b.rows[i][j])
			color := idle
			if i < m.b.currentRow {
				color = m.GetColorStyleAtIndex(l, j)
			}
			s += color.Render(l)
		}
		s += "\n"
	}

	if gameOver {
		s += "Game over!\n"
	} else if win {
		s += "You win!\n"
	}

	// Send the UI for rendering
	return s
}

func (m model) GetColorStyleAtIndex(s string, i int) lipgloss.Style {

	if string(m.fasit[i]) == s { // is the letter the same as in the fasit?
		return correct
	} else if strings.Contains(m.fasit, s) { // ...or do we at least have the letter?
		return present
	}
	return idle // we don't have the letter
}

func main() {
	p := tea.NewProgram(initialModel())

	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
