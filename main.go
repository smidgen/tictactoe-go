/*
	Tic-Tac-Toe web app
	By Nolan Ching <nolan@nolanching.com>

	This program allows the user to play against another player or against
	the computer. If the computer is a player, it uses the Minimax algorithm
	to beat its opponent every time. It allows the user to pick whether the
	computer goes first or not.

	The program is optimized for speed, using alpha-beta pruning to prevent
	the algorithm from uselessly examining too many future moves.
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

const (
	DEBUG = true
	BLANK = 0
	X     = 1
	O     = 2
	DRAW  = 3
)

type GameState struct {
	GameType       string
	Board          *Board
	NextTurn       uint8
	IsComputerTurn bool
}

type TicTacToe struct{}

func (self *TicTacToe) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	output := NewOutput(true)
	defer output.WriteTo(w)

	defer func() {
		if r := recover(); r != nil {
			output.Add(fmt.Sprint("Error: ", r))
			if DEBUG {
				output.Add(`<pre>`)
				output.AddBytes(debug.Stack())
				output.Add(`</pre>`)
			}
		}
	}()

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

	state.Board.RenderHTML(output, state.GameType, state.NextTurn, roboMessage)
}

func (self *TicTacToe) ParseParams(urlpath string) *GameState {

	var params []string = strings.SplitN(urlpath, "/", 4)

	state := new(GameState)
	state.GameType = "p"
	var size uint8 = 3

	if len(params) >= 3 {
		if params[1] != "p" {
			state.GameType = "c"
			state.IsComputerTurn = true
		}
		switch params[2] {
		case `3`:
			size = 3
		case `4`:
			size = 4
		case `5`:
			size = 5
		}
		if len(params) == 4 {
			state.Board, state.NextTurn = NewBoard(params[3], size)
		} else {
			// If no board in URL, human wants to go first.
			state.IsComputerTurn = false
		}
	}

	if state.Board == nil {
		state.Board = NewBlankBoard(size)
		state.NextTurn = X
	}

	return state
}

func (self *TicTacToe) ComputerMove(state *GameState) int8 {

	bestMove, bestMoveValue := NegaMax(state.NextTurn, state.Board, -2, -2)
	if bestMove != -1 {
		(*state.Board)[bestMove] = state.NextTurn
	}

	state.NextTurn = (state.NextTurn % 2) + 1

	return bestMoveValue
}

func NegaMax(player uint8, board *Board, alpha int8, beta int8) (int, int8) {
	result := board.CheckWin()

	var enemy uint8 = (player % 2) + 1

	if result == player { // win
		return -1, 1
	} else if result == enemy { // lose
		return -1, -1
	} else if result == DRAW {
		return -1, 0
	}

	var bestMove int = -1

	// bestMoveValue represents the best possible choice score we can make during this round.
	// We start with something lower than -1, so even if we
	// end up losing, we'll still end up making a move.
	var bestMoveValue int8 = -2

	for i := 0; i < len(*board); i++ {
		if (*board)[i] == BLANK { // for each possible move
			(*board)[i] = player // make the move

			// what is the worst possible move the enemy can do to me?
			// (also happens to be the best move for that enemy)
			_, moveValue := NegaMax(enemy, board, beta, alpha)

			(*board)[i] = BLANK // undo the move

			// This is a zero-sum game: a win for me is a lose for you, and vice versa.
			moveValue = -moveValue

			// if the worst the enemy can do to me is better than the the worst he can do
			// to me if I pick another square, then this square is better than that one.
			if moveValue > bestMoveValue {
				bestMove = i
				bestMoveValue = moveValue
				alpha = moveValue
			}

			// If this branch is better for me than a branch the enemy has already examined,
			// he's not going to pick this branch anyway, so we might as well not look at
			// the rest of our options.
			if -alpha < beta {
				break
			}
		}
	}

	return bestMove, bestMoveValue
}

func main() {

	http.Handle("/", &TicTacToe{})
	log.Fatal(http.ListenAndServe(":4000", nil))
}
