# Tic-Tac-Toe web app written in Go
Tic Tac Toe web application with AI algorithm, written in Go, optimized for efficiency.

Uses the Negamax variant of the Minimax algorithm to compute the best possible move. Uses alpha-beta pruning to optimize for epeed without sacrificing accuracy.

Supports 3x3, 4x4, and 5x5 boards. The move tree can be fully searched in 3x3 mode. In 4x4 and 5x5 modes, search is limited to a reasonably computable depth.

Todo list:
* Parallel processing
* Add arbitrary board size (within reasonable limits; 50x50 is probably too big)

Just dreaming todo list:
* GPU processing
* Add other game types like checkers
