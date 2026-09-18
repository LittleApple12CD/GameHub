package games

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Brick struct {
	Rect  sdl.Rect
	Color sdl.Color
	Alive bool
}

type BreakoutGame struct {
	renderer   *sdl.Renderer
	width      int32
	height     int32
	offsetX    int32
	offsetY    int32
	font       *ttf.Font
	fontBig    *ttf.Font
	paddle     sdl.Rect
	paddleSpd  int32
	ball       sdl.Rect
	ballDx     int32
	ballDy     int32
	ballSpeed  int32
	bricks     []Brick
	score      int32
	lives      int32
	gameOver   bool
	waiting    bool
	clock      uint32
	running    bool
}

func NewBreakoutGame(renderer *sdl.Renderer) *BreakoutGame {
	g := &BreakoutGame{
		renderer:  renderer,
		width:     600,
		height:    600,
		paddleSpd: 9,
		ballSpeed: 6,
	}
	winW, winH, _ := renderer.GetOutputSize()
	g.offsetX = (int32(winW) - g.width) / 2
	g.offsetY = (int32(winH) - g.height) / 2
	g.font, _ = ttf.OpenFont("assets/fonts/arial.ttf", 36)
	if g.font == nil {
		g.font, _ = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 36)
	}
	g.fontBig, _ = ttf.OpenFont("assets/fonts/arial.ttf", 48)
	if g.fontBig == nil {
		g.fontBig, _ = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 48)
	}
	g.Reset()
	return g
}

func (g *BreakoutGame) Reset() {
	g.paddle = sdl.Rect{X: g.width/2 - 60, Y: g.height - 40, W: 120, H: 16}
	g.ball = sdl.Rect{X: g.width/2 - 10, Y: g.height - 70, W: 20, H: 20}
	g.ballDx = 5
	g.ballDy = -6
	g.bricks = nil
	g.score = 0
	g.lives = 3
	g.gameOver = false
	g.waiting = true
	g.running = true

	rows, cols := 5, 8
	brickWidth := (g.width - 20) / int32(cols) - 4
	brickHeight := int32(22)
	colors := []sdl.Color{
		{R: 230, G: 50, B: 50, A: 255},
		{R: 230, G: 150, B: 50, A: 255},
		{R: 230, G: 230, B: 50, A: 255},
		{R: 50, G: 230, B: 50, A: 255},
		{R: 50, G: 150, B: 230, A: 255},
	}
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			x := 10 + int32(col)*(brickWidth+4)
			y := 40 + int32(row)*(brickHeight+4)
			g.bricks = append(g.bricks, Brick{
				Rect:  sdl.Rect{X: x, Y: y, W: brickWidth, H: brickHeight},
				Color: colors[row%len(colors)],
				Alive: true,
			})
		}
	}
}

func rectCollide(a, b *sdl.Rect) bool {
	return a.X < b.X+b.W && a.X+a.W > b.X &&
		a.Y < b.Y+b.H && a.Y+a.H > b.Y
}

func (g *BreakoutGame) HandleEvents() bool {
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
				case sdl.K_SPACE:
					if g.waiting && !g.gameOver {
						g.waiting = false
					}
				}
			}
		}
	}

	if !g.gameOver && !g.waiting {
		keys := sdl.GetKeyboardState()
		if keys[sdl.SCANCODE_LEFT] != 0 && g.paddle.X > 0 {
			g.paddle.X -= g.paddleSpd
		}
		if keys[sdl.SCANCODE_RIGHT] != 0 && g.paddle.X+g.paddle.W < g.width {
			g.paddle.X += g.paddleSpd
		}
	}
	return true
}

