package games

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type MinesweeperGame struct {
	renderer   *sdl.Renderer
	rows       int32
	cols       int32
	mines      int32
	cellSize   int32
	width      int32
	height     int32
	offsetX    int32
	offsetY    int32
	font       *ttf.Font
	fontBig    *ttf.Font
	board      [][]int
	revealed   [][]bool
	flagged    [][]bool
	gameOver   bool
	won        bool
	firstClick bool
	mineCount  int32
	running    bool
}

func NewMinesweeperGame(renderer *sdl.Renderer) *MinesweeperGame {
	g := &MinesweeperGame{
		renderer: renderer,
		rows:     16,
		cols:     16,
		mines:    40,
		cellSize: 30,
	}
	g.width = g.cols * g.cellSize
	g.height = g.rows * g.cellSize
	winW, winH, _ := renderer.GetOutputSize()
	g.offsetX = (int32(winW) - g.width) / 2
	g.offsetY = (int32(winH) - g.height) / 2
	g.font, _ = ttf.OpenFont("assets/fonts/arial.ttf", 22)
	if g.font == nil {
		g.font, _ = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 22)
	}
	g.fontBig, _ = ttf.OpenFont("assets/fonts/arial.ttf", 40)
	if g.fontBig == nil {
		g.fontBig, _ = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 40)
	}
	g.Reset()
	return g
}

func (g *MinesweeperGame) Reset() {
	g.board = make([][]int, g.rows)
	for i := range g.board {
		g.board[i] = make([]int, g.cols)
	}
	g.revealed = make([][]bool, g.rows)
	for i := range g.revealed {
		g.revealed[i] = make([]bool, g.cols)
	}
	g.flagged = make([][]bool, g.rows)
	for i := range g.flagged {
		g.flagged[i] = make([]bool, g.cols)
	}
	g.gameOver = false
	g.won = false
	g.firstClick = true
	g.mineCount = g.mines
	g.running = true
}

func (g *MinesweeperGame) placeMines(safeX, safeY int32) {
	rand.Seed(time.Now().UnixNano())
	placed := int32(0)
	for placed < g.mines {
		x := int32(rand.Intn(int(g.cols)))
		y := int32(rand.Intn(int(g.rows)))
		if g.board[y][x] == -1 {
			continue
		}
		if x >= safeX-1 && x <= safeX+1 && y >= safeY-1 && y <= safeY+1 {
			continue
		}
		g.board[y][x] = -1
		placed++
	}
	for y := int32(0); y < g.rows; y++ {
		for x := int32(0); x < g.cols; x++ {
			if g.board[y][x] == -1 {
				continue
			}
			count := 0
			for dy := int32(-1); dy <= 1; dy++ {
				for dx := int32(-1); dx <= 1; dx++ {
					nx, ny := x+dx, y+dy
					if nx >= 0 && nx < g.cols && ny >= 0 && ny < g.rows {
						if g.board[ny][nx] == -1 {
							count++
						}
					}
				}
			}
			g.board[y][x] = count
		}
	}
}

func (g *MinesweeperGame) reveal(x, y int32) {
	if x < 0 || x >= g.cols || y < 0 || y >= g.rows {
		return
	}
	if g.revealed[y][x] || g.flagged[y][x] {
		return
	}
	g.revealed[y][x] = true
	if g.board[y][x] == -1 {
		g.gameOver = true
		return
	}
	if g.board[y][x] == 0 {
		for dy := int32(-1); dy <= 1; dy++ {
			for dx := int32(-1); dx <= 1; dx++ {
				g.reveal(x+dx, y+dy)
			}
		}
	}
	revealedCount := int32(0)
	for y2 := int32(0); y2 < g.rows; y2++ {
		for x2 := int32(0); x2 < g.cols; x2++ {
			if g.revealed[y2][x2] {
				revealedCount++
			}
		}
	}
	if revealedCount == g.rows*g.cols-g.mines {
		g.won = true
	}
}

func (g *MinesweeperGame) HandleEvents() bool {
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
				}
			}
		case *sdl.MouseButtonEvent:
			if e.Type == sdl.MOUSEBUTTONDOWN && !g.gameOver && !g.won {
				x := (e.X - g.offsetX) / g.cellSize
				y := (e.Y - g.offsetY) / g.cellSize
				if x < 0 || x >= g.cols || y < 0 || y >= g.rows {
					continue
				}
				if e.Button == sdl.BUTTON_LEFT {
					if g.firstClick {
						g.placeMines(x, y)
						g.firstClick = false
					}
					if !g.flagged[y][x] {
						g.reveal(x, y)
					}
				} else if e.Button == sdl.BUTTON_RIGHT {
					if !g.revealed[y][x] {
						g.flagged[y][x] = !g.flagged[y][x]
					}
				}
			}
		}
	}
	return true
}

