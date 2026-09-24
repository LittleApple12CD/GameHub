package games

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Tank struct {
	X     int32
	Y     int32
	Dir   [2]int32
	Alive bool
	Timer int32
}

type Bullet struct {
	X        float32
	Y        float32
	Dir      [2]int32
	IsPlayer bool
	Alive    bool
}

type TankGame struct {
	width           int32
	height          int32
	offsetX         int32
	offsetY         int32
	cellSize        int32
	gridSize        int32
	bulletSpeed     float32
	walls           map[[2]int32]bool
	player          Tank
	enemies         []Tank
	bullets         []Bullet
	score           int32
	gameOver        bool
	enemySpawnTimer int32
	enemySpawnDelay int32
	moveTimer       float64
	moveDelay       float64
	enemyMoveTimer  float64
	enemyMoveDelay  float64
	wallCache       *TileCache
}

func NewTankGame() *TankGame {
	g := &TankGame{
		width:           600,
		height:          600,
		cellSize:        30,
		bulletSpeed:     0.05,
		enemySpawnDelay: 120,
		moveDelay:       0.2,
		enemyMoveDelay:  0.25,
	}
	g.gridSize = g.width / g.cellSize
	g.offsetX = (ScreenWidth - g.width) / 2
	g.offsetY = (ScreenHeight - g.height) / 2
	g.wallCache = NewTileCache(float32(g.cellSize-2), float32(g.cellSize-2), 4)
	g.Reset()
	return g
}

func (g *TankGame) Reset() {
	g.walls = g.generateWalls()
	g.player = Tank{X: 1, Y: 1, Dir: [2]int32{1, 0}, Alive: true}
	g.enemies = []Tank{{X: g.gridSize - 2, Y: g.gridSize - 2, Dir: [2]int32{-1, 0}, Alive: true}}
	g.bullets = []Bullet{}
	g.score = 0
	g.gameOver = false
	g.enemySpawnTimer = 0
	g.moveTimer = 0
	g.enemyMoveTimer = 0
}

func (g *TankGame) generateWalls() map[[2]int32]bool {
	walls := make(map[[2]int32]bool)
	for i := int32(0); i < g.gridSize; i++ {
		walls[[2]int32{i, 0}] = true
		walls[[2]int32{i, g.gridSize - 1}] = true
		walls[[2]int32{0, i}] = true
		walls[[2]int32{g.gridSize - 1, i}] = true
	}
	for i := 0; i < 12; i++ {
		x := int32(rand.Intn(int(g.gridSize-4))) + 2
		y := int32(rand.Intn(int(g.gridSize-4))) + 2
		if !(x == 1 && y == 1) && !(x == g.gridSize-2 && y == g.gridSize-2) {
			walls[[2]int32{x, y}] = true
		}
	}
	return walls
}

