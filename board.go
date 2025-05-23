package main

import (
	"fmt"
	"math/rand"
)

const (
	MINE = 9
)

type Board struct {
	width    int
	height   int
	numMines int
	fields   []uint8
}

func NewBoard(width, height, numMines int) *Board {
	b := &Board{
		width:    width,
		height:   height,
		numMines: numMines,
		fields:   make([]uint8, width*height),
	}
	b.DistributeMines()
	return b
}

func (b *Board) At(x, y int) uint8 {
	return b.fields[y*b.width+x]
}

func (b *Board) Set(x, y int, value uint8) {
	b.fields[y*b.width+x] = value
}

func (b *Board) inBounds(x, y int) bool {
	return x >= 0 && x < b.width && y >= 0 && y < b.height
}

func (b *Board) DistributeMines() {
	total := b.width * b.height
	if b.numMines >= total {
		panic("Too many mines")
	}

	positions := rand.Perm(total)[:b.numMines]
	for _, pos := range positions {
		b.fields[pos] = MINE
	}

	dx := []int{-1, 0, 1, -1, 1, -1, 0, 1}
	dy := []int{-1, -1, -1, 0, 0, 1, 1, 1}

	for i := 0; i < total; i++ {
		if b.fields[i] == MINE {
			x := i % b.width
			y := i / b.width
			for d := 0; d < 8; d++ {
				nx, ny := x+dx[d], y+dy[d]
				if b.inBounds(nx, ny) {
					nIdx := ny*b.width + nx
					if b.fields[nIdx] != MINE {
						b.fields[nIdx]++
					}
				}
			}
		}
	}
}

func (b *Board) Print() {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
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
