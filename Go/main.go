package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"

	"gamehub/games"
)

const (
	ScreenWidth  = 850
	ScreenHeight = 650
)

var (
	ColorBg      = sdl.Color{R: 25, G: 25, B: 35, A: 255}
	ColorWhite   = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	ColorHover   = sdl.Color{R: 80, G: 200, B: 255, A: 255}
	ColorHoverBg = sdl.Color{R: 60, G: 80, B: 120, A: 255}
	ColorButton  = sdl.Color{R: 55, G: 55, B: 75, A: 255}
	ColorBorder  = sdl.Color{R: 120, G: 120, B: 160, A: 255}
	ColorBorderHover = sdl.Color{R: 150, G: 220, B: 255, A: 255}
	ColorQuit    = sdl.Color{R: 200, G: 50, B: 50, A: 255}
	ColorQuitHov = sdl.Color{R: 255, G: 80, B: 80, A: 255}
)

type GameInfo struct {
	Name string
	New  func(*sdl.Renderer) games.Game
	instance games.Game
}

var gameList = []GameInfo{
	{"Snake", func(r *sdl.Renderer) games.Game { return games.NewSnakeGame(r) }, nil},
	{"Tetris", func(r *sdl.Renderer) games.Game { return games.NewTetrisGame(r) }, nil},
	{"Minesweeper", func(r *sdl.Renderer) games.Game { return games.NewMinesweeperGame(r) }, nil},
	{"Flappy Bird", func(r *sdl.Renderer) games.Game { return games.NewFlappyGame(r) }, nil},
	{"Tank", func(r *sdl.Renderer) games.Game { return games.NewTankGame(r) }, nil},
	{"Breakout", func(r *sdl.Renderer) games.Game { return games.NewBreakoutGame(r) }, nil},
	{"Tic Tac Toe", func(r *sdl.Renderer) games.Game { return games.NewTicTacToeGame(r) }, nil},
}

func main() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %v\n", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize TTF: %v\n", err)
		os.Exit(1)
	}
	defer ttf.Quit()

	window, err := sdl.CreateWindow("Game Hub", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		ScreenWidth, ScreenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %v\n", err)
		os.Exit(1)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %v\n", err)
		os.Exit(1)
	}
	defer renderer.Destroy()

	font, err := ttf.OpenFont("assets/fonts/arial.ttf", 52)
	if err != nil {
		fmt.Printf("Warning: Cannot open font: %v\n", err)
		font = nil
	}
	if font != nil {
		defer font.Close()
	}

	fontSmall, err := ttf.OpenFont("assets/fonts/arial.ttf", 36)
	if err != nil {
		fontSmall = nil
	}
	if fontSmall != nil {
		defer fontSmall.Close()
	}

	hub := &GameHub{
		renderer:   renderer,
		window:     window,
		font:       font,
		fontSmall:  fontSmall,
		selected:   0,
		quitSel:    false,
		running:    true,
	}

	hub.Run()
}

type GameHub struct {
	renderer  *sdl.Renderer
	window    *sdl.Window
	font      *ttf.Font
	fontSmall *ttf.Font
	selected  int
	quitSel   bool
	running   bool
}

func (h *GameHub) Run() {
	for h.running {
		h.HandleEvents()
		h.Draw()
		sdl.Delay(16)
	}
}

func (h *GameHub) HandleEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			h.running = false
			return
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				switch e.Keysym.Sym {
				case sdl.K_ESCAPE:
					h.running = false
					return
				case sdl.K_UP:
					if h.quitSel {
						h.quitSel = false
						h.selected = len(gameList) - 1
					} else {
						h.selected = (h.selected - 1 + len(gameList)) % len(gameList)
					}
				case sdl.K_DOWN:
					if h.selected == len(gameList)-1 {
						h.quitSel = true
					} else {
						h.selected = (h.selected + 1) % len(gameList)
						h.quitSel = false
					}
				case sdl.K_RETURN:
					if h.quitSel {
						h.running = false
						return
					} else {
						h.LaunchGame(h.selected)
					}
				}
			}
		case *sdl.MouseMotionEvent:
			mx, my := int(e.X), int(e.Y)
			h.quitSel = false
			for i := range gameList {
				y := 130 + i*52
				rect := sdl.Rect{X: ScreenWidth/2 - 140, Y: int32(y), W: 280, H: 46}
				if mx >= int(rect.X) && mx <= int(rect.X+rect.W) &&
					my >= int(rect.Y) && my <= int(rect.Y+rect.H) {
					h.selected = i
					break
				}
			}
			quitRect := sdl.Rect{X: ScreenWidth/2 - 60, Y: int32(130 + len(gameList)*52 + 20), W: 120, H: 40}
			if mx >= int(quitRect.X) && mx <= int(quitRect.X+quitRect.W) &&
				my >= int(quitRect.Y) && my <= int(quitRect.Y+quitRect.H) {
				h.quitSel = true
			}
		case *sdl.MouseButtonEvent:
			if e.Type == sdl.MOUSEBUTTONDOWN && e.Button == sdl.BUTTON_LEFT {
				mx, my := int(e.X), int(e.Y)
				quitRect := sdl.Rect{X: ScreenWidth/2 - 60, Y: int32(130 + len(gameList)*52 + 20), W: 120, H: 40}
				if mx >= int(quitRect.X) && mx <= int(quitRect.X+quitRect.W) &&
					my >= int(quitRect.Y) && my <= int(quitRect.Y+quitRect.H) {
					h.running = false
					return
				}
				for i := range gameList {
					y := 130 + i*52
					rect := sdl.Rect{X: ScreenWidth/2 - 140, Y: int32(y), W: 280, H: 46}
					if mx >= int(rect.X) && mx <= int(rect.X+rect.W) &&
						my >= int(rect.Y) && my <= int(rect.Y+rect.H) {
						h.LaunchGame(i)
						break
					}
				}
			}
		}
	}
}