func (g *TankGame) canMove(x, y int32, isPlayer bool) bool {
	if g.walls[[2]int32{x, y}] {
		return false
	}
	if isPlayer {
		for _, e := range g.enemies {
			if e.Alive && e.X == x && e.Y == y {
				return false
			}
		}
	} else {
		if g.player.Alive && g.player.X == x && g.player.Y == y {
			return false
		}
		for _, e := range g.enemies {
			if e.Alive && e.X == x && e.Y == y {
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

func (g *TankGame) Update() error {
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

	dt := 1.0 / 60.0

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.player.Alive {
		g.shootBullet(g.player.X, g.player.Y, g.player.Dir, true)
	}

	if g.player.Alive {
		dx, dy := int32(0), int32(0)
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			dy = -1
		} else if ebiten.IsKeyPressed(ebiten.KeyS) {
			dy = 1
		} else if ebiten.IsKeyPressed(ebiten.KeyA) {
			dx = -1
		} else if ebiten.IsKeyPressed(ebiten.KeyD) {
			dx = 1
		}

		if dx != 0 || dy != 0 {
			g.player.Dir = [2]int32{dx, dy}
			g.moveTimer += dt
			if g.moveTimer >= g.moveDelay {
				g.moveTimer = 0
				nx := g.player.X + dx
				ny := g.player.Y + dy
				if g.canMove(nx, ny, true) {
					g.player.X = nx
					g.player.Y = ny
				}
			}
		} else {
			g.moveTimer = 0
		}
	}

	for i := range g.bullets {
		b := &g.bullets[i]
		if !b.Alive {
			continue
		}
		b.X += float32(b.Dir[0]) * g.bulletSpeed * 60
		b.Y += float32(b.Dir[1]) * g.bulletSpeed * 60
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

	newBullets := g.bullets[:0]
	for _, b := range g.bullets {
		if b.Alive {
			newBullets = append(newBullets, b)
		}
	}
	g.bullets = newBullets

	for i := range g.enemies {
		e := &g.enemies[i]
		if !e.Alive {
			continue
		}
		e.Timer++
		if e.Timer >= 20 {
			e.Timer = 0
			if rand.Float32() < 0.25 {
				dirs := [][2]int32{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
				e.Dir = dirs[rand.Intn(4)]
			}
			g.enemyMoveTimer += dt
			if g.enemyMoveTimer >= g.enemyMoveDelay {
				g.enemyMoveTimer = 0
				nx := e.X + e.Dir[0]
				ny := e.Y + e.Dir[1]
				if g.canMove(nx, ny, false) {
					e.X = nx
					e.Y = ny
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
		alive := 0
		for _, e := range g.enemies {
			if e.Alive {
				alive++
			}
		}
		if alive < 3 {
			g.enemies = append(g.enemies, Tank{
				X: g.gridSize - 2, Y: g.gridSize - 2,
				Dir: [2]int32{-1, 0}, Alive: true,
			})
		}
	}
	return nil
}

func (g *TankGame) drawTank(screen *ebiten.Image, x, y int32, body, gun color.Color, dir [2]int32) {
	cx := float32(g.offsetX + x*g.cellSize + 15)
	cy := float32(g.offsetY + y*g.cellSize + 15)

	rectX := float32(g.offsetX + x*g.cellSize + 2)
	rectY := float32(g.offsetY + y*g.cellSize + 2)
	DrawRoundedRect(screen, rectX, rectY, 26, 26, 4, body)

	var gx, gy, gw, gh float32
	if dir[0] == 1 {
		gx, gy, gw, gh = cx+4, cy-4, 14, 8
	} else if dir[0] == -1 {
		gx, gy, gw, gh = cx-18, cy-4, 14, 8
	} else if dir[1] == -1 {
		gx, gy, gw, gh = cx-4, cy-18, 8, 14
	} else {
		gx, gy, gw, gh = cx-4, cy+4, 8, 14
	}
	vector.DrawFilledRect(screen, gx, gy, gw, gh, gun, true)
	vector.DrawFilledCircle(screen, cx, cy, 6, body, true)
}

func (g *TankGame) Draw(screen *ebiten.Image) {
	screen.Fill(ColorBg)

	DrawRoundedRect(screen, float32(g.offsetX), float32(g.offsetY),
		float32(g.width), float32(g.height), 12, color.RGBA{30, 30, 42, 255})

	wallImg := g.wallCache.Get(color.RGBA{100, 100, 130, 255})
	for pos := range g.walls {
		DrawImageAt(screen, wallImg,
			float32(g.offsetX+pos[0]*g.cellSize+1),
			float32(g.offsetY+pos[1]*g.cellSize+1))
	}

	if g.player.Alive {
		g.drawTank(screen, g.player.X, g.player.Y,
			color.RGBA{50, 230, 50, 255}, color.RGBA{30, 200, 30, 255}, g.player.Dir)
	}
	for _, e := range g.enemies {
		if e.Alive {
			g.drawTank(screen, e.X, e.Y,
				color.RGBA{230, 50, 50, 255}, color.RGBA{200, 30, 30, 255}, e.Dir)
		}
	}

	for _, b := range g.bullets {
		clr := color.RGBA{255, 255, 80, 255}
		if !b.IsPlayer {
			clr = color.RGBA{255, 150, 50, 255}
		}
		DrawRoundedRect(screen,
			float32(g.offsetX+int32(b.X)*g.cellSize+10),
			float32(g.offsetY+int32(b.Y)*g.cellSize+10),
			10, 10, 5, clr)
	}

	if FontSmall != nil {
		text.Draw(screen, "Score: "+fmt.Sprint(g.score), FontSmall,
			int(g.offsetX)+10, int(g.offsetY)+40, ColorWhite)
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

func (g *TankGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}