func (g *MinesweeperGame) Update() {}

func (g *MinesweeperGame) Draw() {
	winW, winH, _ := g.renderer.GetOutputSize()
	g.offsetX = (int32(winW) - g.width) / 2
	g.offsetY = (int32(winH) - g.height) / 2

	g.renderer.SetDrawColor(25, 25, 35, 255)
	g.renderer.Clear()

	g.renderer.SetDrawColor(30, 30, 42, 255)
	DrawRoundedRect(g.renderer, g.offsetX, g.offsetY, g.width, g.height, 12)

	colors := []sdl.Color{
		{R: 0, G: 0, B: 0, A: 255},
		{R: 0, G: 0, B: 255, A: 255},
		{R: 0, G: 150, B: 0, A: 255},
		{R: 255, G: 0, B: 0, A: 255},
		{R: 0, G: 0, B: 180, A: 255},
		{R: 150, G: 0, B: 0, A: 255},
		{R: 0, G: 150, B: 150, A: 255},
		{R: 0, G: 0, B: 0, A: 255},
	}

	for y := int32(0); y < g.rows; y++ {
		for x := int32(0); x < g.cols; x++ {
			rect := sdl.Rect{
				X: g.offsetX + x*g.cellSize + 1,
				Y: g.offsetY + y*g.cellSize + 1,
				W: g.cellSize - 2,
				H: g.cellSize - 2,
			}
			if g.revealed[y][x] {
				if g.board[y][x] == -1 {
					g.renderer.SetDrawColor(200, 50, 50, 255)
					DrawRoundedRect(g.renderer, rect.X, rect.Y, rect.W, rect.H, 4)
					centerX := rect.X + g.cellSize/2
					centerY := rect.Y + g.cellSize/2
					g.renderer.SetDrawColor(0, 0, 0, 255)
					for i := int32(-8); i <= 8; i++ {
						g.renderer.DrawPoint(centerX+i, centerY)
						g.renderer.DrawPoint(centerX, centerY+i)
					}
				} else {
					g.renderer.SetDrawColor(190, 190, 200, 255)
					DrawRoundedRect(g.renderer, rect.X, rect.Y, rect.W, rect.H, 4)
					if g.board[y][x] > 0 {
						if g.font != nil {
							text := strconv.Itoa(g.board[y][x])
							surface, err := g.font.RenderUTF8Solid(text, colors[g.board[y][x]])
							if err == nil {
								defer surface.Free()
								texture, err := g.renderer.CreateTextureFromSurface(surface)
								if err == nil {
									defer texture.Destroy()
									_, _, w, h, _ := texture.Query()
									g.renderer.Copy(texture, nil, &sdl.Rect{
										X: rect.X + g.cellSize/2 - int32(w)/2,
										Y: rect.Y + g.cellSize/2 - int32(h)/2,
										W: w, H: h,
									})
								}
							}
						}
					}
				}
			} else {
				if g.flagged[y][x] {
					g.renderer.SetDrawColor(210, 210, 60, 255)
				} else {
					g.renderer.SetDrawColor(80, 80, 105, 255)
				}
				DrawRoundedRect(g.renderer, rect.X, rect.Y, rect.W, rect.H, 4)
				if g.flagged[y][x] && g.font != nil {
					surface, err := g.font.RenderUTF8Solid("F", sdl.Color{R: 255, G: 50, B: 50, A: 255})
					if err == nil {
						defer surface.Free()
						texture, err := g.renderer.CreateTextureFromSurface(surface)
						if err == nil {
							defer texture.Destroy()
							_, _, w, h, _ := texture.Query()
							g.renderer.Copy(texture, nil, &sdl.Rect{
								X: rect.X + g.cellSize/2 - int32(w)/2,
								Y: rect.Y + g.cellSize/2 - int32(h)/2,
								W: w, H: h,
							})
						}
					}
				}
			}
		}
	}

	if g.gameOver || g.won {
		overlay := sdl.Rect{X: g.offsetX, Y: g.offsetY, W: g.width, H: g.height}
		g.renderer.SetDrawColor(0, 0, 0, 150)
		g.renderer.FillRect(&overlay)
		var text string
		var color sdl.Color
		if g.gameOver {
			text = "GAME OVER"
			color = sdl.Color{R: 255, G: 255, B: 255, A: 255}
		} else {
			text = "YOU WIN!"
			color = sdl.Color{R: 0, G: 255, B: 0, A: 255}
		}
		if g.fontBig != nil {
			surface, err := g.fontBig.RenderUTF8Solid(text, color)
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

func (g *MinesweeperGame) Run() {
	g.running = true
	for g.running {
		if !g.HandleEvents() {
			break
		}
		g.Update()
		g.Draw()
		sdl.Delay(16)
	}
}