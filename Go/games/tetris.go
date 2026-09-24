package games

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type TetrisGame struct {
	width      int32
	height     int32
	offsetX    int32
	offsetY    int32
	cols       int32
	rows       int32
	cellSize   int32
	board      [][]string
	score      int32
	gameOver   bool
	fallTimer  float64
	fallDelay  float64
	current    string
	pieceShape [][]int
	pieceX     int32
	pieceY     int32
	shapes     map[string][][]int
	colors     map[string]color.RGBA
	tileCache  *TileCache
}

func NewTetrisGame() *TetrisGame {
	g := &TetrisGame{
		width:     300,
		height:    600,
		cols:      10,
		rows:      20,
		cellSize:  30,
		fallDelay: 0.5,
	}
	g.offsetX = (ScreenWidth - g.width - 140) / 2
	g.offsetY = (ScreenHeight - g.height) / 2
	g.tileCache = NewTileCache(float32(g.cellSize-2), float32(g.cellSize-2), 4)

	g.shapes = map[string][][]int{
		"I": {{0, 0, 0, 0}, {1, 1, 1, 1}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		"O": {{1, 1}, {1, 1}},
		"T": {{0, 1, 0}, {1, 1, 1}, {0, 0, 0}},
		"S": {{0, 1, 1}, {1, 1, 0}, {0, 0, 0}},
		"Z": {{1, 1, 0}, {0, 1, 1}, {0, 0, 0}},
		"L": {{1, 0, 0}, {1, 1, 1}, {0, 0, 0}},
		"J": {{0, 0, 1}, {1, 1, 1}, {0, 0, 0}},
	}
	g.colors = map[string]color.RGBA{
		"I": {0, 240, 240, 255},
		"O": {240, 240, 0, 255},
		"T": {180, 0, 240, 255},
		"S": {0, 240, 0, 255},
		"Z": {240, 0, 0, 255},
		"L": {240, 160, 0, 255},
		"J": {0, 0, 240, 255},
	}
	g.Reset()
	return g
}

func (g *TetrisGame) Reset() {
	g.board = make([][]string, g.rows)
	for i := range g.board {
		g.board[i] = make([]string, g.cols)
	}
	g.score = 0
	g.gameOver = false
	g.fallTimer = 0
	g.spawnPiece()
}

func (g *TetrisGame) spawnPiece() {
	keys := make([]string, 0, len(g.shapes))
	for k := range g.shapes {
		keys = append(keys, k)
	}
	g.current = keys[rand.Intn(len(keys))]
	shape := g.shapes[g.current]
	g.pieceShape = make([][]int, len(shape))
	for i, row := range shape {
		g.pieceShape[i] = make([]int, len(row))
		copy(g.pieceShape[i], row)
	}
	g.pieceX = g.cols/2 - int32(len(g.pieceShape[0]))/2
	g.pieceY = 0
	if g.checkCollision(g.pieceShape, g.pieceX, g.pieceY) {
		g.gameOver = true
	}
}

func (g *TetrisGame) checkCollision(shape [][]int, x, y int32) bool {
	for rowIdx, row := range shape {
		for colIdx, cell := range row {
			if cell == 1 {
				bx := x + int32(colIdx)
				by := y + int32(rowIdx)
				if bx < 0 || bx >= g.cols || by >= g.rows {
					return true
				}
				if by >= 0 && g.board[by][bx] != "" {
					return true
				}
			}
		}
	}
	return false
}

func (g *TetrisGame) rotatePiece() {
	shape := g.pieceShape
	rotated := make([][]int, len(shape[0]))
	for i := range rotated {
		rotated[i] = make([]int, len(shape))
	}
	for y := 0; y < len(shape); y++ {
		for x := 0; x < len(shape[y]); x++ {
			rotated[x][len(shape)-1-y] = shape[y][x]
		}
	}
	if !g.checkCollision(rotated, g.pieceX, g.pieceY) {
		g.pieceShape = rotated
	}
}

func (g *TetrisGame) lockPiece() {
	for rowIdx, row := range g.pieceShape {
		for colIdx, cell := range row {
			if cell == 1 {
				by := g.pieceY + int32(rowIdx)
				bx := g.pieceX + int32(colIdx)
				if by >= 0 {
					g.board[by][bx] = g.current
				}
			}
		}
	}
	g.clearLines()
	g.spawnPiece()
}

func (g *TetrisGame) clearLines() {
	lines := 0
	for row := g.rows - 1; row >= 0; row-- {
		full := true
		for col := int32(0); col < g.cols; col++ {
			if g.board[row][col] == "" {
				full = false
				break
			}
		}
		if full {
			g.board = append(g.board[:row], g.board[row+1:]...)
			g.board = append([][]string{make([]string, g.cols)}, g.board...)
			lines++
			row++
		}
	}
	if lines > 0 {
		points := []int32{0, 100, 300, 500, 800}
		if lines <= 4 {
			g.score += points[lines]
		} else {
			g.score += 800
		}
	}
}

func (g *TetrisGame) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Reset()
		return nil
	}
	if g.gameOver {
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && !g.checkCollision(g.pieceShape, g.pieceX-1, g.pieceY) {
		g.pieceX--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && !g.checkCollision(g.pieceShape, g.pieceX+1, g.pieceY) {
		g.pieceX++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) && !g.checkCollision(g.pieceShape, g.pieceX, g.pieceY+1) {
		g.pieceY++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.rotatePiece()
	}

	g.fallTimer += 1.0 / 60.0
	if g.fallTimer >= g.fallDelay {
		g.fallTimer = 0
		if !g.checkCollision(g.pieceShape, g.pieceX, g.pieceY+1) {
			g.pieceY++
		} else {
			g.lockPiece()
		}
	}
	return nil
}

func (g *TetrisGame) Draw(screen *ebiten.Image) {
	screen.Fill(ColorBg)

	DrawRoundedRect(screen, float32(g.offsetX), float32(g.offsetY),
		float32(g.width), float32(g.height), 12, color.RGBA{20, 20, 30, 255})

	for y, row := range g.board {
		for x, cell := range row {
			if cell != "" {
				img := g.tileCache.Get(g.colors[cell])
				DrawImageAt(screen, img,
					float32(g.offsetX+int32(x)*g.cellSize+1),
					float32(g.offsetY+int32(y)*g.cellSize+1))
			}
		}
	}

	if !g.gameOver {
		for y, row := range g.pieceShape {
			for x, cell := range row {
				if cell == 1 {
					img := g.tileCache.Get(g.colors[g.current])
					DrawImageAt(screen, img,
						float32(g.offsetX+(g.pieceX+int32(x))*g.cellSize+1),
						float32(g.offsetY+(g.pieceY+int32(y))*g.cellSize+1))
				}
			}
		}
	}

	if FontSmall != nil {
		text.Draw(screen, "Score: "+fmt.Sprint(g.score), FontSmall,
			int(g.offsetX+g.width)+20, int(g.offsetY)+40, ColorWhite)
	}

	if g.gameOver {
		DrawOverlay(screen, int(g.offsetX), int(g.offsetY), int(g.width), int(g.height))
		msg := "GAME OVER"
		b := text.BoundString(FontBig, msg)
		text.Draw(screen, msg, FontBig,
			int(g.offsetX)+int(g.width)/2-b.Dx()/2,
			int(g.offsetY)+int(g.height)/2-20, ColorWhite)
		msg2 := "Press R to restart"
		b2 := text.BoundString(FontSmall, msg2)
		text.Draw(screen, msg2, FontSmall,
			int(g.offsetX)+int(g.width)/2-b2.Dx()/2,
			int(g.offsetY)+int(g.height)/2+30, ColorWhite)
	}
}

func (g *TetrisGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}