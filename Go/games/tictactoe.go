package games

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type TicTacToeGame struct {
	renderer  *sdl.Renderer
	width     int32
	height    int32
	offsetX   int32
	offsetY   int32
	cellSize  int32
	font      *ttf.Font
	fontBig   *ttf.Font
	fontSmall *ttf.Font
	board     [3][3]string
	current   string
	winner    string
	gameOver  bool
	moveCount int32
	running   bool
}

func NewTicTacToeGame(renderer *sdl.Renderer) *TicTacToeGame {
	g := &TicTacToeGame{
		renderer: renderer,
		width:    500,
		height:   500,
		current:  "X",
	}
	g.cellSize = g.width / 3
	winW, winH, _ := renderer.GetOutputSize()
	g.offsetX = (int32(winW) - g.width) / 2
	g.offsetY = (int32(winH) - g.height) / 2
	g.font, _ = ttf.OpenFont("assets/fonts/arial.ttf", 80)
	if g.font == nil {
		g.font, _ = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 80)
	}
	g.fontBig, _ = ttf.OpenFont("assets/fonts/arial.ttf", 48)
	if g.fontBig == nil {
		g.fontBig, _ = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 48)
	}
	g.fontSmall, _ = ttf.OpenFont("assets/fonts/arial.ttf", 28)
	if g.fontSmall == nil {
		g.fontSmall, _ = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 28)
	}
	g.Reset()
	return g
}

func (g *TicTacToeGame) Reset() {
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			g.board[r][c] = ""
		}
	}
	g.current = "X"
	g.winner = ""
	g.gameOver = false
	g.moveCount = 0
	g.running = true
}

func (g *TicTacToeGame) checkWinner() string {
	for row := 0; row < 3; row++ {
		if g.board[row][0] != "" && g.board[row][0] == g.board[row][1] && g.board[row][1] == g.board[row][2] {
			return g.board[row][0]
		}
	}
	for col := 0; col < 3; col++ {
		if g.board[0][col] != "" && g.board[0][col] == g.board[1][col] && g.board[1][col] == g.board[2][col] {
			return g.board[0][col]
		}
	}
	if g.board[0][0] != "" && g.board[0][0] == g.board[1][1] && g.board[1][1] == g.board[2][2] {
		return g.board[0][0]
	}
	if g.board[0][2] != "" && g.board[0][2] == g.board[1][1] && g.board[1][1] == g.board[2][0] {
		return g.board[0][2]
	}
	return ""
}

func (g *TicTacToeGame) HandleEvents() bool {
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
			if e.Type == sdl.MOUSEBUTTONDOWN && e.Button == sdl.BUTTON_LEFT && !g.gameOver {
				x := (e.X - g.offsetX) / g.cellSize
				y := (e.Y - g.offsetY) / g.cellSize
				if x >= 0 && x < 3 && y >= 0 && y < 3 && g.board[y][x] == "" {
					g.board[y][x] = g.current
					g.moveCount++
					g.winner = g.checkWinner()
					if g.winner != "" {
						g.gameOver = true
					} else if g.moveCount == 9 {
						g.gameOver = true
					} else {
						if g.current == "X" {
							g.current = "O"
						} else {
							g.current = "X"
						}
					}
				}
			}
		}
	}
	return true
}

func (g *TicTacToeGame) Update() {}

func (g *TicTacToeGame) Draw() {
	winW, winH, _ := g.renderer.GetOutputSize()
	g.offsetX = (int32(winW) - g.width) / 2
	g.offsetY = (int32(winH) - g.height) / 2

	g.renderer.SetDrawColor(25, 25, 35, 255)
	g.renderer.Clear()

	g.renderer.SetDrawColor(30, 30, 42, 255)
	DrawRoundedRect(g.renderer, g.offsetX, g.offsetY, g.width, g.height, 12)

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			rect := sdl.Rect{
				X: g.offsetX + int32(col)*g.cellSize + 2,
				Y: g.offsetY + int32(row)*g.cellSize + 2,
				W: g.cellSize - 4,
				H: g.cellSize - 4,
			}
			g.renderer.SetDrawColor(50, 50, 65, 255)
			DrawRoundedRect(g.renderer, rect.X, rect.Y, rect.W, rect.H, 8)

			if g.board[row][col] != "" {
				var color sdl.Color
				if g.board[row][col] == "X" {
					color = sdl.Color{R: 60, G: 200, B: 255, A: 255}
				} else {
					color = sdl.Color{R: 255, G: 200, B: 60, A: 255}
				}
				if g.font != nil {
					surface, err := g.font.RenderUTF8Solid(g.board[row][col], color)
					if err == nil {
						defer surface.Free()
						texture, err := g.renderer.CreateTextureFromSurface(surface)
						if err == nil {
							defer texture.Destroy()
							_, _, w, h, _ := texture.Query()
							x := rect.X + rect.W/2 - int32(w)/2
							y := rect.Y + rect.H/2 - int32(h)/2
							g.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: w, H: h})
						}
					}
				}
			}
		}
	}

	for i := int32(1); i < 3; i++ {
		g.renderer.SetDrawColor(80, 80, 100, 255)
		g.renderer.FillRect(&sdl.Rect{
			X: g.offsetX + i*g.cellSize - 1,
			Y: g.offsetY + 5,
			W: 3,
			H: g.height - 10,
		})
		g.renderer.FillRect(&sdl.Rect{
			X: g.offsetX + 5,
			Y: g.offsetY + i*g.cellSize - 1,
			W: g.width - 10,
			H: 3,
		})
	}

	statusY := g.offsetY + g.height + 20
	if g.gameOver {
		var text string
		var color sdl.Color
		if g.winner != "" {
			text = "Player " + g.winner + " Wins!"
			color = sdl.Color{R: 0, G: 255, B: 0, A: 255}
		} else {
			text = "Draw!"
			color = sdl.Color{R: 255, G: 255, B: 100, A: 255}
		}
		if g.fontBig != nil {
			surface, err := g.fontBig.RenderUTF8Solid(text, color)
			if err == nil {
				defer surface.Free()
				texture, err := g.renderer.CreateTextureFromSurface(surface)
				if err == nil {
					defer texture.Destroy()
					_, _, w, h, _ := texture.Query()
					x := int32(winW/2) - int32(w)/2
					g.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: statusY, W: w, H: h})
				}
			}
		}
		if g.fontSmall != nil {
			surface, err := g.fontSmall.RenderUTF8Solid("Press R to restart", sdl.Color{R: 200, G: 200, B: 200, A: 255})
			if err == nil {
				defer surface.Free()
				texture, err := g.renderer.CreateTextureFromSurface(surface)
				if err == nil {
					defer texture.Destroy()
					_, _, w, h, _ := texture.Query()
					x := int32(winW/2) - int32(w)/2
					g.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: statusY + 40, W: w, H: h})
				}
			}
		}
	} else {
		if g.fontBig != nil {
			text := "Player " + g.current + "'s Turn"
			surface, err := g.fontBig.RenderUTF8Solid(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
			if err == nil {
				defer surface.Free()
				texture, err := g.renderer.CreateTextureFromSurface(surface)
				if err == nil {
					defer texture.Destroy()
					_, _, w, h, _ := texture.Query()
					x := int32(winW/2) - int32(w)/2
					g.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: statusY, W: w, H: h})
				}
			}
		}
	}

	g.renderer.Present()
}

func (g *TicTacToeGame) Run() {
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