package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dmac/gg"
)

type Tetris struct {
	bg       *gg.Poly
	board    *Board
	gameOver bool
	ticker   *time.Ticker
	score    int
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
		bg:     bg,
		board:  board,
		ticker: time.NewTicker(time.Second),
	}

	return t
}

func (t *Tetris) Update() {
	select {
	case <-t.ticker.C:
		t.HandleInput(inputDown)
	default:
	}
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
		movedPiece.orientation = (movedPiece.orientation + 1) % len(orientations[movedPiece.typ])
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
	}
	// Anchor piece, clear lines, new piece
	if input == inputDown {
		t.board.AnchorCurrent()
		cleared := t.board.ClearLines()
		if cleared > 0 {
			t.score += cleared
			fmt.Println("Score:", t.score)
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
	grid    [][]string
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

	grid := make([][]string, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]string, width)
	}

	board := &Board{
		bg:   bg,
		grid: grid,
		sprs: make(map[string]*gg.Sprite),
	}
	board.current = NewPiece(board)
	board.sprs["I"] = gg.NewSpriteFromTexture(textures["cyan"])
	board.sprs["J"] = gg.NewSpriteFromTexture(textures["blue"])
	board.sprs["L"] = gg.NewSpriteFromTexture(textures["orange"])
	board.sprs["S"] = gg.NewSpriteFromTexture(textures["green"])
	board.sprs["Z"] = gg.NewSpriteFromTexture(textures["red"])
	board.sprs["O"] = gg.NewSpriteFromTexture(textures["yellow"])
	board.sprs["T"] = gg.NewSpriteFromTexture(textures["purple"])

	return board
}

func (b *Board) AnchorCurrent() {
	cells := orientations[b.current.typ][b.current.orientation]
	for row, cols := range cells {
		for col, cell := range cols {
			if cell {
				b.current.board.grid[b.current.row+row][b.current.col+col] = b.current.typ
			}
		}
	}
}

func (b *Board) ClearLines() (cleared int) {
	for row := len(b.grid) - 1; row >= 0; {
		full := true
		for _, cell := range b.grid[row] {
			if cell == "" {
				full = false
				break
			}
		}
		if !full {
			row--
			continue
		}
		cleared++
		for r := row; r > 0; r-- {
			for c := 0; c < len(b.grid[r]); c++ {
				b.grid[r][c] = b.grid[r-1][c]
			}
		}
		for c := 0; c < len(b.grid[0]); c++ {
			b.grid[0][c] = ""
		}
	}
	return
}

func (b *Board) Draw() {
	b.bg.Draw()
	b.current.Draw()

	for row, cols := range b.grid {
		for col, cell := range cols {
			if cell != "" {
				spr := b.sprs[b.grid[row][col]]
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
	spr         *gg.Sprite
	board       *Board
	row, col    int
	typ         string
	orientation int
}

func NewPiece(board *Board) *Piece {
	var typ, color string
	switch rand.Intn(len(orientations)) {
	case 0:
		typ = "L"
		color = "orange"
	case 1:
		typ = "J"
		color = "blue"
	case 2:
		typ = "O"
		color = "yellow"
	case 3:
		typ = "I"
		color = "cyan"
	case 4:
		typ = "S"
		color = "green"
	case 5:
		typ = "Z"
		color = "red"
	case 6:
		typ = "T"
		color = "purple"
	}

	return &Piece{
		spr:         gg.NewSpriteFromTexture(textures[color]),
		row:         1,
		col:         2,
		board:       board,
		typ:         typ,
		orientation: 0,
	}
}

func (p *Piece) Draw() {
	boardPosition := p.board.bg.Position
	cells := orientations[p.typ][p.orientation]
	for row, cols := range cells {
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
	cells := orientations[p.typ][p.orientation]
	if p.col+len(cells[0]) > len(p.board.grid[0]) {
		return false
	}
	if p.row+len(cells) > len(p.board.grid) {
		return false
	}
	for row, cols := range cells {
		for col, cell := range cols {
			if cell && p.board.grid[p.row+row][p.col+col] != "" {
				return false
			}
		}
	}
	return true
}

func (p *Piece) Copy() *Piece {
	return &Piece{
		spr:         p.spr,
		board:       p.board,
		row:         p.row,
		col:         p.col,
		typ:         p.typ,
		orientation: p.orientation,
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

type Orientation [][]bool

var orientations = map[string][]Orientation{
	"L": {
		{
			{true, false, false},
			{true, true, true},
		},
		{
			{true, true},
			{true, false},
			{true, false},
		},
		{
			{true, true, true},
			{false, false, true},
		},
		{
			{false, true},
			{false, true},
			{true, true},
		},
	},
	"J": {
		{
			{false, false, true},
			{true, true, true},
		},
		{
			{true, false},
			{true, false},
			{true, true},
		},
		{
			{true, true, true},
			{true, false, false},
		},
		{
			{true, true},
			{false, true},
			{false, true},
		},
	},
	"O": {
		{
			{true, true},
			{true, true},
		},
	},
	"I": {
		{
			{true, true, true, true},
		},
		{
			{true},
			{true},
			{true},
			{true},
		},
	},
	"S": {
		{
			{false, true, true},
			{true, true, false},
		},
		{
			{true, false},
			{true, true},
			{false, true},
		},
	},
	"Z": {
		{
			{true, true, false},
			{false, true, true},
		},
		{
			{false, true},
			{true, true},
			{true, false},
		},
	},
	"T": {
		{
			{false, true, false},
			{true, true, true},
		},
		{
			{true, false},
			{true, true},
			{true, false},
		},
		{
			{true, true, true},
			{false, true, false},
		},
		{
			{false, true},
			{true, true},
			{false, true},
		},
	},
}
