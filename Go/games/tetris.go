package games

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type TetrisGame struct {
	renderer     *sdl.Renderer
	width        int32
	height       int32
	offsetX      int32
	offsetY      int32
	cols         int32
	rows         int32
	cellSize     int32
	font         *ttf.Font
	fontBig      *ttf.Font
	board        [][]string
	score        int32
	gameOver     bool
	fallTimer    uint32
	fallDelay    uint32
	clock        uint32
	currentPiece string
	pieceShape   [][]int
	pieceX       int32
	pieceY       int32
	shapes       map[string][][]int
	colors       map[string]sdl.Color
	running      bool
}

func NewTetrisGame(renderer *sdl.Renderer) *TetrisGame {
	g := &TetrisGame{
		renderer:  renderer,
		width:     300,
		height:    600,
		cols:      10,
		rows:      20,
		cellSize:  30,
		fallDelay: 500,
	}
	winW, winH, _ := renderer.GetOutputSize()
	g.offsetX = (int32(winW) - g.width - 140) / 2
	g.offsetY = (int32(winH) - g.height) / 2
	g.font, _ = ttf.OpenFont("assets/fonts/arial.ttf", 32)
	if g.font == nil {
		g.font, _ = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 32)
	}
	g.fontBig, _ = ttf.OpenFont("assets/fonts/arial.ttf", 48)
	if g.fontBig == nil {
		g.fontBig, _ = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 48)
	}
	g.shapes = map[string][][]int{
		"I": {{0, 0, 0, 0}, {1, 1, 1, 1}, {0, 0, 0, 0}, {0, 0, 0, 0}},
		"O": {{1, 1}, {1, 1}},
		"T": {{0, 1, 0}, {1, 1, 1}, {0, 0, 0}},
		"S": {{0, 1, 1}, {1, 1, 0}, {0, 0, 0}},
		"Z": {{1, 1, 0}, {0, 1, 1}, {0, 0, 0}},
		"L": {{1, 0, 0}, {1, 1, 1}, {0, 0, 0}},
		"J": {{0, 0, 1}, {1, 1, 1}, {0, 0, 0}},
	}
	g.colors = map[string]sdl.Color{
		"I": {R: 0, G: 240, B: 240, A: 255},
		"O": {R: 240, G: 240, B: 0, A: 255},
		"T": {R: 180, G: 0, B: 240, A: 255},
		"S": {R: 0, G: 240, B: 0, A: 255},
		"Z": {R: 240, G: 0, B: 0, A: 255},
		"L": {R: 240, G: 160, B: 0, A: 255},
		"J": {R: 0, G: 0, B: 240, A: 255},
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
	g.running = true
}

func (g *TetrisGame) spawnPiece() {
	rand.Seed(time.Now().UnixNano())
	keys := make([]string, 0, len(g.shapes))
	for k := range g.shapes {
		keys = append(keys, k)
	}
	g.currentPiece = keys[rand.Intn(len(keys))]
	shape := g.shapes[g.currentPiece]
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
					g.board[by][bx] = g.currentPiece
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

func (g *TetrisGame) HandleEvents() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			g.running = false
			return false
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				switch e.Keysym.Sym {
				case sdl.K_ESCAPE:
					g.running = false
					return false
				case sdl.K_r:
					g.Reset()
				case sdl.K_LEFT:
					if !g.gameOver && !g.checkCollision(g.pieceShape, g.pieceX-1, g.pieceY) {
						g.pieceX--
					}
				case sdl.K_RIGHT:
					if !g.gameOver && !g.checkCollision(g.pieceShape, g.pieceX+1, g.pieceY) {
						g.pieceX++
					}
				case sdl.K_DOWN:
					if !g.gameOver && !g.checkCollision(g.pieceShape, g.pieceX, g.pieceY+1) {
						g.pieceY++
					}
				case sdl.K_SPACE, sdl.K_UP:
					if !g.gameOver {
						g.rotatePiece()
					}
				}
			}
		}
	}
	return true
}

