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

type Pipe struct {
	X      int32
	GapY   int32
	Passed bool
}

type FlappyGame struct {
	width      int32
	height     int32
	offsetX    int32
	offsetY    int32
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
}

func NewFlappyGame() *FlappyGame {
	g := &FlappyGame{
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
	g.offsetX = (ScreenWidth - g.width) / 2
	g.offsetY = (ScreenHeight - g.height) / 2
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
}

func (g *FlappyGame) addPipe() {
	gapY := int32(rand.Intn(int(g.height-80-g.pipeGap))) + 40
	g.pipes = append(g.pipes, Pipe{X: g.width, GapY: gapY, Passed: false})
}

func (g *FlappyGame) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Reset()
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if !g.gameOver {
			g.birdVY = g.jumpStr
		}
	}

	if g.gameOver {
		return nil
	}

	g.birdVY += g.gravity
	g.birdY += g.birdVY

	if g.birdY-float32(g.birdRadius) < 0 {
		g.birdY = float32(g.birdRadius)
		g.birdVY = 0
	} else if g.birdY+float32(g.birdRadius) > float32(g.height) {
		g.gameOver = true
		return nil
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
				return nil
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
	return nil
}

func (g *FlappyGame) Draw(screen *ebiten.Image) {
	screen.Fill(ColorBg)

	DrawRoundedRect(screen, float32(g.offsetX), float32(g.offsetY),
		float32(g.width), float32(g.height), 12, color.RGBA{135, 206, 235, 255})

	bx := float32(g.offsetX + g.birdX)
	by := float32(g.offsetY) + g.birdY
	vector.DrawFilledCircle(screen, bx, by, float32(g.birdRadius), color.RGBA{255, 255, 50, 255}, true)
	vector.DrawFilledRect(screen, bx+8, by-5, 4, 4, color.Black, true)
	vector.DrawFilledRect(screen, bx+10, by-7, 2, 2, color.White, true)

	for _, pipe := range g.pipes {
		topRect := float32(g.offsetY)
		topH := float32(pipe.GapY)
		DrawRoundedRect(screen, float32(g.offsetX+pipe.X), topRect,
			float32(g.pipeWidth), topH, 0, color.RGBA{40, 200, 40, 255})
		DrawRoundedRect(screen, float32(g.offsetX+pipe.X-6),
			topRect+topH-22, float32(g.pipeWidth+12), 22, 6,
			color.RGBA{30, 160, 30, 255})

		bottomY := float32(g.offsetY + pipe.GapY + g.pipeGap)
		bottomH := float32(g.height) - float32(pipe.GapY+g.pipeGap)
		DrawRoundedRect(screen, float32(g.offsetX+pipe.X), bottomY,
			float32(g.pipeWidth), bottomH, 0, color.RGBA{40, 200, 40, 255})
		DrawRoundedRect(screen, float32(g.offsetX+pipe.X-6), bottomY,
			float32(g.pipeWidth+12), 22, 6, color.RGBA{30, 160, 30, 255})
	}

	if FontBig != nil {
		s := fmt.Sprint(g.score)
		b := text.BoundString(FontBig, s)
		text.Draw(screen, s, FontBig,
			int(g.offsetX)+int(g.width)/2-b.Dx()/2,
			int(g.offsetY)+50, ColorWhite)
	}

	if g.gameOver {
		DrawOverlay(screen, int(g.offsetX), int(g.offsetY), int(g.width), int(g.height))
		msg := "GAME OVER"
		b := text.BoundString(FontBig, msg)
		text.Draw(screen, msg, FontBig,
			int(g.offsetX)+int(g.width)/2-b.Dx()/2,
			int(g.offsetY)+int(g.height)/2-30, ColorWhite)
		msg2 := "Press R to restart"
		b2 := text.BoundString(FontSmall, msg2)
		text.Draw(screen, msg2, FontSmall,
			int(g.offsetX)+int(g.width)/2-b2.Dx()/2,
			int(g.offsetY)+int(g.height)/2+20, ColorWhite)
	}
}

func (g *FlappyGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}