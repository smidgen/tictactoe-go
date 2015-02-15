package main

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

/*
type Output has methods for adding common HTML features (like header and footer).
It uses a bytes.Buffer for efficient string concatenation.

Output methods are chainable; ie., output.Header().Add("Hello, ").Add("World!").Footer()
*/

type Output struct {
	mainBuffer    *bytes.Buffer
	headBuffer    *bytes.Buffer
	footBuffer    *bytes.Buffer
	includeHeader bool
	includeFooter bool
}

func NewOutput(htmlWrapping bool) *Output {
	return &Output{new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer), htmlWrapping, htmlWrapping}
}

func (out *Output) Add(s ...string) *Output {
	for _, str := range s {
		out.mainBuffer.WriteString(str)
	}
	return out
}

func (out *Output) AddBytes(b ...[]byte) *Output {
	for _, bytes := range b {
		out.mainBuffer.Write(bytes)
	}
	return out
}

func (out *Output) AddToHead(s ...string) *Output {
	for _, str := range s {
		out.headBuffer.WriteString(str)
	}
	return out
}

func (out *Output) AddToFoot(s ...string) *Output {
	for _, str := range s {
		out.footBuffer.WriteString(str)
	}
	return out
}

func (out *Output) AddHeader() *Output {
	out.includeHeader = true
	return out
}
func (out *Output) RemoveHeader() *Output {
	out.includeHeader = false
	return out
}

func (out *Output) AddFooter() *Output {
	out.includeFooter = true
	return out
}
func (out *Output) RemoveFooter() *Output {
	out.includeFooter = false
	return out
}

func (out *Output) WriteTo(w io.Writer) int64 {
	var next string
	var n int64

	processReturns64 := func(x int64, e error) {
		n += x
		if e != nil {
			panic(e)
		}
	}
	processReturns := func(x int, e error) {
		processReturns64(int64(x), e)
	}

	if out.includeHeader {
		next = `<!DOCTYPE html>
<html>
<head>
	<title>Tic Tac Toe</title>
	<meta charset="utf-8" />
	<meta name="author" content="Nolan Ching" />
	<style type="text/css">
		body {
			font-family: arial;
		}
		a, a.visited {
			color: #0000ee;
		}
		.board a {
			display: inline-block;
			color: transparent;
			text-decoration: none;
			line-height: 75px;
			padding-top: 12px;
			padding-bottom: 12px;
			width: 99px;
			border: 1px solid black;
		}
		.board a:hover {
			color: #cccccc;
		}
		.board {
			font-size: 50pt;
			text-align: center;
			border: 1px solid black;
		}
		.board3 {
			width: 303px;
		}
		.board4 {
			width: 404px;
		}
		.board5 {
			width: 505px;
		}
		.board span {
			display: inline-block;
			line-height: 75px;
			padding-top: 12px;
			padding-bottom: 12px;
			width: 99px;
			border: 1px solid black;
		}
	</style>
`
		processReturns(w.Write([]byte(next)))
		processReturns64(out.headBuffer.WriteTo(w))
		next = `</head>

<body>
<h3>Tic Tac Toe</h3>

<p>New 3x3 Game:	<a rel="nofollow" href="/c/3">Player vs. Computer</a>
			&bull;	<a rel="nofollow" href="/c/3/_________">Computer vs. Player</a>
			&bull;	<a rel="nofollow" href="/p/3">Player vs. Player</a></p>

<p>New 4x4 Game:	<a rel="nofollow" href="/c/4">Player vs. Computer</a>
			&bull;	<a rel="nofollow" href="/c/4/________________">Computer vs. Player</a>
			&bull;	<a rel="nofollow" href="/p/4">Player vs. Player</a></p>

<p>New 5x5 Game:	<a rel="nofollow" href="/c/5">Player vs. Computer</a>
			&bull;	<a rel="nofollow" href="/c/5/_________________________">Computer vs. Player</a>
			&bull;	<a rel="nofollow" href="/p/5">Player vs. Player</a></p>

`
		processReturns(w.Write([]byte(next)))
	}

	if out.includeFooter {
		out.mainBuffer.WriteString("\n\n<br /><br />")
		out.mainBuffer.WriteString(fmt.Sprintf("<footer>&#169; %v Nolan Ching</footer>", time.Now().Year()))
		out.footBuffer.WriteString("\n</body>\n</html>")
	}
	processReturns64(out.mainBuffer.WriteTo(w))
	processReturns64(out.footBuffer.WriteTo(w))

	return n
}
