package games

import "github.com/hajimehoshi/ebiten/v2"

type Game interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (int, int)
	Reset()
}