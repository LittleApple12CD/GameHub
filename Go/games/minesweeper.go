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

type MinesweeperGame struct {
	rows         int32
	cols         int32
	mines        int32
	cellSize     int32
	width        int32
	height       int32
	offsetX      int32
	offsetY      int32
	board        [][]int
	revealed     [][]bool
	flagged      [][]bool
	gameOver     bool
	won          bool
	firstClick   bool
	running      bool
	tileHidden   *TileCache
	tileRevealed *TileCache
	tileFlagged  *TileCache
	tileMine     *TileCache
}

func NewMinesweeperGame() *MinesweeperGame {
	g := &MinesweeperGame{
		rows:     16,
		cols:     16,
		mines:    40,
		cellSize: 30,
	}
	g.width = g.cols * g.cellSize
	g.height = g.rows * g.cellSize
	g.offsetX = (ScreenWidth - g.width) / 2
	g.offsetY = (ScreenHeight - g.height) / 2
	g.tileHidden = NewTileCache(float32(g.cellSize-2), float32(g.cellSize-2), 4)
	g.tileRevealed = NewTileCache(float32(g.cellSize-2), float32(g.cellSize-2), 4)
	g.tileFlagged = NewTileCache(float32(g.cellSize-2), float32(g.cellSize-2), 4)
	g.tileMine = NewTileCache(float32(g.cellSize-2), float32(g.cellSize-2), 4)
	g.Reset()
	return g
}

func (g *MinesweeperGame) Reset() {
	g.board = make([][]int, g.rows)
	for i := range g.board {
		g.board[i] = make([]int, g.cols)
	}
	g.revealed = make([][]bool, g.rows)
	for i := range g.revealed {
		g.revealed[i] = make([]bool, g.cols)
	}
	g.flagged = make([][]bool, g.rows)
	for i := range g.flagged {
		g.flagged[i] = make([]bool, g.cols)
	}
	g.gameOver = false
	g.won = false
	g.firstClick = true
	g.running = true
}

func (g *MinesweeperGame) placeMines(safeX, safeY int32) {
	placed := int32(0)
	for placed < g.mines {
		x := int32(rand.Intn(int(g.cols)))
		y := int32(rand.Intn(int(g.rows)))
		if g.board[y][x] == -1 {
			continue
		}
		if x >= safeX-1 && x <= safeX+1 && y >= safeY-1 && y <= safeY+1 {
			continue
		}
		g.board[y][x] = -1
		placed++
	}
	for y := int32(0); y < g.rows; y++ {
		for x := int32(0); x < g.cols; x++ {
			if g.board[y][x] == -1 {
				continue
			}
			count := 0
			for dy := int32(-1); dy <= 1; dy++ {
				for dx := int32(-1); dx <= 1; dx++ {
					nx, ny := x+dx, y+dy
					if nx >= 0 && nx < g.cols && ny >= 0 && ny < g.rows {
						if g.board[ny][nx] == -1 {
							count++
						}
					}
				}
			}
			g.board[y][x] = count
		}
	}
}

func (g *MinesweeperGame) reveal(x, y int32) {
	if x < 0 || x >= g.cols || y < 0 || y >= g.rows {
		return
	}
	if g.revealed[y][x] || g.flagged[y][x] {
		return
	}
	g.revealed[y][x] = true
	if g.board[y][x] == -1 {
		g.gameOver = true
		return
	}
	if g.board[y][x] == 0 {
		for dy := int32(-1); dy <= 1; dy++ {
			for dx := int32(-1); dx <= 1; dx++ {
				g.reveal(x+dx, y+dy)
			}
		}
	}
	revealedCount := int32(0)
	for y2 := int32(0); y2 < g.rows; y2++ {
		for x2 := int32(0); x2 < g.cols; x2++ {
			if g.revealed[y2][x2] {
				revealedCount++
			}
		}
	}
	if revealedCount == g.rows*g.cols-g.mines {
		g.won = true
	}
}

