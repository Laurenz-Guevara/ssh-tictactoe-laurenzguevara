package game

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	PlayerX = "X"
	PlayerO = "O"
)

type Game struct {
	board         [3][3]string
	currentPlayer string
	winner        string
}

func (g *Game) ResetBoard() {
	g.board = [3][3]string{
		{" ", " ", " "},
		{" ", " ", " "},
		{" ", " ", " "},
	}
	g.winner = ""
}

func (g *Game) CheckWin() string {
	b := g.board

	for i := range 3 {
		if b[i][0] != " " && b[i][0] == b[i][1] && b[i][1] == b[i][2] {
			return b[i][0]
		}
	}

	for i := range 3 {
		if b[0][i] != " " && b[0][i] == b[1][i] && b[1][i] == b[2][i] {
			return b[0][i]
		}
	}

	if b[0][0] != " " && b[0][0] == b[1][1] && b[1][1] == b[2][2] {
		return b[0][0]
	}
	if b[0][2] != " " && b[0][2] == b[1][1] && b[1][1] == b[2][0] {
		return b[0][2]
	}

	return ""
}

func (g *Game) IsBoardFull() bool {
	for _, row := range g.board {
		for _, cell := range row {
			if cell == " " {
				return false
			}
		}
	}
	return true
}

type model struct {
	game      Game
	cursor    [2]int
	gameOver  bool
	askReplay bool
}

func InitialModel() model {
	var g Game
	g.ResetBoard()
	g.currentPlayer = PlayerO

	return model{
		cursor:    [2]int{0, 0},
		game:      g,
		gameOver:  false,
		askReplay: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		x, y := m.cursor[0], m.cursor[1]

		if m.askReplay {
			if key == "y" {
				m.game.ResetBoard()
				m.game.currentPlayer = PlayerO
				m.cursor = [2]int{0, 0}
				m.gameOver = false
				m.askReplay = false
			} else if key == "n" || key == "q" {
				return m, tea.Quit
			}
			return m, nil
		}

		if m.gameOver {
			m.askReplay = true
			return m, nil
		}

		switch key {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if y > 0 {
				y--
			}

		case "down", "j":
			if y < 2 {
				y++
			}

		case "left", "h":
			if x > 0 {
				x--
			}

		case "right", "l":
			if x < 2 {
				x++
			}

		case "enter", " ":
			if m.game.board[y][x] == " " && m.game.winner == "" {
				m.game.board[y][x] = m.game.currentPlayer

				if winner := m.game.CheckWin(); winner != "" {
					m.game.winner = winner
					m.gameOver = true
					m.askReplay = true
				} else if m.game.IsBoardFull() {
					m.game.winner = "Tie"
					m.gameOver = true
					m.askReplay = true
				} else {
					if m.game.currentPlayer == PlayerX {
						m.game.currentPlayer = PlayerO
					} else {
						m.game.currentPlayer = PlayerX
					}
				}
			}
		}

		m.cursor[0] = x
		m.cursor[1] = y
	}

	return m, nil
}

func (m model) View() string {
	s := "Tic Tac Toe\n\n"

	for rowIdx, row := range m.game.board {
		for colIdx, cell := range row {
			val := cell
			if val == " " {
				val = "_"
			}

			if m.cursor[0] == colIdx && m.cursor[1] == rowIdx {
				val = "[" + val + "]"
			} else {
				val = " " + val + " "
			}

			s += val
			if colIdx < len(row)-1 {
				s += "|"
			}
		}
		s += "\n"
	}

	if m.game.winner == "Tie" {
		s += "\nIt's a tie!"
	}
	if m.game.winner != "" && m.game.winner != "Tie" {
		s += fmt.Sprintf("\nPlayer %s wins!", m.game.winner)
	}
	if m.askReplay {
		s += "\n\nPlay again? (y/n): "
	} else {
		s += fmt.Sprintf("\n\nCurrent Player: %s", m.game.currentPlayer)
		s += "\n\nUse arrow keys or hjkl. Press space/enter to place. Press q to quit.\n"
	}

	return s
}

func main() {
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
