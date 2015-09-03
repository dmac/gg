package main

import (
	"fmt"
	"time"

	"github.com/dmac/gg"
)

type Tetris struct {
	bg       *gg.Poly
	board    *Board
	gameOver bool
}

func NewTetris(windowWidth, windowHeight int) *Tetris {
	bg := gg.NewPoly([][2]float32{
		{0, 0},
		{0, float32(windowHeight)},
		{float32(windowWidth), float32(windowHeight)},
		{float32(windowWidth), 0},
	})
	bg.SetColor(0.3, 0.2, 0.3, 1)
	board := NewBoard(10, 24)
	board.bg.Position = [2]float32{100, 50}

	t := &Tetris{
		bg:    bg,
		board: board,
	}

	go func() {
		for range time.NewTicker(1 * time.Second).C {
			t.HandleInput(inputDown)
			if t.gameOver {
				return
			}
		}
	}()

	return t
}

func (t *Tetris) Draw() {
	t.bg.Draw()
	t.board.Draw()
}

func (t *Tetris) HandleInput(input Input) {
	if t.gameOver {
		return
	}
	current := t.board.current
	movedPiece := current.Copy()
	switch input {
	case inputUp:
		movedPiece.row -= 1
	case inputDown:
		movedPiece.row += 1
	case inputLeft:
		movedPiece.col -= 1
	case inputRight:
		movedPiece.col += 1
	}

	if movedPiece.Valid() {
		t.board.current = movedPiece
		return
	} else if input == inputDown {
		for row, cols := range current.cells {
			for col, cell := range cols {
				if cell {
					current.board.grid[current.row+row][current.col+col] = true
				}
			}
		}
		t.board.current = NewPiece(t.board)
		if !t.board.current.Valid() {
			fmt.Println("Game over!")
			t.gameOver = true
		}
	}
}

type Board struct {
	bg      *gg.Poly
	sprs    map[string]*gg.Sprite
	grid    [][]bool
	current *Piece
}

const cellsize = 16

func NewBoard(width, height int) *Board {
	bg := gg.NewPoly([][2]float32{
		{0, 0},
		{0, float32(height * cellsize)},
		{float32(width * cellsize), float32(height * cellsize)},
		{float32(width * cellsize), 0},
	})
	bg.SetColor(0, 0, 0, 1)

	grid := make([][]bool, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]bool, width)
	}

	board := &Board{
		bg:   bg,
		grid: grid,
		sprs: make(map[string]*gg.Sprite),
	}
	board.current = NewPiece(board)
	board.sprs["orange"] = gg.NewSpriteFromTexture(textures["orange"])

	return board
}

func (b *Board) Draw() {
	b.bg.Draw()
	b.current.Draw()

	for row, cols := range b.grid {
		for col, cell := range cols {
			if cell {
				spr := b.sprs["orange"]
				spr.SetPosition(
					b.bg.Position[0]+float32(cellsize*col),
					b.bg.Position[1]+float32(cellsize*row),
				)
				spr.Draw()
			}
		}
	}
}

type Piece struct {
	spr      *gg.Sprite
	board    *Board
	row, col int
	cells    [][]bool
}

func NewPiece(board *Board) *Piece {
	// #
	// ###
	cells := make([][]bool, 2)
	for i := 0; i < len(cells); i++ {
		cells[i] = make([]bool, 3)
	}
	cells[0][0] = true
	cells[1][0] = true
	cells[1][1] = true
	cells[1][2] = true

	spr := gg.NewSpriteFromTexture(textures["orange"])

	return &Piece{
		spr:   spr,
		row:   1,
		col:   2,
		cells: cells,
		board: board,
	}
}

func (p *Piece) Draw() {
	boardPosition := p.board.bg.Position
	for row, cols := range p.cells {
		for col, cell := range cols {
			if cell {
				p.spr.SetPosition(
					boardPosition[0]+float32(cellsize*(p.col+col)),
					boardPosition[1]+float32(cellsize*(p.row+row)),
				)
				p.spr.Draw()
			}
		}
	}
}

func (p *Piece) Valid() bool {
	if p.row < 0 {
		return false
	}
	if p.col < 0 {
		return false
	}
	if p.col+len(p.cells[0]) > len(p.board.grid[0]) {
		return false
	}
	if p.row+len(p.cells) > len(p.board.grid) {
		return false
	}
	for row, cols := range p.cells {
		for col, cell := range cols {
			if cell && p.board.grid[p.row+row][p.col+col] {
				return false
			}
		}
	}
	return true
}

func (p *Piece) Copy() *Piece {
	return &Piece{
		spr:   p.spr,
		board: p.board,
		row:   p.row,
		col:   p.col,
		cells: p.cells,
	}
}

type Input byte

const (
	inputUp Input = iota
	inputRight
	inputDown
	inputLeft
	inputSpace
)

func (i Input) String() string {
	switch i {
	case inputUp:
		return "up"
	case inputRight:
		return "right"
	case inputDown:
		return "down"
	case inputLeft:
		return "left"
	case inputSpace:
		return "space"
	default:
		return "unknown"
	}
}
