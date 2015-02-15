package main

import (
	"fmt"
)

type Board []byte

func NewBlankBoard(size uint8) *Board {
	if size < 3 || size > 5 {
		panic("Unsupported board size: " + fmt.Sprint(size))
	}

	b := make(Board, size*size)

	return &b
}

func NewBoard(encodedBoard string, size uint8) (b *Board, nextTurn uint8) {
	b = NewBlankBoard(size)
	xCount := 0
	oCount := 0
	nextTurn = X

	for i, char := range encodedBoard {
		if char == 'X' {
			(*b)[i] = X
			xCount++
		} else if char == 'O' {
			(*b)[i] = O
			oCount++
		} else {
			(*b)[i] = BLANK
		}
	}

	if oCount < xCount {
		nextTurn = O
	}

	return b, nextTurn
}

func (b *Board) Size() uint8 {
	switch len(*b) {
	case 9:
		return 3
	case 16:
		return 4
	case 25:
		return 5
	}
	panic("Invalid board length")
}

func (b *Board) Length() int {
	return len(*b)
}

func (b *Board) CheckWin() uint8 {
	size := b.Size()

	var xWins bool
	var oWins bool

	type loopParams struct {
		iTo,
		iIncrement,
		jFrom,
		jTo,
		jIncrement uint8
	}

	tictactoeWinPatterns := [4]loopParams{
		loopParams{size * size, size, 0, size, 1},           // rows
		loopParams{size, 1, 0, size * size, size},           // columns
		loopParams{1, 1, 0, size * size, size + 1},          // diagonal top-left to right-bottom
		loopParams{1, 1, size - 1, size*size - 1, size - 1}, // diagonal top-right to left-bottom
	}

	for _, params := range tictactoeWinPatterns {
		for i := uint8(0); i < params.iTo; i += params.iIncrement {
			xWins = true
			oWins = true
			for j := params.jFrom; j < params.jTo; j += params.jIncrement {
				if (*b)[i+j] != X {
					xWins = false
				}
				if (*b)[i+j] != O {
					oWins = false
				}
			}
			if xWins {
				return X
			}
			if oWins {
				return O
			}
		}
	}

	draw := true
	for _, spot := range *b {
		if spot == BLANK {
			draw = false
		}
	}
	if draw {
		return DRAW
	}

	return BLANK
}

func (b *Board) Encode() string {
	// we're only inserting single-byte runes, so we can use a byte array.
	var buf = make([]byte, len(*b))
	for i, value := range *b {
		switch value {
		case X:
			buf[i] = 'X'
		case O:
			buf[i] = 'O'
		case BLANK:
			buf[i] = '_'
		}
	}
	return string(buf)
}

func (b *Board) RenderHTML(output *Output, gameType string, turn uint8, message string) {

	winner := b.CheckWin()
	size := b.Size()
	sTurn := `X`
	if turn == O {
		sTurn = `O`
	}

	output.Add(`<p style="font-weight: bold;">`)
	if winner == X {
		output.Add(`The winner is X.`)
	} else if winner == O {
		output.Add(`The winner is O.`)
	} else if winner == DRAW {
		output.Add(`It's a draw!`)
	} else {
		output.Add(`It is now Player `, sTurn, `'s turn.`)
	}
	output.Add("</p>\n")
	if len(message) > 0 {
		output.Add(`<p>`, message, "</p>\n")
	}

	output.Add(`<div class="board board`, fmt.Sprint(size), `">`)

	for i := 0; i < len(*b); i++ {
		if (*b)[i] == BLANK && winner == BLANK {
			(*b)[i] = turn
			output.Add(`<a rel="nofollow" href="/`, gameType, `/`, fmt.Sprint(size), `/`, b.Encode(), `">`, sTurn, `</a>`)
			(*b)[i] = BLANK
		} else {
			output.Add(`<span>`)
			if (*b)[i] == X {
				output.Add(`X`)
			} else if (*b)[i] == O {
				output.Add(`O`)
			} else {
				output.Add(`&nbsp;`)
			}
			output.Add(`</span>`)
		}
	}

	output.Add("</div>\n")
}
