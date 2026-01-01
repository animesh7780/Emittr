package main

import (
	"fmt"
)

const (
	ROWS    = 6
	COLS    = 7
	PLAYER1 = 1
	PLAYER2 = 2
	BOT     = 2
	EMPTY   = 0
)

type Board struct {
	Grid [ROWS][COLS]int
}

type GameState struct {
	ID        string
	Board     *Board
	Player1   string
	Player2   string
	CurrentPlayer int
	Status    string // "waiting", "active", "finished"
	Winner    string // "", "player1", "player2", "draw"
	CreatedAt string
	UpdatedAt string
	IsBot     bool
}

func NewBoard() *Board {
	return &Board{
		Grid: [ROWS][COLS]int{},
	}
}

func (b *Board) CanDropDisc(col int) bool {
	if col < 0 || col >= COLS {
		return false
	}
	return b.Grid[0][col] == EMPTY
}

func (b *Board) DropDisc(col int, player int) (int, error) {
	if !b.CanDropDisc(col) {
		return -1, fmt.Errorf("column %d is full", col)
	}

	row := -1
	for r := ROWS - 1; r >= 0; r-- {
		if b.Grid[r][col] == EMPTY {
			row = r
			break
		}
	}

	if row == -1 {
		return -1, fmt.Errorf("column %d is full", col)
	}

	b.Grid[row][col] = player
	return row, nil
}

func (b *Board) CheckWin(row int, col int, player int) bool {
	if row < 0 || row >= ROWS || col < 0 || col >= COLS {
		return false
	}

	// Check horizontal
	if b.checkDirection(row, col, player, 0, 1) {
		return true
	}

	// Check vertical
	if b.checkDirection(row, col, player, 1, 0) {
		return true
	}

	// Check diagonal (top-left to bottom-right)
	if b.checkDirection(row, col, player, 1, 1) {
		return true
	}

	// Check diagonal (top-right to bottom-left)
	if b.checkDirection(row, col, player, 1, -1) {
		return true
	}

	return false
}

func (b *Board) checkDirection(row, col, player, dRow, dCol int) bool {
	count := 1

	// Check positive direction
	r, c := row+dRow, col+dCol
	for r >= 0 && r < ROWS && c >= 0 && c < COLS && b.Grid[r][c] == player {
		count++
		r += dRow
		c += dCol
	}

	// Check negative direction
	r, c = row-dRow, col-dCol
	for r >= 0 && r < ROWS && c >= 0 && c < COLS && b.Grid[r][c] == player {
		count++
		r -= dRow
		c -= dCol
	}

	return count >= 4
}

func (b *Board) IsBoardFull() bool {
	for col := 0; col < COLS; col++ {
		if b.Grid[0][col] == EMPTY {
			return false
		}
	}
	return true
}

func (b *Board) Copy() *Board {
	newBoard := NewBoard()
	for r := 0; r < ROWS; r++ {
		for c := 0; c < COLS; c++ {
			newBoard.Grid[r][c] = b.Grid[r][c]
		}
	}
	return newBoard
}

func (b *Board) GetValidMoves() []int {
	var moves []int
	for col := 0; col < COLS; col++ {
		if b.CanDropDisc(col) {
			moves = append(moves, col)
		}
	}
	return moves
}