func (g *TetrisGame) Update() {
	if g.gameOver || !g.running {
		return
	}
	now := uint32(sdl.GetTicks())
	dt := now - g.clock
	g.clock = now
	g.fallTimer += dt
	if g.fallTimer >= g.fallDelay {
		g.fallTimer = 0
		if !g.checkCollision(g.pieceShape, g.pieceX, g.pieceY+1) {
			g.pieceY++
		} else {
			g.lockPiece()
		}
	}
}

func (g *TetrisGame) Draw() {
	winW, winH, _ := g.renderer.GetOutputSize()
	g.offsetX = (int32(winW) - g.width - 140) / 2
	g.offsetY = (int32(winH) - g.height) / 2

	g.renderer.SetDrawColor(25, 25, 35, 255)
	g.renderer.Clear()

	g.renderer.SetDrawColor(20, 20, 30, 255)
	DrawRoundedRect(g.renderer, g.offsetX, g.offsetY, g.width, g.height, 12)

	for y, row := range g.board {
		for x, cell := range row {
			if cell != "" {
				color := g.colors[cell]
				rect := sdl.Rect{
					X: g.offsetX + int32(x)*g.cellSize + 1,
					Y: g.offsetY + int32(y)*g.cellSize + 1,
					W: g.cellSize - 2,
					H: g.cellSize - 2,
				}
				g.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
				DrawRoundedRect(g.renderer, rect.X, rect.Y, rect.W, rect.H, 4)
			}
		}
	}

	if !g.gameOver {
		for y, row := range g.pieceShape {
			for x, cell := range row {
				if cell == 1 {
					color := g.colors[g.currentPiece]
					rect := sdl.Rect{
						X: g.offsetX + (g.pieceX+int32(x))*g.cellSize + 1,
						Y: g.offsetY + (g.pieceY+int32(y))*g.cellSize + 1,
						W: g.cellSize - 2,
						H: g.cellSize - 2,
					}
					g.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
					DrawRoundedRect(g.renderer, rect.X, rect.Y, rect.W, rect.H, 4)
				}
			}
		}
	}

	if g.font != nil {
		text := "Score: " + strconv.Itoa(int(g.score))
		surface, err := g.font.RenderUTF8Solid(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err == nil {
			defer surface.Free()
			texture, err := g.renderer.CreateTextureFromSurface(surface)
			if err == nil {
				defer texture.Destroy()
				_, _, w, h, _ := texture.Query()
				g.renderer.Copy(texture, nil, &sdl.Rect{X: g.offsetX + g.width + 20, Y: g.offsetY + 20, W: w, H: h})
			}
		}
	}

	if g.gameOver {
		overlay := sdl.Rect{X: g.offsetX, Y: g.offsetY, W: g.width, H: g.height}
		g.renderer.SetDrawColor(0, 0, 0, 180)
		g.renderer.FillRect(&overlay)
		if g.fontBig != nil {
			surface, err := g.fontBig.RenderUTF8Solid("GAME OVER", sdl.Color{R: 255, G: 255, B: 255, A: 255})
			if err == nil {
				defer surface.Free()
				texture, err := g.renderer.CreateTextureFromSurface(surface)
				if err == nil {
					defer texture.Destroy()
					_, _, w, h, _ := texture.Query()
					x := g.offsetX + g.width/2 - int32(w)/2
					y := g.offsetY + g.height/2 - int32(h)/2 - 20
					g.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: w, H: h})
				}
			}
		}
		if g.font != nil {
			surface, err := g.font.RenderUTF8Solid("Press R to restart", sdl.Color{R: 255, G: 255, B: 255, A: 255})
			if err == nil {
				defer surface.Free()
				texture, err := g.renderer.CreateTextureFromSurface(surface)
				if err == nil {
					defer texture.Destroy()
					_, _, w, h, _ := texture.Query()
					x := g.offsetX + g.width/2 - int32(w)/2
					y := g.offsetY + g.height/2 + 30
					g.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: w, H: h})
				}
			}
		}
	}

	g.renderer.Present()
}

func (g *TetrisGame) Run() {
	g.running = true
	g.clock = uint32(sdl.GetTicks())
	for g.running {
		if !g.HandleEvents() {
			break
		}
		g.Update()
		g.Draw()
		sdl.Delay(16)
	}
}