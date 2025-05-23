package main

import (
	"fmt"
)

type Game struct {
	board    *Board
	revealed []bool
	lost     bool
	won      bool
}

type RevealResult struct {
	Positions []RevealedCell `json:"positions"`
	HitMine   bool           `json:"hitMine"`
	Won       bool           `json:"won"`
}

type RevealedCell struct {
	X     int   `json:"x"`
	Y     int   `json:"y"`
	Value uint8 `json:"value"`
}

func NewGame(board *Board) *Game {
	return &Game{
		board:    board,
		revealed: make([]bool, board.width*board.height),
		lost:     false,
		won:      false,
	}
}

func (g *Game) Reveal(x, y int) RevealResult {
	if g.lost || g.won || !g.board.inBounds(x, y) {
		return RevealResult{}
	}

	index := y*g.board.width + x
	if g.revealed[index] {
		return RevealResult{}
	}

	positions := make([]RevealedCell, 0)
	queue := [][2]int{{x, y}}
	visited := make(map[[2]int]bool)

	for len(queue) > 0 {
		nx, ny := queue[0][0], queue[0][1]
		queue = queue[1:]
		nIdx := ny*g.board.width + nx

		if !g.board.inBounds(nx, ny) || g.revealed[nIdx] {
			continue
		}
		g.revealed[nIdx] = true
		val := g.board.At(nx, ny)
		positions = append(positions, RevealedCell{X: nx, Y: ny, Value: val})

		if val == MINE {
			g.lost = true
			return RevealResult{Positions: positions, HitMine: true, Won: false}
		}

		if val == 0 {
			dx := []int{-1, 0, 1, -1, 1, -1, 0, 1}
			dy := []int{-1, -1, -1, 0, 0, 1, 1, 1}
			for d := 0; d < 8; d++ {
				nxx, nyy := nx+dx[d], ny+dy[d]
				if g.board.inBounds(nxx, nyy) {
					nIndex := nyy*g.board.width + nxx
					if !g.revealed[nIndex] && !visited[[2]int{nxx, nyy}] {
						queue = append(queue, [2]int{nxx, nyy})
						visited[[2]int{nxx, nyy}] = true
					}
				}
			}
		}
	}

	total := g.board.width * g.board.height
	revealedCount := 0
	for _, r := range g.revealed {
		if r {
			revealedCount++
		}
	}
	if revealedCount == total-g.board.numMines {
		g.won = true
	}

	return RevealResult{Positions: positions, HitMine: false, Won: g.won}
}

func (g *Game) PrintVisible() {
	for y := 0; y < g.board.height; y++ {
		for x := 0; x < g.board.width; x++ {
			idx := y*g.board.width + x
			if g.revealed[idx] {
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
