package games

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type TicTacToeGame struct {
	width     int32
	height    int32
	offsetX   int32
	offsetY   int32
	cellSize  int32
	board     [3][3]string
	current   string
	winner    string
	gameOver  bool
	moveCount int32
}

func NewTicTacToeGame() *TicTacToeGame {
	g := &TicTacToeGame{
		width:   500,
		height:  500,
		current: "X",
	}
	g.cellSize = g.width / 3
	g.offsetX = (ScreenWidth - g.width) / 2
	g.offsetY = (ScreenHeight - g.height) / 2
	g.Reset()
	return g
}

func (g *TicTacToeGame) Reset() {
	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			g.board[r][c] = ""
		}
	}
	g.current = "X"
	g.winner = ""
	g.gameOver = false
	g.moveCount = 0
}

func (g *TicTacToeGame) checkWinner() string {
	for row := 0; row < 3; row++ {
		if g.board[row][0] != "" && g.board[row][0] == g.board[row][1] && g.board[row][1] == g.board[row][2] {
			return g.board[row][0]
		}
	}
	for col := 0; col < 3; col++ {
		if g.board[0][col] != "" && g.board[0][col] == g.board[1][col] && g.board[1][col] == g.board[2][col] {
			return g.board[0][col]
		}
	}
	if g.board[0][0] != "" && g.board[0][0] == g.board[1][1] && g.board[1][1] == g.board[2][2] {
		return g.board[0][0]
	}
	if g.board[0][2] != "" && g.board[0][2] == g.board[1][1] && g.board[1][1] == g.board[2][0] {
		return g.board[0][2]
	}
	return ""
}

func (g *TicTacToeGame) Update() error {
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
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		x := (int32(mx) - g.offsetX) / g.cellSize
		y := (int32(my) - g.offsetY) / g.cellSize
		if x >= 0 && x < 3 && y >= 0 && y < 3 && g.board[y][x] == "" {
			g.board[y][x] = g.current
			g.moveCount++
			g.winner = g.checkWinner()
			if g.winner != "" {
				g.gameOver = true
			} else if g.moveCount == 9 {
				g.gameOver = true
			} else {
				if g.current == "X" {
					g.current = "O"
				} else {
					g.current = "X"
				}
			}
		}
	}
	return nil
}

func (g *TicTacToeGame) Draw(screen *ebiten.Image) {
	screen.Fill(ColorBg)

	DrawRoundedRect(screen, float32(g.offsetX), float32(g.offsetY),
		float32(g.width), float32(g.height), 12, color.RGBA{30, 30, 42, 255})

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			rx := float32(g.offsetX + int32(col)*g.cellSize + 2)
			ry := float32(g.offsetY + int32(row)*g.cellSize + 2)
			rw := float32(g.cellSize - 4)
			rh := float32(g.cellSize - 4)
			DrawRoundedRect(screen, rx, ry, rw, rh, 8, color.RGBA{50, 50, 65, 255})

			if g.board[row][col] != "" && FontBig != nil {
				var clr color.Color
				if g.board[row][col] == "X" {
					clr = color.RGBA{60, 200, 255, 255}
				} else {
					clr = color.RGBA{255, 200, 60, 255}
				}
				b := text.BoundString(FontBig, g.board[row][col])
				x := int(rx+rw/2) - b.Dx()/2
				y := int(ry+rh/2) + b.Dy()/2
				text.Draw(screen, g.board[row][col], FontBig, x, y, clr)
			}
		}
	}

	// 网格线
	for i := int32(1); i < 3; i++ {
		x := float32(g.offsetX+i*g.cellSize) - 1
		y := float32(g.offsetY+i*g.cellSize) - 1
		DrawRoundedRect(screen, x, float32(g.offsetY)+5, 3, float32(g.height)-10, 0,
			color.RGBA{80, 80, 100, 255})
		DrawRoundedRect(screen, float32(g.offsetX)+5, y, float32(g.width)-10, 3, 0,
			color.RGBA{80, 80, 100, 255})
	}

	statusY := int(g.offsetY+g.height) + 50
	if g.gameOver {
		var msg string
		var clr color.Color = ColorWhite
		if g.winner != "" {
			msg = "Player " + g.winner + " Wins!"
			clr = color.RGBA{0, 255, 0, 255}
		} else {
			msg = "Draw!"
			clr = color.RGBA{255, 255, 100, 255}
		}
		if FontBig != nil {
			b := text.BoundString(FontBig, msg)
			text.Draw(screen, msg, FontBig,
				ScreenWidth/2-b.Dx()/2, statusY, clr)
		}
		if FontSmall != nil {
			msg2 := "Press R to restart"
			b2 := text.BoundString(FontSmall, msg2)
			text.Draw(screen, msg2, FontSmall,
				ScreenWidth/2-b2.Dx()/2, statusY+40,
				color.RGBA{200, 200, 200, 255})
		}
	} else {
		msg := "Player " + g.current + "'s Turn"
		if FontBig != nil {
			b := text.BoundString(FontBig, msg)
			text.Draw(screen, msg, FontBig,
				ScreenWidth/2-b.Dx()/2, statusY, ColorWhite)
		}
	}
}

func (g *TicTacToeGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}