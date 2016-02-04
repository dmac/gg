package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/dmac/gg"
	mgl "github.com/go-gl/mathgl/mgl32"
)

const (
	Padding      = 100
	CellSize     = 16
	WidthCells   = 10
	HeightCells  = 14
	WindowWidth  = 2*Padding + CellSize*WidthCells
	WindowHeight = 2*Padding + CellSize*HeightCells
)

type Tetris struct {
	textures map[string]*gg.Texture
	bg       *Sprite
	board    *Board
	score    int

	mu sync.Mutex
	gameOver bool
}

func NewTetris(vertShader, fragShader string) (*Tetris, error) {
	gg.Enable(gg.BLEND)
	gg.BlendFunc(gg.SRC_ALPHA, gg.ONE_MINUS_SRC_ALPHA)

	// Compile the global shader program
	vshader, err := gg.CreateShader([]byte(vertShader), gg.VERTEX_SHADER)
	if err != nil {
		return nil, err
	}
	fshader, err := gg.CreateShader([]byte(fragShader), gg.FRAGMENT_SHADER)
	if err != nil {
		return nil, err
	}
	program := gg.CreateProgram()
	gg.AttachShader(program, vshader)
	gg.AttachShader(program, fshader)
	if err := gg.LinkProgram(program); err != nil {
		return nil, err
	}

	// Set the global projection matrix
	gg.UseProgram(program)
	proj := mgl.Ortho(0, float32(WindowWidth), float32(WindowHeight), 0, 0, 1)
	projUniform, err := gg.GetUniformLocation(program, "proj")
	if err != nil {
		return nil, err
	}
	gg.UniformMatrix4fv(projUniform, proj[:])

	tetris := &Tetris{}

	tetris.textures, err = LoadTextures()
	if err != nil {
		return nil, err
	}

	tetris.bg = NewSprite(WindowWidth, WindowHeight, program, tetris.textures["bg"])
	tetris.board = NewBoard(WidthCells, HeightCells, program, tetris.textures)
	tetris.board.x = Padding
	tetris.board.y = Padding

	rand.Seed(time.Now().Unix())
	tetris.board.current = tetris.board.NewRandomPiece()

	go func() {
		t := time.NewTicker(time.Second)
		for range t.C {
			tetris.mu.Lock()
			if tetris.gameOver {
				return
			}
			tetris.mu.Unlock()
			tetris.HandleInput(inputDown)
		}
	}()

	return tetris, nil
}

func (t *Tetris) Draw() {
	gg.ClearColor(0.5, 0.5, 0.5, 1.0)
	gg.Clear(gg.COLOR_BUFFER_BIT | gg.DEPTH_BUFFER_BIT)
	t.bg.Draw()
	t.board.Draw()
}

type Input byte

const (
	inputUp Input = iota
	inputRight
	inputDown
	inputLeft
	inputSpace
)

func (t *Tetris) HandleInput(input Input) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.gameOver {
		return
	}
	movedPiece := t.board.current.Copy()
	switch input {
	case inputUp:
		movedPiece.orientation = (movedPiece.orientation + 1) % len(orientations[movedPiece.kind])
	case inputDown:
		movedPiece.row += 1
	case inputLeft:
		movedPiece.col -= 1
	case inputRight:
		movedPiece.col += 1
	case inputSpace:
		for movedPiece.Valid() {
			t.board.current = movedPiece.Copy()
			movedPiece.row += 1
		}
	}
	if movedPiece.Valid() {
		t.board.current = movedPiece
		return
	}
	if input != inputDown && input != inputSpace {
		return
	}
	fmt.Printf("%#v\n", t.board.current)
	t.board.AnchorCurrent()
	cleared := t.board.ClearLines()
	if cleared > 0 {
		t.score += cleared
		fmt.Println("Score:", t.score)
	}
	t.board.current = t.board.NewRandomPiece()
	if !t.board.current.Valid() {
		fmt.Println("Game over!")
		t.gameOver = true
	}
}

type Color byte

const (
	NoColor Color = iota
	Red
	Orange
	Yellow
	Green
	Blue
	Cyan
	Purple
)

