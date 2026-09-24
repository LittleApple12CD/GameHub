package games

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

const (
	ScreenWidth  = 850
	ScreenHeight = 650
)

var (
	FontBig   font.Face
	FontSmall font.Face
)

var (
	ColorBg     = color.RGBA{25, 25, 35, 255}
	ColorWhite  = color.RGBA{255, 255, 255, 255}
	ColorRed    = color.RGBA{220, 60, 60, 255}
	ColorGreen  = color.RGBA{60, 220, 60, 255}
	ColorYellow = color.RGBA{255, 220, 60, 255}
)

func InitFonts() {
	var data []byte
	if b, err := os.ReadFile("assets/fonts/arial.ttf"); err == nil {
		data = b
	} else {
		data = goregular.TTF
	}
	f, err := opentype.Parse(data)
	if err != nil {
		log.Fatalf("parse font: %v", err)
	}
	FontBig, err = opentype.NewFace(f, &opentype.FaceOptions{Size: 48, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		log.Fatalf("new face big: %v", err)
	}
	FontSmall, err = opentype.NewFace(f, &opentype.FaceOptions{Size: 28, DPI: 72, Hinting: font.HintingFull})
	if err != nil {
		log.Fatalf("new face small: %v", err)
	}
}

type TileCache struct {
	imgs map[color.RGBA]*ebiten.Image
	w, h float32
	r    float32
}

func NewTileCache(w, h, r float32) *TileCache {
	return &TileCache{
		imgs: make(map[color.RGBA]*ebiten.Image),
		w:    w,
		h:    h,
		r:    r,
	}
}

func (c *TileCache) Get(clr color.RGBA) *ebiten.Image {
	if img, ok := c.imgs[clr]; ok {
		return img
	}
	img := ebiten.NewImage(int(c.w), int(c.h))
	DrawRoundedRect(img, 0, 0, c.w, c.h, c.r, clr)
	c.imgs[clr] = img
	return img
}

func (c *TileCache) GetWithBorder(clr color.RGBA, border color.Color) *ebiten.Image {
	if img, ok := c.imgs[clr]; ok {
		return img
	}
	img := ebiten.NewImage(int(c.w), int(c.h))
	
	DrawRoundedRect(img, 0, 0, c.w, c.h, c.r, border)

	bt := float32(2)
	DrawRoundedRect(img,
		bt, bt,
		c.w-2*bt, c.h-2*bt,
		c.r-bt,
		clr)

	c.imgs[clr] = img
	return img
}

func DrawImageAt(screen *ebiten.Image, img *ebiten.Image, x, y float32) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, op)
}