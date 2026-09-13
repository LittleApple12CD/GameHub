package games

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Tank struct {
	X         int32
	Y         int32
	Dir       [2]int32
	Alive     bool
	Timer     int32
	MoveTimer uint32
}

type Bullet struct {
	X        float32
	Y        float32
	Dir      [2]int32
	IsPlayer bool
	Alive    bool
}

type TankGame struct {
	renderer        *sdl.Renderer
	width           int32
	height          int32
	offsetX         int32
	offsetY         int32
	cellSize        int32
	gridSize        int32
	bulletSpeed     float32
	font            *ttf.Font
	fontBig         *ttf.Font
	walls           map[[2]int32]bool
	player          Tank
	enemies         []Tank
	bullets         []Bullet
	score           int32
	gameOver        bool
	enemySpawnTimer int32
	enemySpawnDelay int32
	moveDelay       uint32
	enemyMoveDelay  uint32
	clock           uint32
	running         bool
}

func NewTankGame(renderer *sdl.Renderer) *TankGame {
	g := &TankGame{
		renderer:        renderer,
		width:           600,
		height:          600,
		cellSize:        30,
		bulletSpeed:     1.0,
		enemySpawnDelay: 120,
		moveDelay:       200,
		enemyMoveDelay:  250,
	}
	g.gridSize = g.width / g.cellSize
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

func (g *TankGame) Reset() {
	g.walls = g.generateWalls()
	g.player = Tank{X: 1, Y: 1, Dir: [2]int32{1, 0}, Alive: true, MoveTimer: 0}
	g.enemies = []Tank{{X: g.gridSize - 2, Y: g.gridSize - 2, Dir: [2]int32{-1, 0}, Alive: true, Timer: 0, MoveTimer: 0}}
	g.bullets = []Bullet{}
	g.score = 0
	g.gameOver = false
	g.enemySpawnTimer = 0
	g.clock = uint32(sdl.GetTicks())
	g.running = true
}

func (g *TankGame) generateWalls() map[[2]int32]bool {
	walls := make(map[[2]int32]bool)
	for i := int32(0); i < g.gridSize; i++ {
		walls[[2]int32{i, 0}] = true
		walls[[2]int32{i, g.gridSize - 1}] = true
		walls[[2]int32{0, i}] = true
		walls[[2]int32{g.gridSize - 1, i}] = true
	}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 12; i++ {
		x := int32(rand.Intn(int(g.gridSize-4))) + 2
		y := int32(rand.Intn(int(g.gridSize-4))) + 2
		if !(x == 1 && y == 1) && !(x == g.gridSize-2 && y == g.gridSize-2) {
			walls[[2]int32{x, y}] = true
		}
	}
	return walls
}

func (g *TankGame) canMove(x, y int32, dir [2]int32, isPlayer bool) bool {
	if g.walls[[2]int32{x, y}] {
		return false
	}
	if isPlayer {
		for _, enemy := range g.enemies {
			if enemy.Alive && enemy.X == x && enemy.Y == y {
				return false
			}
		}
	} else {
		if g.player.Alive && g.player.X == x && g.player.Y == y {
			return false
		}
		for _, enemy := range g.enemies {
			if enemy.Alive && enemy.X == x && enemy.Y == y {
				return false
			}
		}
	}
	return x >= 0 && x < g.gridSize && y >= 0 && y < g.gridSize
}

func (g *TankGame) shootBullet(x, y int32, dir [2]int32, isPlayer bool) {
	g.bullets = append(g.bullets, Bullet{
		X:        float32(x + dir[0]),
		Y:        float32(y + dir[1]),
		Dir:      dir,
		IsPlayer: isPlayer,
		Alive:    true,
	})
}

func (g *TankGame) HandleEvents() bool {
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
					if !g.gameOver && g.player.Alive {
						g.shootBullet(g.player.X, g.player.Y, g.player.Dir, true)
					}
				}
			}
		}
	}

	if !g.gameOver && g.player.Alive {
		keys := sdl.GetKeyboardState()
		dx, dy := int32(0), int32(0)
		if keys[sdl.SCANCODE_W] != 0 {
			dy = -1
		} else if keys[sdl.SCANCODE_S] != 0 {
			dy = 1		
		} else if keys[sdl.SCANCODE_A] != 0 {
			dx = -1
		} else if keys[sdl.SCANCODE_D] != 0 {
			dx = 1
		}

		if dx != 0 || dy != 0 {
			g.player.Dir = [2]int32{dx, dy}
			now := uint32(sdl.GetTicks())
			dt := now - g.clock
			g.player.MoveTimer += dt
			if g.player.MoveTimer >= g.moveDelay {
				g.player.MoveTimer = 0
				newX := g.player.X + dx
				newY := g.player.Y + dy
				if g.canMove(newX, newY, g.player.Dir, true) {
					g.player.X = newX
					g.player.Y = newY
				}
			}
		} else {
			g.player.MoveTimer = 0
		}
	}
	return true
}

