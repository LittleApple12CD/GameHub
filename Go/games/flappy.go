package games

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Pipe struct {
	X      int32
	GapY   int32
	Passed bool
}

type FlappyGame struct {
	renderer   *sdl.Renderer
	width      int32
	height     int32
	offsetX    int32
	offsetY    int32
	font       *ttf.Font
	fontSmall  *ttf.Font
	birdY      float32
	birdVY     float32
	birdRadius int32
	birdX      int32
	gravity    float32
	jumpStr    float32
	pipes      []Pipe
	pipeWidth  int32
	pipeGap    int32
	pipeSpeed  int32
	score      int32
	gameOver   bool
	pipeTimer  int32
	pipeDelay  int32
	running    bool
}

func NewFlappyGame(renderer *sdl.Renderer) *FlappyGame {
	g := &FlappyGame{
		renderer:   renderer,
		width:      400,
		height:     600,
		birdRadius: 15,
		birdX:      80,
		gravity:    0.6,
		jumpStr:    -11,
		pipeWidth:  60,
		pipeGap:    160,
		pipeSpeed:  3,
		pipeDelay:  90,
	}
	winW, winH, _ := renderer.GetOutputSize()
	g.offsetX = (int32(winW) - g.width) / 2
	g.offsetY = (int32(winH) - g.height) / 2
	g.font, _ = ttf.OpenFont("assets/fonts/arial.ttf", 40)
	if g.font == nil {
		g.font, _ = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 40)
	}
	g.fontSmall, _ = ttf.OpenFont("assets/fonts/arial.ttf", 26)
	if g.fontSmall == nil {
		g.fontSmall, _ = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 26)
	}
	g.Reset()
	return g
}

func (g *FlappyGame) Reset() {
	g.birdY = float32(g.height / 2)
	g.birdVY = 0
	g.pipes = nil
	g.score = 0
	g.gameOver = false
	g.pipeTimer = 0
	g.running = true
}

func (g *FlappyGame) addPipe() {
	rand.Seed(time.Now().UnixNano())
	gapY := int32(rand.Intn(int(g.height-80-g.pipeGap))) + 40
	g.pipes = append(g.pipes, Pipe{X: g.width, GapY: gapY, Passed: false})
}

func (g *FlappyGame) HandleEvents() bool {
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
					if !g.gameOver {
						g.birdVY = g.jumpStr
					}
				}
			}
		case *sdl.MouseButtonEvent:
			if e.Type == sdl.MOUSEBUTTONDOWN && !g.gameOver {
				g.birdVY = g.jumpStr
			}
		}
	}
	return true
}

func (g *FlappyGame) Update() {
	if g.gameOver || !g.running {
		return
	}

	g.birdVY += g.gravity
	g.birdY += g.birdVY

	if g.birdY-float32(g.birdRadius) < 0 {
		g.birdY = float32(g.birdRadius)
		g.birdVY = 0
	} else if g.birdY+float32(g.birdRadius) > float32(g.height) {
		g.gameOver = true
		return
	}

	g.pipeTimer++
	if g.pipeTimer >= g.pipeDelay {
		g.pipeTimer = 0
		g.addPipe()
	}

	for i := range g.pipes {
		g.pipes[i].X -= g.pipeSpeed

		if !g.pipes[i].Passed && g.pipes[i].X+g.pipeWidth < g.birdX {
			g.pipes[i].Passed = true
			g.score++
		}

		if g.pipes[i].X < g.birdX+int32(g.birdRadius) &&
			g.pipes[i].X+g.pipeWidth > g.birdX-int32(g.birdRadius) {
			if g.birdY-float32(g.birdRadius) < float32(g.pipes[i].GapY) ||
				g.birdY+float32(g.birdRadius) > float32(g.pipes[i].GapY+g.pipeGap) {
				g.gameOver = true
				return
			}
		}
	}

	newPipes := make([]Pipe, 0)
	for _, p := range g.pipes {
		if p.X+g.pipeWidth > 0 {
			newPipes = append(newPipes, p)
		}
	}
	g.pipes = newPipes
}