type Board struct {
	x, y          float32
	width, height int
	bg            *Sprite
	blockSprites  map[Color]*Sprite
	cells         [][]Color
	current       *Piece
}

func NewBoard(width, height int, program *gg.Program, textures map[string]*gg.Texture) *Board {
	const boardOriginX, boardOriginY = CellSize * 4, CellSize * 4
	board := &Board{
		x:      boardOriginX,
		y:      boardOriginY,
		width:  width,
		height: height,
	}
	board.bg = NewSprite(float32(width)*CellSize, float32(height)*CellSize, program, textures["board"])
	board.cells = make([][]Color, board.height)
	for row := 0; row < board.height; row++ {
		board.cells[row] = make([]Color, board.width)
	}

	board.blockSprites = make(map[Color]*Sprite)
	board.blockSprites[Red] = NewSprite(CellSize, CellSize, program, textures["red"])
	board.blockSprites[Orange] = NewSprite(CellSize, CellSize, program, textures["orange"])
	board.blockSprites[Yellow] = NewSprite(CellSize, CellSize, program, textures["yellow"])
	board.blockSprites[Green] = NewSprite(CellSize, CellSize, program, textures["green"])
	board.blockSprites[Blue] = NewSprite(CellSize, CellSize, program, textures["blue"])
	board.blockSprites[Cyan] = NewSprite(CellSize, CellSize, program, textures["cyan"])
	board.blockSprites[Purple] = NewSprite(CellSize, CellSize, program, textures["purple"])

	return board
}

func (b *Board) Draw() {
	b.bg.x = b.x
	b.bg.y = b.y
	b.bg.Draw()
	b.current.Draw()
	for row := 0; row < b.height; row++ {
		for col := 0; col < b.width; col++ {
			spr, ok := b.blockSprites[b.cells[row][col]]
			if !ok {
				continue
			}
			spr.x = b.x + CellSize*float32(col)
			spr.y = b.y + CellSize*float32(row)
			spr.Draw()
		}
	}
}

type PieceKind byte

const (
	I PieceKind = iota
	O
	T
	S
	Z
	L
	J
)

type Piece struct {
	board       *Board
	kind        PieceKind
	row, col    int
	orientation int
}

func (b *Board) NewPiece(kind PieceKind, row, col int) *Piece {
	return &Piece{
		board: b,
		kind:  kind,
		row:   row,
		col:   col,
	}
}

func (b *Board) NewRandomPiece() *Piece {
	kind := PieceKind(rand.Intn(7))
	return b.NewPiece(kind, 0, len(b.cells[0])/2)
}

