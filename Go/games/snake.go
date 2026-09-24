package games

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type SnakeGame struct {
	width      int32
	height     int32
	offsetX    int32
	offsetY    int32
	cellSize   int32
	gridWidth  int32
	gridHeight int32
	snake      [][2]int32
	dir        [2]int32
	nextDir    [2]int32
	food       [2]int32
	score      int32
	gameOver   bool
	moveTimer  float64
	moveDelay  float64
	headCache  *TileCache
	bodyCache  *TileCache
	foodCache  *TileCache
}

func NewSnakeGame() *SnakeGame {
	g := &SnakeGame{
		width:     600,
		height:    600,
		cellSize:  20,
		moveDelay: 0.15,
	}
	g.offsetX = (ScreenWidth - g.width) / 2
	g.offsetY = (ScreenHeight - g.height) / 2
	g.gridWidth = g.width / g.cellSize
	g.gridHeight = g.height / g.cellSize
	g.headCache = NewTileCache(float32(g.cellSize-2), float32(g.cellSize-2), 4)
	g.bodyCache = NewTileCache(float32(g.cellSize-2), float32(g.cellSize-2), 4)
	g.foodCache = NewTileCache(float32(g.cellSize-2), float32(g.cellSize-2), 6)
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
}

func (g *SnakeGame) spawnFood() [2]int32 {
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

func (g *SnakeGame) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Reset()
	}

	if g.gameOver {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) && g.dir != [2]int32{0, 1} {
		g.nextDir = [2]int32{0, -1}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) && g.dir != [2]int32{0, -1} {
		g.nextDir = [2]int32{0, 1}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && g.dir != [2]int32{1, 0} {
		g.nextDir = [2]int32{-1, 0}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && g.dir != [2]int32{-1, 0} {
		g.nextDir = [2]int32{1, 0}
	}

	g.moveTimer += 1.0 / 60.0
	if g.moveTimer >= g.moveDelay {
		g.moveTimer = 0
		g.dir = g.nextDir
		head := g.snake[0]
		newHead := [2]int32{head[0] + g.dir[0], head[1] + g.dir[1]}

		if newHead[0] < 0 || newHead[0] >= g.gridWidth ||
			newHead[1] < 0 || newHead[1] >= g.gridHeight {
			g.gameOver = true
			return nil
		}
		for _, seg := range g.snake {
			if seg[0] == newHead[0] && seg[1] == newHead[1] {
				g.gameOver = true
				return nil
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
	return nil
}

func (g *SnakeGame) Draw(screen *ebiten.Image) {
	screen.Fill(ColorBg)

	DrawRoundedRect(screen, float32(g.offsetX), float32(g.offsetY),
		float32(g.width), float32(g.height), 12, color.RGBA{20, 20, 30, 255})

	for i, seg := range g.snake {
		var img *ebiten.Image
		if i == 0 {
			img = g.headCache.Get(color.RGBA{60, 220, 60, 255})
		} else {
			img = g.bodyCache.Get(color.RGBA{40, 180, 40, 255})
		}
		DrawImageAt(screen, img,
			float32(g.offsetX+seg[0]*g.cellSize+1),
			float32(g.offsetY+seg[1]*g.cellSize+1))
	}

	DrawImageAt(screen, g.foodCache.Get(color.RGBA{255, 60, 60, 255}),
		float32(g.offsetX+g.food[0]*g.cellSize+1),
		float32(g.offsetY+g.food[1]*g.cellSize+1))

	if FontSmall != nil {
		text.Draw(screen, "Score: "+fmt.Sprint(g.score), FontSmall,
			int(g.offsetX)+10, int(g.offsetY)+40, ColorWhite)
	}

	if g.gameOver {
		DrawOverlay(screen, int(g.offsetX), int(g.offsetY), int(g.width), int(g.height))
		msg := "GAME OVER - Press R to Reset"
		b := text.BoundString(FontBig, msg)
		text.Draw(screen, msg, FontBig,
			int(g.offsetX)+int(g.width)/2-b.Dx()/2,
			int(g.offsetY)+int(g.height)/2, ColorWhite)
	}
}

func (g *SnakeGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func DrawOverlay(screen *ebiten.Image, x, y, w, h int) {
	DrawRoundedRect(screen, float32(x), float32(y), float32(w), float32(h), 0,
		color.RGBA{0, 0, 0, 180})
}