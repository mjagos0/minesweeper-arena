package main

import (
	"fmt"
)

const (
	UNREVEALED = 0
	REVEALED   = 1
	FLAGGED    = 2
)

type Game struct {
	board    *Board
	state    []int
	revealed int
	finished bool
	won      bool
}

type CellUpdate struct {
	X     int   `json:"x"`
	Y     int   `json:"y"`
	Value uint8 `json:"value"`
}

func (g *Game) IsRevealed(x, y int) bool {
	return g.state[g.board.getIndex(x, y)] == REVEALED
}

func (g *Game) IsFlagged(x, y int) bool {
	return g.state[g.board.getIndex(x, y)] == FLAGGED
}

func (g *Game) UpdateState(x, y, state int) {
	g.state[g.board.getIndex(x, y)] = state
}

func NewGame(board *Board) *Game {
	return &Game{
		board:    board,
		state:    make([]int, board.Width*board.Height),
		revealed: 0,
		finished: false,
		won:      false,
	}
}

func (g *Game) processFlag(x, y int, cellUpdates *[]CellUpdate) {
	if g.IsRevealed(x, y) {
		return
	}

	if g.IsFlagged(x, y) {
		g.UpdateState(x, y, UNREVEALED)
		*cellUpdates = append(*cellUpdates, CellUpdate{X: x, Y: y, Value: WIPE})
	} else {
		g.UpdateState(x, y, FLAGGED)
		*cellUpdates = append(*cellUpdates, CellUpdate{X: x, Y: y, Value: FLAG})
	}
}

func (g *Game) processReveal(x, y int, cellUpdates *[]CellUpdate) {
	queue := [][2]int{{x, y}}
	visited := make(map[[2]int]bool)
	visited[[2]int{x, y}] = true

	for len(queue) > 0 {
		nx, ny := queue[0][0], queue[0][1]
		queue = queue[1:]

		if !g.board.inBounds(nx, ny) || g.IsRevealed(nx, ny) {
			continue
		}
		g.UpdateState(nx, ny, REVEALED)
		val := g.board.At(nx, ny)
		*cellUpdates = append(*cellUpdates, CellUpdate{X: nx, Y: ny, Value: val})

		if val == MINE {
			g.finished = true
			g.won = false
			return
		} else {
			g.revealed++
		}

		if val == 0 {
			dx := []int{-1, 0, 1, -1, 1, -1, 0, 1}
			dy := []int{-1, -1, -1, 0, 0, 1, 1, 1}
			for d := 0; d < 8; d++ {
				nxx, nyy := nx+dx[d], ny+dy[d]
				if g.board.inBounds(nxx, nyy) {
					if !g.IsRevealed(nxx, nyy) && !visited[[2]int{nxx, nyy}] {
						queue = append(queue, [2]int{nxx, nyy})
						visited[[2]int{nxx, nyy}] = true
					}
				}
			}
		}
	}

	if g.revealed == len(g.board.Fields)-g.board.NumMines {
		g.finished = true
		g.won = true
	}
}

func (g *Game) processMove(x, y int, flag bool) []CellUpdate {
	cellUpdates := make([]CellUpdate, 0)
	if g.finished || !g.board.inBounds(x, y) || g.IsRevealed(x, y) {
		return cellUpdates
	} else if flag {
		g.processFlag(x, y, &cellUpdates)
	} else {
		g.processReveal(x, y, &cellUpdates)
	}

	return cellUpdates
}

func (g *Game) Print() {
	for y := 0; y < g.board.Height; y++ {
		for x := 0; x < g.board.Width; x++ {
			if g.IsRevealed(x, y) {
				val := g.board.At(x, y)
				if val == MINE {
					fmt.Print("* ")
				} else {
					fmt.Printf("%d ", val)
				}
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
}
