package main

import (
	"fmt"
	"image"
	_ "image/png"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"

	"gamehub/games"
)

const (
	ScreenWidth  = 850
	ScreenHeight = 650
)

var (
	ColorBg          = color.RGBA{25, 25, 35, 255}
	ColorWhite       = color.RGBA{255, 255, 255, 255}
	ColorHover       = color.RGBA{80, 200, 255, 255}
	ColorHoverBg     = color.RGBA{60, 80, 120, 255}
	ColorButton      = color.RGBA{55, 55, 75, 255}
	ColorBorder      = color.RGBA{120, 120, 160, 255}
	ColorBorderHover = color.RGBA{150, 220, 255, 255}
	ColorQuit        = color.RGBA{200, 50, 50, 255}
	ColorQuitHov     = color.RGBA{255, 80, 80, 255}
)

type GameInfo struct {
	Name     string
	New      func() games.Game
	instance games.Game
}

var gameList = []GameInfo{
	{"SnakeGame", func() games.Game { return games.NewSnakeGame() }, nil},
	{"TetrisGame", func() games.Game { return games.NewTetrisGame() }, nil},
	{"Minesweeper", func() games.Game { return games.NewMinesweeperGame() }, nil},
	{"Flappy Bird", func() games.Game { return games.NewFlappyGame() }, nil},
	{"TankBattle", func() games.Game { return games.NewTankGame() }, nil},
	{"Breakout", func() games.Game { return games.NewBreakoutGame() }, nil},
	{"Tic Tac Toe", func() games.Game { return games.NewTicTacToeGame() }, nil},
}

func loadIcon() []image.Image {
	icon, _, err := ebitenutil.NewImageFromFile("assets/icons/icon.png")
	if err == nil {
		return []image.Image{icon}
	}

	f, err := os.Open("assets/icons/icon.png")
	if err != nil {
		log.Printf("Warning: cannot open icon: %v", err)
		return nil
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Printf("Warning: cannot decode icon: %v", err)
		return nil
	}
	return []image.Image{img}
}

func loadFonts() {
	games.InitFonts()
}

type GameHub struct {
	selected int
	quitSel  bool
	inGame   bool
	current  games.Game
}

func (h *GameHub) Update() error {	
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if h.inGame {
			h.inGame = false
			h.current = nil
			return nil
		}
		return ebiten.Termination
	}

	if h.inGame && h.current != nil {
		if err := h.current.Update(); err != nil {
			if err == ebiten.Termination {
				h.inGame = false
				h.current = nil
			} else {
				return err
			}
		}
		return nil
	}

	// Hub 主菜单输入
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if h.quitSel {
			h.quitSel = false
			h.selected = len(gameList) - 1
		} else {
			h.selected = (h.selected - 1 + len(gameList)) % len(gameList)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if h.selected == len(gameList)-1 {
			h.quitSel = true
		} else {
			h.selected = (h.selected + 1) % len(gameList)
			h.quitSel = false
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if h.quitSel {
			return ebiten.Termination
		}
		h.launch(h.selected)
	}

	// 鼠标
	mx, my := ebiten.CursorPosition()
	h.quitSel = false
	for i := range gameList {
		y := 130 + i*52
		if mx >= ScreenWidth/2-140 && mx <= ScreenWidth/2+140 &&
			my >= y && my <= y+46 {
			h.selected = i
			break
		}
	}
	quitY := 130 + len(gameList)*52 + 20
	if mx >= ScreenWidth/2-60 && mx <= ScreenWidth/2+60 &&
		my >= quitY && my <= quitY+40 {
		h.quitSel = true
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if mx >= ScreenWidth/2-60 && mx <= ScreenWidth/2+60 &&
			my >= quitY && my <= quitY+40 {
			return ebiten.Termination
		}
		for i := range gameList {
			y := 130 + i*52
			if mx >= ScreenWidth/2-140 && mx <= ScreenWidth/2+140 &&
				my >= y && my <= y+46 {
				h.launch(i)
				break
			}
		}
	}
	return nil
}

func (h *GameHub) launch(index int) {
	info := &gameList[index]
	info.instance = info.New()
	info.instance.Reset()
	h.current = info.instance
	h.inGame = true
}

func (h *GameHub) Draw(screen *ebiten.Image) {
	if h.inGame && h.current != nil {
		h.current.Draw(screen)
		return
	}

	screen.Fill(ColorBg)

	// 标题
	drawCenteredText(screen, "GAME HUB", games.FontBig, ColorWhite, ScreenWidth/2, 55)

    for i, info := range gameList {
        y := float32(130 + i*52)
        x := float32(ScreenWidth/2 - 140)
        bw, bh := float32(280), float32(46)

        if i == h.selected {
            glow := color.RGBA{ColorHover.R, ColorHover.G, ColorHover.B, 80}
            games.DrawRoundedRect(screen, x-4, y-4, bw+8, bh+8, 14, glow)

            games.DrawRoundedRect(screen, x, y, bw, bh, 12, ColorHoverBg)
            games.DrawRoundedRectBorder(screen, x, y, bw, bh, 12, 3, ColorHover)
        } else {
            games.DrawRoundedRect(screen, x, y, bw, bh, 12, ColorButton)
            games.DrawRoundedRectBorder(screen, x, y, bw, bh, 12, 2, ColorBorder)
        }

        var textColor color.Color
        if i == h.selected {
            textColor = ColorWhite
        } else {
            textColor = color.RGBA{200, 200, 210, 255}
        }
        drawCenteredText(screen, info.Name, games.FontSmall, textColor, ScreenWidth/2, int(y)+23)
    }

	quitY := float32(130 + len(gameList)*52 + 20)
	qx := float32(ScreenWidth/2 - 60)
	qw, qh := float32(120), float32(40)
	qColor := ColorQuit
	if h.quitSel {
		qColor = ColorQuitHov
	}
	games.DrawRoundedRect(screen, qx, quitY, qw, qh, 10, qColor)
	border := color.RGBA{150, 40, 40, 255}
	if h.quitSel {
		border = color.RGBA{220, 60, 60, 255}
	}
	games.DrawRoundedRectBorder(screen, qx, quitY, qw, qh, 10, 2, border)
	drawCenteredText(screen, "EXIT", games.FontSmall, ColorWhite, ScreenWidth/2, int(quitY)+20)
}

func (h *GameHub) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func drawCenteredText(screen *ebiten.Image, s string, face font.Face, clr color.Color, cx, cy int) {
	if face == nil {
		return
	}
	bounds := text.BoundString(face, s)
	w := bounds.Dx()
	h := bounds.Dy()
	text.Draw(screen, s, face, cx-w/2, cy+h/2, clr)
}

func main() {
	loadFonts()

	if icon := loadIcon(); icon != nil {
		ebiten.SetWindowIcon(icon)
	}

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Game Hub")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	hub := &GameHub{}
	if err := ebiten.RunGame(hub); err != nil {
		fmt.Fprintf(os.Stderr, "Game exited: %v\n", err)
		os.Exit(1)
	}
}