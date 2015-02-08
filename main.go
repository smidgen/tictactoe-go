package main

import (
	"log"
	"net/http"
	"strings"
)

type GameState struct {
	GameType       string
	Board          string
	NextTurn       rune
	IsComputerTurn bool
	Winner         rune
}

type TicTacToe struct{}

func (self *TicTacToe) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	output := NewOutput(true)
	defer output.WriteTo(w)

	state := self.ParseParams(r.URL.Path)

	var roboMessage string

	if state.IsComputerTurn {
		result := self.ComputerMove(state)
		switch result {
		case 1:
			roboMessage = "Nice try, human, but I am going to win."
		case 0:
			roboMessage = "If you're smart, you can still make it a draw, but you can't win."
		case -1:
			roboMessage = "No fair, you cheated!"
		}
	}
	state.Winner = self.CheckWin(state.Board)

	self.RenderBoard(state, roboMessage, output)
}

func (self *TicTacToe) ParseParams(urlpath string) *GameState {

	var params []string = strings.SplitN(urlpath, "/", 3)

	state := &GameState{
		"c",
		"_________",
		'X',
		true,
		'_',
	}

	xCount := 0
	oCount := 0
	if len(params) >= 2 {
		if params[1] != "c" {
			state.GameType = "p"
			state.IsComputerTurn = false
		}
		if len(params) == 3 {
			if len(params[2]) == 9 {
				state.Board = strings.ToUpper(params[2])
				for i := 0; i < 9; i++ { // Expect all characters to be 8-bit ASCII
					if state.Board[i] == 'X' {
						xCount++
					} else if state.Board[i] == 'O' {
						oCount++
					}
				}
			}
		} else {
			// If no board in URL, human wants to go first.
			state.IsComputerTurn = false
		}
	}

	if oCount < xCount {
		state.NextTurn = 'O'
	}

	return state
}

func (self *TicTacToe) ComputerMove(state *GameState) int8 {
	var miniMax func(player rune, board string, alpha int8, beta int8) (string, int8)
	miniMax = func(player rune, board string, alpha int8, beta int8) (string, int8) {
		result := self.CheckWin(board)
		if result != '_' {
			if result == player {
				return board, 1 // win
			} else {
				return board, -1 // lose
			}
		} else if !strings.ContainsRune(board, '_') {
			return board, 0 // draw
		}

		enemy := 'X'
		if player == 'X' {
			enemy = 'O'
		}

		var bestMove string

		// $bestMoveValue represents the best possible choice score we can make during this round.
		// We start with something lower than -1, so even if we
		// end up losing, we'll still end up making a move.
		var bestMoveValue int8 = -2;

		for i := 0; i < 9; i++ {
			if board[i] == '_' { // for each possible move
				// figure out what the resulting board looks like
				move := []byte(board)
				move[i] = byte(player)

				// what is the worst possible move the enemy can do to me?
				// (also happens to be the best move for that enemy)
				_, moveValue := miniMax(enemy, string(move), beta, alpha)

				// This is a zero-sum game: a win for me is a lose for you, and vice versa.
				moveValue = -moveValue

				// if the worst the enemy can do to me is better than the the worst he can do
				// to me if I pick another square, then this square is better than that one.
				if moveValue > bestMoveValue {
					bestMove = string(move)
					bestMoveValue = moveValue
					alpha = moveValue
				}

				// If this path is better for me than a path the enemy has already examined,
				// he's not going to pick this path anyway, so we might as well not look at
				// the rest of our options.
				if -alpha < beta {
					break
				}
			}
		}

		return bestMove, bestMoveValue
	}

	newBoard, bestMoveValue := miniMax(state.NextTurn, state.Board, -2, -2)
	state.Board = newBoard

	if state.NextTurn == 'X' {
		state.NextTurn = 'O'
	} else {
		state.NextTurn = 'X'
	}
	return bestMoveValue
}

func (self *TicTacToe) CheckWin(b string) rune {

	if 'X' == b[0] && 'X' == b[1] && 'X' == b[2] ||
		'X' == b[3] && 'X' == b[4] && 'X' == b[5] ||
		'X' == b[6] && 'X' == b[7] && 'X' == b[8] ||
		'X' == b[0] && 'X' == b[3] && 'X' == b[6] ||
		'X' == b[1] && 'X' == b[4] && 'X' == b[7] ||
		'X' == b[2] && 'X' == b[5] && 'X' == b[8] ||
		'X' == b[0] && 'X' == b[4] && 'X' == b[8] ||
		'X' == b[2] && 'X' == b[4] && 'X' == b[6] {
		return 'X'
	} else if 'O' == b[0] && 'O' == b[1] && 'O' == b[2] ||
		'O' == b[3] && 'O' == b[4] && 'O' == b[5] ||
		'O' == b[6] && 'O' == b[7] && 'O' == b[8] ||
		'O' == b[0] && 'O' == b[3] && 'O' == b[6] ||
		'O' == b[1] && 'O' == b[4] && 'O' == b[7] ||
		'O' == b[2] && 'O' == b[5] && 'O' == b[8] ||
		'O' == b[0] && 'O' == b[4] && 'O' == b[8] ||
		'O' == b[2] && 'O' == b[4] && 'O' == b[6] {
		return 'O'
	} else {
		return '_'
	}
}

func (self *TicTacToe) RenderBoard(state *GameState, roboMessage string, output *Output) {

	output.Add(`<p style="font-weight: bold;">`)
	if state.Winner != '_' {
		output.Add(`The winner is `, string(state.Winner))
	} else if !strings.ContainsRune(state.Board, '_') {
		output.Add(`It's a draw!`)
	} else {
		output.Add(`It is now Player `, string(state.NextTurn), `'s turn.`)
	}
	output.Add("</p>\n")
	if len(roboMessage) > 0 {
		output.Add(`<p>`, roboMessage, "</p>\n")
	}

	output.Add(`<div class="board">`)

	for i := 0; i < 9; i++ {
		if state.Board[i] == '_' && state.Winner == '_' {
			newBoard := []byte(state.Board)
			newBoard[i] = byte(state.NextTurn)
			output.Add(`<a rel="nofollow" href="/`, state.GameType, `/`, string(newBoard), `">`, string(state.NextTurn), `</a>`)
		} else {
			output.Add(`<span>`)
			if state.Board[i] != '_' {
				output.Add(string(state.Board[i]))
			} else {
				output.Add(`&nbsp;`)
			}
			output.Add(`</span>`)
		}
	}

	output.Add("</div>\n")
}

func main() {

	http.Handle("/", &TicTacToe{})
	log.Fatal(http.ListenAndServe("localhost:4000", nil))
}