func (g *FlappyGame) Draw() {
	winW, winH, _ := g.renderer.GetOutputSize()
	g.offsetX = (int32(winW) - g.width) / 2
	g.offsetY = (int32(winH) - g.height) / 2

	g.renderer.SetDrawColor(25, 25, 35, 255)
	g.renderer.Clear()

	g.renderer.SetDrawColor(135, 206, 235, 255)
	DrawRoundedRect(g.renderer, g.offsetX, g.offsetY, g.width, g.height, 12)

	bx := g.offsetX + g.birdX
	by := g.offsetY + int32(g.birdY)
	g.renderer.SetDrawColor(255, 255, 50, 255)
	for dy := -g.birdRadius; dy <= g.birdRadius; dy++ {
		for dx := -g.birdRadius; dx <= g.birdRadius; dx++ {
			if dx*dx+dy*dy <= g.birdRadius*g.birdRadius {
				g.renderer.DrawPoint(bx+dx, by+dy)
			}
		}
	}
	g.renderer.SetDrawColor(0, 0, 0, 255)
	g.renderer.FillRect(&sdl.Rect{X: bx + 8, Y: by - 5, W: 4, H: 4})
	g.renderer.SetDrawColor(255, 255, 255, 255)
	g.renderer.FillRect(&sdl.Rect{X: bx + 10, Y: by - 7, W: 2, H: 2})

	for _, pipe := range g.pipes {
		g.renderer.SetDrawColor(40, 200, 40, 255)
		topRect := sdl.Rect{X: g.offsetX + pipe.X, Y: g.offsetY, W: g.pipeWidth, H: pipe.GapY}
		g.renderer.FillRect(&topRect)
		g.renderer.SetDrawColor(30, 160, 30, 255)
		capRect := sdl.Rect{X: g.offsetX + pipe.X - 6, Y: g.offsetY + pipe.GapY - 22, W: g.pipeWidth + 12, H: 22}
		DrawRoundedRect(g.renderer, capRect.X, capRect.Y, capRect.W, capRect.H, 6)

		bottomY := pipe.GapY + g.pipeGap
		g.renderer.SetDrawColor(40, 200, 40, 255)
		bottomRect := sdl.Rect{X: g.offsetX + pipe.X, Y: g.offsetY + bottomY, W: g.pipeWidth, H: g.height - bottomY}
		g.renderer.FillRect(&bottomRect)
		g.renderer.SetDrawColor(30, 160, 30, 255)
		capRect2 := sdl.Rect{X: g.offsetX + pipe.X - 6, Y: g.offsetY + bottomY, W: g.pipeWidth + 12, H: 22}
		DrawRoundedRect(g.renderer, capRect2.X, capRect2.Y, capRect2.W, capRect2.H, 6)
	}

	if g.font != nil {
		text := strconv.Itoa(int(g.score))
		surface, err := g.font.RenderUTF8Solid(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err == nil {
			defer surface.Free()
			texture, err := g.renderer.CreateTextureFromSurface(surface)
			if err == nil {
				defer texture.Destroy()
				_, _, w, h, _ := texture.Query()
				x := g.offsetX + g.width/2 - int32(w)/2
				y := g.offsetY + 50 - int32(h)/2
				g.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: w, H: h})
			}
		}
	}

	if g.gameOver {
		overlay := sdl.Rect{X: g.offsetX, Y: g.offsetY, W: g.width, H: g.height}
		g.renderer.SetDrawColor(0, 0, 0, 150)
		g.renderer.FillRect(&overlay)
		if g.font != nil {
			surface, err := g.font.RenderUTF8Solid("GAME OVER", sdl.Color{R: 255, G: 255, B: 255, A: 255})
			if err == nil {
				defer surface.Free()
				texture, err := g.renderer.CreateTextureFromSurface(surface)
				if err == nil {
					defer texture.Destroy()
					_, _, w, h, _ := texture.Query()
					x := g.offsetX + g.width/2 - int32(w)/2
					y := g.offsetY + g.height/2 - int32(h)/2 - 30
					g.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: w, H: h})
				}
			}
		}
		if g.fontSmall != nil {
			surface, err := g.fontSmall.RenderUTF8Solid("Press R to restart", sdl.Color{R: 255, G: 255, B: 255, A: 255})
			if err == nil {
				defer surface.Free()
				texture, err := g.renderer.CreateTextureFromSurface(surface)
				if err == nil {
					defer texture.Destroy()
					_, _, w, h, _ := texture.Query()
					x := g.offsetX + g.width/2 - int32(w)/2
					y := g.offsetY + g.height/2 + 20
					g.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: w, H: h})
				}
			}
		}
	}

	g.renderer.Present()
}

func (g *FlappyGame) Run() {
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