func (g *MinesweeperGame) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Reset()
		return nil
	}

	if g.gameOver || g.won {
		return nil
	}

	mx, my := ebiten.CursorPosition()
	x := (int32(mx) - g.offsetX) / g.cellSize
	y := (int32(my) - g.offsetY) / g.cellSize
	inBounds := x >= 0 && x < g.cols && y >= 0 && y < g.rows

	if inBounds && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.firstClick {
			g.placeMines(x, y)
			g.firstClick = false
		}
		if !g.flagged[y][x] {
			g.reveal(x, y)
		}
	}
	if inBounds && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if !g.revealed[y][x] {
			g.flagged[y][x] = !g.flagged[y][x]
		}
	}
	return nil
}

func (g *MinesweeperGame) Draw(screen *ebiten.Image) {
	screen.Fill(ColorBg)

	DrawRoundedRect(screen, float32(g.offsetX), float32(g.offsetY),
		float32(g.width), float32(g.height), 12, color.RGBA{30, 30, 42, 255})

	numberColors := []color.RGBA{
		{0, 0, 0, 255},
		{0, 0, 255, 255},
		{0, 150, 0, 255},
		{255, 0, 0, 255},
		{0, 0, 180, 255},
		{150, 0, 0, 255},
		{0, 150, 150, 255},
		{0, 0, 0, 255},
	}

	for y := int32(0); y < g.rows; y++ {
		for x := int32(0); x < g.cols; x++ {
			rx := float32(g.offsetX + x*g.cellSize + 1)
			ry := float32(g.offsetY + y*g.cellSize + 1)

			if g.revealed[y][x] {
				if g.board[y][x] == -1 {
					DrawImageAt(screen, g.tileMine.Get(color.RGBA{200, 50, 50, 255}), rx, ry)
					cx := rx + float32(g.cellSize)/2
					cy := ry + float32(g.cellSize)/2
					vector.StrokeLine(screen, cx-8, cy, cx+8, cy, 2, color.Black, true)
					vector.StrokeLine(screen, cx, cy-8, cx, cy+8, 2, color.Black, true)
				} else {
					DrawImageAt(screen, g.tileRevealed.Get(color.RGBA{190, 190, 200, 255}), rx, ry)
					if g.board[y][x] > 0 && FontSmall != nil {
						text.Draw(screen, fmt.Sprint(g.board[y][x]), FontSmall,
							int(rx)+int(g.cellSize)/2-8,
							int(ry)+int(g.cellSize)/2+10,
							numberColors[g.board[y][x]])
					}
				}
			} else if g.flagged[y][x] {
				DrawImageAt(screen, g.tileFlagged.Get(color.RGBA{210, 210, 60, 255}), rx, ry)
				if FontSmall != nil {
					text.Draw(screen, "F", FontSmall,
						int(rx)+int(g.cellSize)/2-8,
						int(ry)+int(g.cellSize)/2+10,
						color.RGBA{255, 50, 50, 255})
				}
			} else {
				DrawImageAt(screen, g.tileHidden.Get(color.RGBA{80, 80, 105, 255}), rx, ry)
			}
		}
	}

	if g.gameOver || g.won {
		DrawOverlay(screen, int(g.offsetX), int(g.offsetY), int(g.width), int(g.height))
		var msg string
		var clr color.Color = ColorWhite
		if g.gameOver {
			msg = "GAME OVER"
		} else {
			msg = "YOU WIN!"
			clr = color.RGBA{0, 255, 0, 255}
		}
		b := text.BoundString(FontBig, msg)
		text.Draw(screen, msg, FontBig,
			int(g.offsetX)+int(g.width)/2-b.Dx()/2,
			int(g.offsetY)+int(g.height)/2-20, clr)
		msg2 := "Press R to restart"
		b2 := text.BoundString(FontSmall, msg2)
		text.Draw(screen, msg2, FontSmall,
			int(g.offsetX)+int(g.width)/2-b2.Dx()/2,
			int(g.offsetY)+int(g.height)/2+30, ColorWhite)
	}
}

func (g *MinesweeperGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}