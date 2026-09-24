package games

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Brick struct {
	Rect  [4]int32
	Color color.RGBA
	Alive bool
}

type BreakoutGame struct {
	width       int32
	height      int32
	offsetX     int32
	offsetY     int32
	paddle      [4]float32
	paddleSpd   float32
	ball        [4]float32
	ballDx      float32
	ballDy      float32
	ballSpeed   float32
	bricks      []Brick
	score       int32
	lives       int32
	gameOver    bool
	waiting     bool
	brickCache  *TileCache
	paddleCache *TileCache
	ballCache   *TileCache
}

func NewBreakoutGame() *BreakoutGame {
	g := &BreakoutGame{
		width:     600,
		height:    600,
		paddleSpd: 7,
		ballSpeed: 5,
	}
	g.offsetX = (ScreenWidth - g.width) / 2
	g.offsetY = (ScreenHeight - g.height) / 2
	bw := (g.width - 20) / 8 - 4
	bh := int32(22)
	g.brickCache = NewTileCache(float32(bw), float32(bh), 6)
	g.paddleCache = NewTileCache(120, 16, 10)
	g.ballCache = NewTileCache(20, 20, 10)
	g.Reset()
	return g
}

func (g *BreakoutGame) Reset() {
	g.paddle = [4]float32{float32(g.width)/2 - 60, float32(g.height) - 40, 120, 16}
	g.ball = [4]float32{float32(g.width)/2 - 10, float32(g.height) - 70, 20, 20}
	g.ballDx = 5
	g.ballDy = -6
	g.bricks = nil
	g.score = 0
	g.lives = 3
	g.gameOver = false
	g.waiting = true

	rows, cols := 5, 8
	brickWidth := (g.width - 20) / int32(cols) - 4
	brickHeight := int32(22)
	colors := []color.RGBA{
		{230, 50, 50, 255},
		{230, 150, 50, 255},
		{230, 230, 50, 255},
		{50, 230, 50, 255},
		{50, 150, 230, 255},
	}
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			x := 10 + int32(col)*(brickWidth+4)
			y := 40 + int32(row)*(brickHeight+4)
			g.bricks = append(g.bricks, Brick{
				Rect:  [4]int32{x, y, brickWidth, brickHeight},
				Color: colors[row%len(colors)],
				Alive: true,
			})
		}
	}
}

func rectCollide(a, b [4]float32) bool {
	return a[0] < b[0]+b[2] && a[0]+a[2] > b[0] &&
		a[1] < b[1]+b[3] && a[1]+a[3] > b[1]
}