func (g *TankGame) Update() {
	if g.gameOver || !g.running {
		return
	}

	now := uint32(sdl.GetTicks())
	dt := now - g.clock
	g.clock = now

	for i := range g.bullets {
		if !g.bullets[i].Alive {
			continue
		}
		b := &g.bullets[i]
		b.X += float32(b.Dir[0]) * g.bulletSpeed
		b.Y += float32(b.Dir[1]) * g.bulletSpeed

		bx := int32(b.X)
		by := int32(b.Y)

		if bx < 0 || bx >= g.gridSize || by < 0 || by >= g.gridSize {
			b.Alive = false
			continue
		}

		if g.walls[[2]int32{bx, by}] {
			b.Alive = false
			continue
		}

		if b.IsPlayer {
			for j := range g.enemies {
				if g.enemies[j].Alive && g.enemies[j].X == bx && g.enemies[j].Y == by {
					g.enemies[j].Alive = false
					b.Alive = false
					g.score += 10
					break
				}
			}
		} else {
			if g.player.Alive && g.player.X == bx && g.player.Y == by {
				g.player.Alive = false
				b.Alive = false
				g.gameOver = true
			}
		}
	}

	newBullets := make([]Bullet, 0)
	for _, b := range g.bullets {
		if b.Alive {
			newBullets = append(newBullets, b)
		}
	}
	g.bullets = newBullets

	for i := range g.enemies {
		if !g.enemies[i].Alive {
			continue
		}
		e := &g.enemies[i]
		e.Timer++
		if e.Timer >= 20 {
			e.Timer = 0
			if rand.Float32() < 0.25 {
				dirs := [][2]int32{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
				e.Dir = dirs[rand.Intn(4)]
			}
			e.MoveTimer += dt
			if e.MoveTimer >= g.enemyMoveDelay {
				e.MoveTimer = 0
				newX := e.X + e.Dir[0]
				newY := e.Y + e.Dir[1]
				if g.canMove(newX, newY, e.Dir, false) {
					e.X = newX
					e.Y = newY
				}
			}
		}
		if rand.Float32() < 0.015 {
			g.shootBullet(e.X, e.Y, e.Dir, false)
		}
	}

	g.enemySpawnTimer++
	if g.enemySpawnTimer >= g.enemySpawnDelay {
		g.enemySpawnTimer = 0
		aliveCount := 0
		for _, e := range g.enemies {
			if e.Alive {
				aliveCount++
			}
		}
		if aliveCount < 3 {
			g.enemies = append(g.enemies, Tank{
				X:         g.gridSize - 2,
				Y:         g.gridSize - 2,
				Dir:       [2]int32{-1, 0},
				Alive:     true,
				Timer:     0,
				MoveTimer: 0,
			})
		}
	}
}

func (g *TankGame) drawTank(x, y int32, color, gunColor sdl.Color, dir [2]int32) {
	cx := g.offsetX + x*g.cellSize + 15
	cy := g.offsetY + y*g.cellSize + 15

	rect := sdl.Rect{
		X: g.offsetX + x*g.cellSize + 2,
		Y: g.offsetY + y*g.cellSize + 2,
		W: 26, H: 26,
	}
	g.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	DrawRoundedRect(g.renderer, rect.X, rect.Y, rect.W, rect.H, 4)

	var gunRect sdl.Rect
	if dir[0] == 1 {
		gunRect = sdl.Rect{X: cx + 4, Y: cy - 4, W: 14, H: 8}
	} else if dir[0] == -1 {
		gunRect = sdl.Rect{X: cx - 18, Y: cy - 4, W: 14, H: 8}
	} else if dir[1] == -1 {
		gunRect = sdl.Rect{X: cx - 4, Y: cy - 18, W: 8, H: 14}
	} else {
		gunRect = sdl.Rect{X: cx - 4, Y: cy + 4, W: 8, H: 14}
	}
	g.renderer.SetDrawColor(gunColor.R, gunColor.G, gunColor.B, gunColor.A)
	g.renderer.FillRect(&gunRect)

	g.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
	for dy := int32(-6); dy <= 6; dy++ {
		for dx := int32(-6); dx <= 6; dx++ {
			if dx*dx+dy*dy <= 36 {
				g.renderer.DrawPoint(cx+dx, cy+dy)
			}
		}
	}
}

func (g *TankGame) Draw() {
	winW, winH, _ := g.renderer.GetOutputSize()
	g.offsetX = (int32(winW) - g.width) / 2
	g.offsetY = (int32(winH) - g.height) / 2

	g.renderer.SetDrawColor(25, 25, 35, 255)
	g.renderer.Clear()

	g.renderer.SetDrawColor(30, 30, 42, 255)
	DrawRoundedRect(g.renderer, g.offsetX, g.offsetY, g.width, g.height, 12)

	for pos := range g.walls {
		rect := sdl.Rect{
			X: g.offsetX + pos[0]*g.cellSize + 1,
			Y: g.offsetY + pos[1]*g.cellSize + 1,
			W: g.cellSize - 2,
			H: g.cellSize - 2,
		}
		g.renderer.SetDrawColor(100, 100, 130, 255)
		DrawRoundedRect(g.renderer, rect.X, rect.Y, rect.W, rect.H, 4)
	}

	if g.player.Alive {
		g.drawTank(g.player.X, g.player.Y,
			sdl.Color{R: 50, G: 230, B: 50, A: 255},
			sdl.Color{R: 30, G: 200, B: 30, A: 255},
			g.player.Dir)
	}

	for _, enemy := range g.enemies {
		if enemy.Alive {
			g.drawTank(enemy.X, enemy.Y,
				sdl.Color{R: 230, G: 50, B: 50, A: 255},
				sdl.Color{R: 200, G: 30, B: 30, A: 255},
				enemy.Dir)
		}
	}

	for _, bullet := range g.bullets {
		var color sdl.Color
		if bullet.IsPlayer {
			color = sdl.Color{R: 255, G: 255, B: 80, A: 255}
		} else {
			color = sdl.Color{R: 255, G: 150, B: 50, A: 255}
		}
		rect := sdl.Rect{
			X: g.offsetX + int32(bullet.X)*g.cellSize + 10,
			Y: g.offsetY + int32(bullet.Y)*g.cellSize + 10,
			W: 10, H: 10,
		}
		g.renderer.SetDrawColor(color.R, color.G, color.B, color.A)
		DrawRoundedRect(g.renderer, rect.X, rect.Y, rect.W, rect.H, 5)
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
				g.renderer.Copy(texture, nil, &sdl.Rect{X: g.offsetX + 10, Y: g.offsetY + 10, W: w, H: h})
			}
		}
	}

	if g.gameOver {
		overlay := sdl.Rect{X: g.offsetX, Y: g.offsetY, W: g.width, H: g.height}
		g.renderer.SetDrawColor(0, 0, 0, 150)
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

func (g *TankGame) Run() {
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