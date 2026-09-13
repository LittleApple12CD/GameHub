package games

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type SnakeGame struct {
	renderer   *sdl.Renderer
	width      int32
	height     int32
	offsetX    int32
	offsetY    int32
	cellSize   int32
	gridWidth  int32
	gridHeight int32
	font       *ttf.Font
	snake      [][2]int32
	dir        [2]int32
	nextDir    [2]int32
	food       [2]int32
	score      int32
	gameOver   bool
	moveTimer  uint32
	moveDelay  uint32
	clock      uint32
	running    bool
}

func NewSnakeGame(renderer *sdl.Renderer) *SnakeGame {
	g := &SnakeGame{
		renderer:   renderer,
		width:      600,
		height:     600,
		cellSize:   20,
		moveDelay:  150,
	}
	winW, winH, _ := renderer.GetOutputSize()
	g.offsetX = (int32(winW) - g.width) / 2
	g.offsetY = (int32(winH) - g.height) / 2
	g.gridWidth = g.width / g.cellSize
	g.gridHeight = g.height / g.cellSize
	g.font, _ = ttf.OpenFont("assets/fonts/arial.ttf", 36)
	if g.font == nil {
		g.font, _ = ttf.OpenFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 36)
	}
	g.Reset()
	return g
}

func (g *SnakeGame) Reset() {
	g.snake = [][2]int32{{g.gridWidth / 2, g.gridHeight / 2}}
	g.dir = [2]int32{1, 0}
	g.nextDir = [2]int32{1, 0}
	g.food = g.spawnFood()
	g.score = 0
	g.gameOver = false
	g.moveTimer = 0
	g.running = true
}

func (g *SnakeGame) spawnFood() [2]int32 {
	rand.Seed(time.Now().UnixNano())
	for {
		pos := [2]int32{int32(rand.Intn(int(g.gridWidth))), int32(rand.Intn(int(g.gridHeight)))}
		ok := true
		for _, seg := range g.snake {
			if seg[0] == pos[0] && seg[1] == pos[1] {
				ok = false
				break
			}
		}
		if ok {
			return pos
		}
	}
}

func (g *SnakeGame) HandleEvents() bool {
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
				case sdl.K_UP:
					if g.dir != [2]int32{0, 1} {
						g.nextDir = [2]int32{0, -1}
					}
				case sdl.K_DOWN:
					if g.dir != [2]int32{0, -1} {
						g.nextDir = [2]int32{0, 1}
					}
				case sdl.K_LEFT:
					if g.dir != [2]int32{1, 0} {
						g.nextDir = [2]int32{-1, 0}
					}
				case sdl.K_RIGHT:
					if g.dir != [2]int32{-1, 0} {
						g.nextDir = [2]int32{1, 0}
					}
				}
			}
		}
	}
	return true
}

func (g *SnakeGame) Update() {
    if g.gameOver || !g.running {
        return
    }

    g.moveTimer += 16
    
    if g.moveTimer >= g.moveDelay {
        g.moveTimer = 0
        g.dir = g.nextDir
        head := g.snake[0]
        newHead := [2]int32{head[0] + g.dir[0], head[1] + g.dir[1]}

        if newHead[0] < 0 || newHead[0] >= g.gridWidth ||
            newHead[1] < 0 || newHead[1] >= g.gridHeight {
            g.gameOver = true
            return
        }

        for _, seg := range g.snake {
            if seg[0] == newHead[0] && seg[1] == newHead[1] {
                g.gameOver = true
                return
            }
        }

        g.snake = append([][2]int32{newHead}, g.snake...)

        if newHead[0] == g.food[0] && newHead[1] == g.food[1] {
            g.score += 10
            g.food = g.spawnFood()
        } else {
            g.snake = g.snake[:len(g.snake)-1]
        }
    }
}

func (g *SnakeGame) Draw() {
	winW, winH, _ := g.renderer.GetOutputSize()
	g.offsetX = (int32(winW) - g.width) / 2
	g.offsetY = (int32(winH) - g.height) / 2

	g.renderer.SetDrawColor(25, 25, 35, 255)
	g.renderer.Clear()

	g.renderer.SetDrawColor(20, 20, 30, 255)
	DrawRoundedRect(g.renderer, g.offsetX, g.offsetY, g.width, g.height, 12)

	for i, seg := range g.snake {
		var color sdl.Color
		if i == 0 {
			color = sdl.Color{R: 60, G: 220, B: 60, A: 255}
		} else {
			color = sdl.Color{R: 40, G: 180, B: 40, A: 255}
		}
		rect := sdl.Rect{
			X: g.offsetX + seg[0]*g.cellSize + 1,
			Y: g.offsetY + seg[1]*g.cellSize + 1,
			W: g.cellSize - 2,
			H: g.cellSize - 2,
		}
		g.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
		DrawRoundedRect(g.renderer, rect.X, rect.Y, rect.W, rect.H, 4)
	}

	foodRect := sdl.Rect{
		X: g.offsetX + g.food[0]*g.cellSize + 1,
		Y: g.offsetY + g.food[1]*g.cellSize + 1,
		W: g.cellSize - 2,
		H: g.cellSize - 2,
	}
	g.renderer.SetDrawColor(255, 60, 60, 255)
	DrawRoundedRect(g.renderer, foodRect.X, foodRect.Y, foodRect.W, foodRect.H, 6)

	if g.font != nil {
		text := "Score: " + strconv.Itoa(int(g.score))
		surface, err := g.font.RenderUTF8Solid(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
		if err == nil {
			defer surface.Free()
			texture, err := g.renderer.CreateTextureFromSurface(surface)
			if err == nil {
				defer texture.Destroy()
				_, _, w, h, _ := texture.Query()
				g.renderer.Copy(texture, nil, &sdl.Rect{X: g.offsetX + 10, Y: g.offsetY + 10, W: w, H: h})
			}
		}
	}

	if g.gameOver {
		overlay := sdl.Rect{X: g.offsetX, Y: g.offsetY, W: g.width, H: g.height}
		g.renderer.SetDrawColor(0, 0, 0, 180)
		g.renderer.FillRect(&overlay)
		if g.font != nil {
			surface, err := g.font.RenderUTF8Solid("GAME OVER - Press R to Reset", sdl.Color{R: 255, G: 255, B: 255, A: 255})
			if err == nil {
				defer surface.Free()
				texture, err := g.renderer.CreateTextureFromSurface(surface)
				if err == nil {
					defer texture.Destroy()
					_, _, w, h, _ := texture.Query()
					x := g.offsetX + g.width/2 - int32(w)/2
					y := g.offsetY + g.height/2 - int32(h)/2
					g.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: w, H: h})
				}
			}
		}
	}

	g.renderer.Present()
}

func (g *SnakeGame) Run() {
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