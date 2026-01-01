package main

import (
	"math"
)

type Bot struct {
	Difficulty string // "easy", "medium", "hard"
}

func NewBot(difficulty string) *Bot {
	if difficulty == "" {
		difficulty = "medium"
	}
	return &Bot{Difficulty: difficulty}
}

// GetBotMove returns the column where the bot wants to drop a disc
func (bot *Bot) GetBotMove(board *Board, botPlayer int, opponentPlayer int) int {
	validMoves := board.GetValidMoves()
	
	if len(validMoves) == 0 {
		return -1
	}

	bestScore := math.MinInt32
	bestMove := validMoves[0]

	for _, col := range validMoves {
		score := bot.evaluateMove(board, col, botPlayer, opponentPlayer)
		if score > bestScore {
			bestScore = score
			bestMove = col
		}
	}

	return bestMove
}

func (bot *Bot) evaluateMove(board *Board, col int, botPlayer int, opponentPlayer int) int {
	testBoard := board.Copy()
	
	row, err := testBoard.DropDisc(col, botPlayer)
	if err != nil {
		return -1000
	}

	// Winning move (highest priority)
	if testBoard.CheckWin(row, col, botPlayer) {
		return 10000
	}

	// Block opponent's winning move
	testBoard2 := board.Copy()
	if row2, err := testBoard2.DropDisc(col, opponentPlayer); err == nil {
		if testBoard2.CheckWin(row2, col, opponentPlayer) {
			return 9000
		}
	}

	// Score based on position and threats
	score := bot.scorePosition(testBoard, col, botPlayer, opponentPlayer)
	
	return score
}

func (bot *Bot) scorePosition(board *Board, col int, botPlayer int, opponentPlayer int) int {
	score := 0
	row := -1
	
	// Find the row where disc would land
	for r := ROWS - 1; r >= 0; r-- {
		if board.Grid[r][col] == botPlayer {
			row = r
			break
		}
	}

	if row == -1 {
		return 0
	}

	// Center columns are more valuable
	centerDistance := abs(col - COLS/2)
	score += (COLS - centerDistance) * 10

	// Count adjacent discs
	score += bot.countAdjacentDiscs(board, row, col, botPlayer) * 50
	
	// Count opponent discs to block
	score += bot.countAdjacentDiscs(board, row, col, opponentPlayer) * 40

	return score
}

func (bot *Bot) countAdjacentDiscs(board *Board, row, col, player int) int {
	count := 0
	
	// Directions: horizontal, vertical, diagonal1, diagonal2
	directions := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	
	for _, dir := range directions {
		// Count in positive direction
		r, c := row+dir[0], col+dir[1]
		for r >= 0 && r < ROWS && c >= 0 && c < COLS && board.Grid[r][c] == player {
			count++
			r += dir[0]
			c += dir[1]
		}
		
		// Count in negative direction
		r, c = row-dir[0], col-dir[1]
		for r >= 0 && r < ROWS && c >= 0 && c < COLS && board.Grid[r][c] == player {
			count++
			r -= dir[0]
			c -= dir[1]
		}
	}
	
	return count
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