func (g *BreakoutGame) Update() {
	if g.gameOver || g.waiting || !g.running {
		return
	}

	g.ball.X += g.ballDx
	g.ball.Y += g.ballDy

	if g.ball.X <= 0 || g.ball.X+g.ball.W >= g.width {
		g.ballDx = -g.ballDx
	}
	if g.ball.Y <= 0 {
		g.ballDy = -g.ballDy
	}

	if g.ball.Y+g.ball.H >= g.height {
		g.lives--
		if g.lives <= 0 {
			g.gameOver = true
		} else {
			g.waiting = true
			g.ball.X = g.width/2 - 10
			g.ball.Y = g.height - 70
			rand.Seed(time.Now().UnixNano())
			if rand.Intn(2) == 0 {
				g.ballDx = g.ballSpeed
			} else {
				g.ballDx = -g.ballSpeed
			}
			g.ballDy = -g.ballSpeed
		}
		return
	}

	if rectCollide(&g.ball, &g.paddle) {
		g.ballDy = -abs32(g.ballDy)
		hitPos := float32(g.ball.X+g.ball.W/2-g.paddle.X-g.paddle.W/2) / float32(g.paddle.W/2)
		g.ballDx = int32(hitPos * float32(g.ballSpeed) * 0.9)
		if abs32(g.ballDx) < 2 {
			if g.ballDx >= 0 {
				g.ballDx = 3
			} else {
				g.ballDx = -3
			}
		}
	}

	for i := range g.bricks {
		if !g.bricks[i].Alive {
			continue
		}
		if rectCollide(&g.ball, &g.bricks[i].Rect) {
			g.bricks[i].Alive = false
			g.score += 10

			overlapTop := g.ball.Y + g.ball.H - g.bricks[i].Rect.Y
			overlapBottom := g.bricks[i].Rect.Y + g.bricks[i].Rect.H - g.ball.Y
			overlapLeft := g.ball.X + g.ball.W - g.bricks[i].Rect.X
			overlapRight := g.bricks[i].Rect.X + g.bricks[i].Rect.W - g.ball.X

			minOverlap := overlapTop
			if overlapBottom < minOverlap {
				minOverlap = overlapBottom
			}
			if overlapLeft < minOverlap {
				minOverlap = overlapLeft
			}
			if overlapRight < minOverlap {
				minOverlap = overlapRight
			}

			if minOverlap == overlapTop || minOverlap == overlapBottom {
				g.ballDy = -g.ballDy
			} else {
				g.ballDx = -g.ballDx
			}
			break
		}
	}

	allDead := true
	for _, b := range g.bricks {
		if b.Alive {
			allDead = false
			break
		}
	}
	if allDead {
		g.gameOver = true
	}
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func (g *BreakoutGame) Draw() {
	winW, winH, _ := g.renderer.GetOutputSize()
	g.offsetX = (int32(winW) - g.width) / 2
	g.offsetY = (int32(winH) - g.height) / 2

	g.renderer.SetDrawColor(25, 25, 35, 255)
	g.renderer.Clear()

	g.renderer.SetDrawColor(20, 20, 30, 255)
	DrawRoundedRect(g.renderer, g.offsetX, g.offsetY, g.width, g.height, 12)

	for _, brick := range g.bricks {
		if brick.Alive {
			rect := sdl.Rect{
				X: g.offsetX + brick.Rect.X,
				Y: g.offsetY + brick.Rect.Y,
				W: brick.Rect.W,
				H: brick.Rect.H,
			}
			g.renderer.SetDrawColor(brick.Color.R, brick.Color.G, brick.Color.B, brick.Color.A)
			DrawRoundedRect(g.renderer, rect.X, rect.Y, rect.W, rect.H, 4)
			g.renderer.SetDrawColor(255, 255, 255, 100)
			DrawRoundedRectBorder(g.renderer, rect.X, rect.Y, rect.W, rect.H, 4, 1)
		}
	}

	paddleRect := sdl.Rect{
		X: g.offsetX + g.paddle.X,
		Y: g.offsetY + g.paddle.Y,
		W: g.paddle.W,
		H: g.paddle.H,
	}
	g.renderer.SetDrawColor(255, 255, 255, 255)
	DrawRoundedRect(g.renderer, paddleRect.X, paddleRect.Y, paddleRect.W, paddleRect.H, 10)
	g.renderer.SetDrawColor(200, 200, 220, 255)
	DrawRoundedRectBorder(g.renderer, paddleRect.X, paddleRect.Y, paddleRect.W, paddleRect.H, 10, 1)

	ballRect := sdl.Rect{
		X: g.offsetX + g.ball.X,
		Y: g.offsetY + g.ball.Y,
		W: g.ball.W,
		H: g.ball.H,
	}
	g.renderer.SetDrawColor(255, 255, 120, 255)
	DrawRoundedRect(g.renderer, ballRect.X, ballRect.Y, ballRect.W, ballRect.H, 10)

	if g.font != nil {
		scoreText := "Score: " + strconv.Itoa(int(g.score))
		surface, err := g.font.RenderUTF8Solid(scoreText, sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err == nil {
			defer surface.Free()
			texture, err := g.renderer.CreateTextureFromSurface(surface)
			if err == nil {
				defer texture.Destroy()
				_, _, w, h, _ := texture.Query()
				g.renderer.Copy(texture, nil, &sdl.Rect{X: g.offsetX + 10, Y: g.offsetY + 10, W: w, H: h})
			}
		}
		livesText := "Lives: " + strconv.Itoa(int(g.lives))
		surface2, err := g.font.RenderUTF8Solid(livesText, sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err == nil {
			defer surface2.Free()
			texture2, err := g.renderer.CreateTextureFromSurface(surface2)
			if err == nil {
				defer texture2.Destroy()
				_, _, w, h, _ := texture2.Query()
				g.renderer.Copy(texture2, nil, &sdl.Rect{X: g.offsetX + g.width - 120, Y: g.offsetY + 10, W: w, H: h})
			}
		}
	}

	if g.waiting && !g.gameOver {
		if g.fontBig != nil {
			surface, err := g.fontBig.RenderUTF8Solid("Press SPACE", sdl.Color{R: 255, G: 255, B: 255, A: 255})
			if err == nil {
				defer surface.Free()
				texture, err := g.renderer.CreateTextureFromSurface(surface)
				if err == nil {
					defer texture.Destroy()
					_, _, w, h, _ := texture.Query()
					x := g.offsetX + g.width/2 - int32(w)/2
					y := g.offsetY + g.height/2 + 30 - int32(h)/2
					g.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: w, H: h})
				}
			}
		}
		if g.font != nil {
			surface, err := g.font.RenderUTF8Solid("to start", sdl.Color{R: 255, G: 255, B: 255, A: 255})
			if err == nil {
				defer surface.Free()
				texture, err := g.renderer.CreateTextureFromSurface(surface)
				if err == nil {
					defer texture.Destroy()
					_, _, w, h, _ := texture.Query()
					x := g.offsetX + g.width/2 - int32(w)/2
					y := g.offsetY + g.height/2 + 70 - int32(h)/2
					g.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: w, H: h})
				}
			}
		}
	} else if g.gameOver {
		overlay := sdl.Rect{X: g.offsetX, Y: g.offsetY, W: g.width, H: g.height}
		g.renderer.SetDrawColor(0, 0, 0, 150)
		g.renderer.FillRect(&overlay)

		var text string
		var color sdl.Color
		if g.lives <= 0 {
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

func (g *BreakoutGame) Run() {
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