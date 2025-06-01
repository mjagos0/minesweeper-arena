package main

import (
	"fmt"
	"math/rand"
)

const (
	MINE = 9
	FLAG = 10
	WIPE = 11 // Fields wiped by pressed mines
)

type Board struct {
	Width    int
	Height   int
	NumMines int
	Fields   []uint8
}

func NewBoard(width, height, numMines int) *Board {
	b := &Board{
		Width:    width,
		Height:   height,
		NumMines: numMines,
		Fields:   make([]uint8, width*height),
	}
	b.DistributeMines(0, 0, b.Width-1, b.Height-1, b.NumMines)
	return b
}

func (b *Board) getIndex(x, y int) int {
	return y*b.Width + x
}

func (b *Board) At(x, y int) uint8 {
	return b.Fields[b.getIndex(x, y)]
}

func (b *Board) Set(x, y int, value uint8) {
	b.Fields[b.getIndex(x, y)] = value
}

func (b *Board) inBounds(x, y int) bool {
	return x >= 0 && x < b.Width && y >= 0 && y < b.Height
}

func (b *Board) DistributeMines(x0, y0, x1, y1, numMines int) {
	if x0 > x1 || y0 > y1 || x0 < 0 || y0 < 0 || x1 >= b.Width || y1 >= b.Height {
		panic("Invalid region bounds")
	}

	regionWidth := x1 - x0 + 1
	regionHeight := y1 - y0 + 1
	totalRegion := regionWidth * regionHeight

	if numMines >= totalRegion {
		panic("Too many mines for the specified region")
	}

	// Generate random positions within the region
	positions := rand.Perm(totalRegion)[:numMines]
	for _, pos := range positions {
		rx := pos % regionWidth
		ry := pos / regionWidth
		b.Set(x0+rx, y0+ry, MINE)
	}

	// Update numbers around mines
	dx := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	dy := []int{-1, -1, -1, 0, 0, 1, 1, 1}

	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if b.At(x, y) != MINE {
				continue
			}
			for d := 0; d < 8; d++ {
				nx := x + dx[d]
				ny := y + dy[d]
				if b.inBounds(nx, ny) && b.At(nx, ny) != MINE {
					b.Set(nx, ny, b.At(nx, ny)+1)
				}
			}
		}
	}
}

func (b *Board) Print() {
	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			v := b.At(x, y)
			if v == MINE {
				fmt.Print("* ")
			} else {
				fmt.Printf("%d ", v)
			}
		}
		fmt.Println()
	}
}