func (b *Board) ClearLines() (cleared int) {
	for row := len(b.cells) - 1; row >= 0; {
		full := true
		for col := 0; col < len(b.cells[0]); col++ {
			if b.cells[row][col] == NoColor {
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
			for c := 0; c < len(b.cells[0]); c++ {
				b.cells[r][c] = b.cells[r-1][c]
			}
		}
		for c := 0; c < len(b.cells[0]); c++ {
			b.cells[0][c] = NoColor
		}
	}
	return cleared
}

func (b *Board) AnchorCurrent() {
	p := b.current
	cells := orientations[p.kind][p.orientation]
	for row := 0; row < len(cells); row++ {
		for col := 0; col < len(cells[0]); col++ {
			if cells[row][col] {
				b.cells[p.row+row][p.col+col] = p.Color()
			}
		}
	}
}

func (p *Piece) Draw() {
	if p == nil {
		return
	}
	spr := p.board.blockSprites[p.Color()]
	cells := orientations[p.kind][p.orientation]
	for row := 0; row < len(cells); row++ {
		for col := 0; col < len(cells[0]); col++ {
			if !cells[row][col] {
				continue
			}
			spr.x = p.board.x + float32(p.col+col)*CellSize
			spr.y = p.board.y + float32(p.row+row)*CellSize
			spr.Draw()
		}
	}
}

func (p *Piece) Color() Color {
	switch p.kind {
	case I:
		return Cyan
	case O:
		return Yellow
	case T:
		return Purple
	case S:
		return Green
	case Z:
		return Red
	case L:
		return Orange
	case J:
		return Blue
	}
	return NoColor
}

func (p *Piece) Copy() *Piece {
	tmp := *p
	return &tmp
}

func (p *Piece) Valid() bool {
	if p.row < 0 || p.col < 0 {
		return false
	}
	cells := orientations[p.kind][p.orientation]
	if p.row+len(cells) > len(p.board.cells) {
		return false
	}
	if p.col+len(cells[0]) > len(p.board.cells[0]) {
		return false
	}
	for row := 0; row < len(cells); row++ {
		for col := 0; col < len(cells[0]); col++ {
			if cells[row][col] && p.board.cells[p.row+row][p.col+col] != NoColor {
				return false
			}
		}
	}
	return true
}

type Sprite struct {
	x, y     float32
	rotation float32
	scale    float32
	width    float32
	height   float32

	program *gg.Program
	pbuf    *gg.Buffer
	tbuf    *gg.Buffer
	tex     *gg.Texture
}

func NewSprite(width, height float32, program *gg.Program, texture *gg.Texture) *Sprite {
	s := &Sprite{
		scale:   1,
		width:   width,
		height:  height,
		tex:     texture,
		program: program,
	}

	pvertices := []float32{
		0, 0, 0,
		0, s.height, 0,
		s.width, s.height, 0,
		s.width, 0, 0,
	}
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, pvertices); err != nil {
		panic(err)
	}
	s.pbuf = gg.CreateBuffer()
	gg.BindBuffer(gg.ARRAY_BUFFER, s.pbuf)
	gg.BufferData(gg.ARRAY_BUFFER, buf.Bytes(), gg.STATIC_DRAW)

	tvertices := []float32{
		0, 0,
		0, 1,
		1, 1,
		1, 0,
	}
	buf.Reset()
	if err := binary.Write(buf, binary.LittleEndian, tvertices); err != nil {
		panic(err)
	}
	s.tbuf = gg.CreateBuffer()
	gg.BindBuffer(gg.ARRAY_BUFFER, s.tbuf)
	gg.BufferData(gg.ARRAY_BUFFER, buf.Bytes(), gg.STATIC_DRAW)
	return s
}

func (s *Sprite) Draw() error {
	gg.UseProgram(s.program)

	model := s.transform()
	modelUniform, err := gg.GetUniformLocation(s.program, "model")
	if err != nil {
		return err
	}
	gg.UniformMatrix4fv(modelUniform, model[:])

	gg.ActiveTexture(gg.TEXTURE0)
	gg.BindTexture(gg.TEXTURE_2D, s.tex)
	textureUniform, err := gg.GetUniformLocation(s.program, "tex_loc")
	if err != nil {
		return err
	}
	gg.Uniform1i(textureUniform, 0)

	vattrib, err := gg.GetAttribLocation(s.program, "vertex_position")
	if err != nil {
		return err
	}
	gg.EnableVertexAttribArray(vattrib)
	gg.BindBuffer(gg.ARRAY_BUFFER, s.pbuf)
	gg.VertexAttribPointer(vattrib, 3, gg.FLOAT, false, 0, 0)

	tattrib, err := gg.GetAttribLocation(s.program, "vertex_texture")
	if err != nil {
		return err
	}
	gg.EnableVertexAttribArray(tattrib)
	gg.BindBuffer(gg.ARRAY_BUFFER, s.tbuf)
	gg.VertexAttribPointer(tattrib, 2, gg.FLOAT, false, 0, 0)

	gg.DrawArrays(gg.TRIANGLE_FAN, 0, 4)
	return nil
}

func (s *Sprite) transform() mgl.Mat4 {
	S := mgl.Scale2D(s.scale, s.scale).Mat4()
	R := mgl.Rotate2D(mgl.DegToRad(s.rotation)).Mat4()
	T := mgl.Translate3D(s.x, s.y, 0)
	return T.Mul4(R).Mul4(S)
}

//TODO(dmac) Replace hard-coded rotations with real matrix rotation
type Orientation [][]bool

var orientations = map[PieceKind][]Orientation{
	L: {
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
	J: {
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
	O: {
		{
			{true, true},
			{true, true},
		},
	},
	I: {
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
	S: {
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
	Z: {
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
	T: {
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
