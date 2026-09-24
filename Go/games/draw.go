package games

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawRoundedRect(screen *ebiten.Image, x, y, w, h, r float32, clr color.Color) {
	if r > w/2 {
		r = w / 2
	}
	if r > h/2 {
		r = h / 2
	}
	if r <= 0 {
		vector.DrawFilledRect(screen, x, y, w, h, clr, true)
		return
	}
	vector.DrawFilledRect(screen, x+r, y, w-2*r, h, clr, true)
	vector.DrawFilledRect(screen, x, y+r, w, h-2*r, clr, true)
	vector.DrawFilledCircle(screen, x+r, y+r, r, clr, true)
	vector.DrawFilledCircle(screen, x+w-r, y+r, r, clr, true)
	vector.DrawFilledCircle(screen, x+r, y+h-r, r, clr, true)
	vector.DrawFilledCircle(screen, x+w-r, y+h-r, r, clr, true)
}

func DrawRoundedRectBorder(screen *ebiten.Image, x, y, w, h, r, border float32, clr color.Color) {
	if r > w/2 {
		r = w / 2
	}
	if r > h/2 {
		r = h / 2
	}
	if r <= 0 {
		vector.StrokeRect(screen, x, y, w, h, border, clr, true)
		return
	}
	vector.StrokeLine(screen, x+r, y, x+w-r, y, border, clr, true)
	vector.StrokeLine(screen, x+r, y+h, x+w-r, y+h, border, clr, true)
	vector.StrokeLine(screen, x, y+r, x, y+h-r, border, clr, true)
	vector.StrokeLine(screen, x+w, y+r, x+w, y+h-r, border, clr, true)

	drawArc(screen, x+r, y+r, r, 180, 270, border, clr)
	drawArc(screen, x+w-r, y+r, r, 270, 360, border, clr)
	drawArc(screen, x+w-r, y+h-r, r, 0, 90, border, clr)
	drawArc(screen, x+r, y+h-r, r, 90, 180, border, clr)
}

func drawArc(screen *ebiten.Image, cx, cy, r, startDeg, endDeg, width float32, clr color.Color) {
	const segments = 16
	step := (endDeg - startDeg) / segments
	toRad := func(d float32) float64 { return float64(d) * math.Pi / 180.0 }
	prevX := cx + r*float32(math.Cos(toRad(startDeg)))
	prevY := cy + r*float32(math.Sin(toRad(startDeg)))
	for i := float32(1); i <= segments; i++ {
		ang := startDeg + step*i
		x := cx + r*float32(math.Cos(toRad(ang)))
		y := cy + r*float32(math.Sin(toRad(ang)))
		vector.StrokeLine(screen, prevX, prevY, x, y, width, clr, true)
		prevX, prevY = x, y
	}
}