func absF32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func (g *BreakoutGame) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Reset()
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && g.waiting && !g.gameOver {
		g.waiting = false
	}

	if !g.gameOver && !g.waiting {
		if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.paddle[0] > 0 {
			g.paddle[0] -= g.paddleSpd
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) && g.paddle[0]+g.paddle[2] < float32(g.width) {
			g.paddle[0] += g.paddleSpd
		}
	}

	if g.gameOver || g.waiting {
		return nil
	}

	g.ball[0] += g.ballDx
	g.ball[1] += g.ballDy

	if g.ball[0] <= 0 || g.ball[0]+g.ball[2] >= float32(g.width) {
		g.ballDx = -g.ballDx
	}
	if g.ball[1] <= 0 {
		g.ballDy = -g.ballDy
	}

	if g.ball[1]+g.ball[3] >= float32(g.height) {
		g.lives--
		if g.lives <= 0 {
			g.gameOver = true
		} else {
			g.waiting = true
			g.ball[0] = float32(g.width)/2 - 10
			g.ball[1] = float32(g.height) - 70
			if rand.Intn(2) == 0 {
				g.ballDx = g.ballSpeed
			} else {
				g.ballDx = -g.ballSpeed
			}
			g.ballDy = -g.ballSpeed
		}
		return nil
	}

	if rectCollide(g.ball, g.paddle) {
		g.ballDy = -absF32(g.ballDy)
		hit := (g.ball[0] + g.ball[2]/2 - g.paddle[0] - g.paddle[2]/2) / (g.paddle[2] / 2)
		g.ballDx = hit * g.ballSpeed * 0.9
		if absF32(g.ballDx) < 2 {
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
		br := g.bricks[i]
		brF := [4]float32{float32(br.Rect[0]), float32(br.Rect[1]), float32(br.Rect[2]), float32(br.Rect[3])}
		if rectCollide(g.ball, brF) {
			g.bricks[i].Alive = false
			g.score += 10

			overlapTop := g.ball[1] + g.ball[3] - brF[1]
			overlapBottom := brF[1] + brF[3] - g.ball[1]
			overlapLeft := g.ball[0] + g.ball[2] - brF[0]
			overlapRight := brF[0] + brF[2] - g.ball[0]

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
	return nil
}

func (g *BreakoutGame) Draw(screen *ebiten.Image) {
	screen.Fill(ColorBg)

	DrawRoundedRect(screen, float32(g.offsetX), float32(g.offsetY),
		float32(g.width), float32(g.height), 12, color.RGBA{20, 20, 30, 255})

	for _, br := range g.bricks {
		if br.Alive {
			img := g.brickCache.GetWithBorder(br.Color, color.RGBA{255, 255, 255, 100})
			DrawImageAt(screen, img,
				float32(g.offsetX+br.Rect[0]),
				float32(g.offsetY+br.Rect[1]))
		}
	}

	DrawImageAt(screen, g.paddleCache.Get(color.RGBA{255, 255, 255, 255}),
		float32(g.offsetX)+g.paddle[0],
		float32(g.offsetY)+g.paddle[1])

	DrawImageAt(screen, g.ballCache.Get(color.RGBA{255, 255, 120, 255}),
		float32(g.offsetX)+g.ball[0],
		float32(g.offsetY)+g.ball[1])

	if FontSmall != nil {
		text.Draw(screen, "Score: "+fmt.Sprint(g.score), FontSmall,
			int(g.offsetX)+10, int(g.offsetY)+40, ColorWhite)
		lives := "Lives: " + fmt.Sprint(g.lives)
		b := text.BoundString(FontSmall, lives)
		text.Draw(screen, lives, FontSmall,
			int(g.offsetX+g.width)-b.Dx()-10, int(g.offsetY)+40, ColorWhite)
	}

	if g.waiting && !g.gameOver {
		if FontBig != nil {
			msg := "Press SPACE"
			b := text.BoundString(FontBig, msg)
			text.Draw(screen, msg, FontBig,
				int(g.offsetX)+int(g.width)/2-b.Dx()/2,
				int(g.offsetY)+int(g.height)/2+30, ColorWhite)
		}
		if FontSmall != nil {
			msg := "to start"
			b := text.BoundString(FontSmall, msg)
			text.Draw(screen, msg, FontSmall,
				int(g.offsetX)+int(g.width)/2-b.Dx()/2,
				int(g.offsetY)+int(g.height)/2+70, ColorWhite)
		}
	} else if g.gameOver {
		DrawOverlay(screen, int(g.offsetX), int(g.offsetY), int(g.width), int(g.height))
		msg := "YOU WIN!"
		clr := color.RGBA{0, 255, 0, 255}
		if g.lives <= 0 {
			msg = "GAME OVER"
			clr = ColorWhite
		}
		if FontBig != nil {
			b := text.BoundString(FontBig, msg)
			text.Draw(screen, msg, FontBig,
				int(g.offsetX)+int(g.width)/2-b.Dx()/2,
				int(g.offsetY)+int(g.height)/2-20, clr)
		}
		if FontSmall != nil {
			msg2 := "Press R to restart"
			b2 := text.BoundString(FontSmall, msg2)
			text.Draw(screen, msg2, FontSmall,
				int(g.offsetX)+int(g.width)/2-b2.Dx()/2,
				int(g.offsetY)+int(g.height)/2+30, ColorWhite)
		}
	}
}

func (g *BreakoutGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}