func (h *GameHub) LaunchGame(index int) {
    info := &gameList[index]
    info.instance = info.New(h.renderer)
	info.instance.Reset()
    info.instance.Run()
}

func (h *GameHub) Draw() {
	h.renderer.SetDrawColor(ColorBg.R, ColorBg.G, ColorBg.B, ColorBg.A)
	h.renderer.Clear()

	if h.font != nil {
		surface, err := h.font.RenderUTF8Solid("GAME HUB", ColorWhite)
		if err == nil {
			defer surface.Free()
			texture, err := h.renderer.CreateTextureFromSurface(surface)
			if err == nil {
				defer texture.Destroy()
				_, _, w, h2, _ := texture.Query()
				x := int32(ScreenWidth/2 - int(w)/2)
				y := int32(55 - int(h2)/2)
				h.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: w, H: h2})
			}
		}
	}

	for i, info := range gameList {
		y := int32(130 + i*52)
		rect := sdl.Rect{X: ScreenWidth/2 - 140, Y: y, W: 280, H: 46}

		if i == h.selected {
			h.renderer.SetDrawColor(ColorHoverBg.R, ColorHoverBg.G, ColorHoverBg.B, ColorHoverBg.A)
			games.DrawRoundedRect(h.renderer, rect.X, rect.Y, rect.W, rect.H, 12)

			h.renderer.SetDrawColor(ColorHover.R, ColorHover.G, ColorHover.B, ColorHover.A)
			games.DrawRoundedRectBorder(h.renderer, rect.X, rect.Y, rect.W, rect.H, 12, 3)

			h.renderer.SetDrawColor(ColorHover.R, ColorHover.G, ColorHover.B, 80)
			glowRect := sdl.Rect{X: rect.X - 4, Y: rect.Y - 4, W: rect.W + 8, H: rect.H + 8}
			games.DrawRoundedRect(h.renderer, glowRect.X, glowRect.Y, glowRect.W, glowRect.H, 14)
		} else {
			h.renderer.SetDrawColor(ColorButton.R, ColorButton.G, ColorButton.B, ColorButton.A)
			games.DrawRoundedRect(h.renderer, rect.X, rect.Y, rect.W, rect.H, 12)

			h.renderer.SetDrawColor(ColorBorder.R, ColorBorder.G, ColorBorder.B, ColorBorder.A)
			games.DrawRoundedRectBorder(h.renderer, rect.X, rect.Y, rect.W, rect.H, 12, 2)
		}

		if h.fontSmall != nil {
			var textColor sdl.Color
			if i == h.selected {
				textColor = ColorWhite
			} else {
				textColor = sdl.Color{R: 200, G: 200, B: 210, A: 255}
			}
			surface, err := h.fontSmall.RenderUTF8Solid(info.Name, textColor)
			if err == nil {
				defer surface.Free()
				texture, err := h.renderer.CreateTextureFromSurface(surface)
				if err == nil {
					defer texture.Destroy()
					_, _, w, h2, _ := texture.Query()
					x := int32(ScreenWidth/2 - int(w)/2)
					yy := y + int32(46/2-int(h2)/2)
					h.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: yy, W: w, H: h2})
				}
			}
		}
	}

	quitRect := sdl.Rect{X: ScreenWidth/2 - 60, Y: int32(130 + len(gameList)*52 + 20), W: 120, H: 40}
	var qColor sdl.Color
	if h.quitSel {
		qColor = ColorQuitHov
	} else {
		qColor = ColorQuit
	}
	h.renderer.SetDrawColor(qColor.R, qColor.G, qColor.B, qColor.A)
	games.DrawRoundedRect(h.renderer, quitRect.X, quitRect.Y, quitRect.W, quitRect.H, 10)

	borderColor := sdl.Color{R: 150, G: 40, B: 40, A: 255}
	if h.quitSel {
		borderColor = sdl.Color{R: 220, G: 60, B: 60, A: 255}
	}
	h.renderer.SetDrawColor(borderColor.R, borderColor.G, borderColor.B, borderColor.A)
	games.DrawRoundedRectBorder(h.renderer, quitRect.X, quitRect.Y, quitRect.W, quitRect.H, 10, 2)

	if h.fontSmall != nil {
		surface, err := h.fontSmall.RenderUTF8Solid("EXIT", ColorWhite)
		if err == nil {
			defer surface.Free()
			texture, err := h.renderer.CreateTextureFromSurface(surface)
			if err == nil {
				defer texture.Destroy()
				_, _, w, h2, _ := texture.Query()
				x := int32(ScreenWidth/2 - int(w)/2)
				y := int32(130 + len(gameList)*52 + 20 + 40/2 - int(h2)/2)
				h.renderer.Copy(texture, nil, &sdl.Rect{X: x, Y: y, W: w, H: h2})
			}
		}
	}

	h.renderer.Present()
}
