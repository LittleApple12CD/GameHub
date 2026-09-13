package games

import "github.com/veandco/go-sdl2/sdl"

func DrawRoundedRect(renderer *sdl.Renderer, x, y, w, h, r int32) {
	if r > w/2 {
		r = w / 2
	}
	if r > h/2 {
		r = h / 2
	}
	if r <= 0 {
		renderer.FillRect(&sdl.Rect{X: x, Y: y, W: w, H: h})
		return
	}

	renderer.FillRect(&sdl.Rect{X: x + r, Y: y, W: w - 2*r, H: h})
	renderer.FillRect(&sdl.Rect{X: x, Y: y + r, W: w, H: h - 2*r})

	for i := int32(0); i <= r; i++ {
		for j := int32(0); j <= r; j++ {
			if i*i+j*j <= r*r {
				renderer.DrawPoint(x+r-i, y+r-j)
				renderer.DrawPoint(x+w-r+i-1, y+r-j)
				renderer.DrawPoint(x+r-i, y+h-r+j-1)
				renderer.DrawPoint(x+w-r+i-1, y+h-r+j-1)
			}
		}
	}
}

func DrawRoundedRectBorder(renderer *sdl.Renderer, x, y, w, h, r, border int32) {
	if r > w/2 {
		r = w / 2
	}
	if r > h/2 {
		r = h / 2
	}
	if r <= 0 {
		for b := int32(0); b < border; b++ {
			renderer.DrawRect(&sdl.Rect{X: x + b, Y: y + b, W: w - 2*b, H: h - 2*b})
		}
		return
	}

	for b := int32(0); b < border; b++ {
		radius := r - b
		if radius < 0 {
			radius = 0
		}
		renderer.DrawLine(x+radius, y+b, x+w-radius, y+b)
		renderer.DrawLine(x+radius, y+h-b-1, x+w-radius, y+h-b-1)
		renderer.DrawLine(x+b, y+radius, x+b, y+h-radius)
		renderer.DrawLine(x+w-b-1, y+radius, x+w-b-1, y+h-radius)

		if radius > 0 {
			for i := int32(0); i <= radius; i++ {
				for j := int32(0); j <= radius; j++ {
					dist := i*i + j*j
					if dist >= (radius-1)*(radius-1) && dist <= radius*radius {
						renderer.DrawPoint(x+radius-i, y+radius-j)
						renderer.DrawPoint(x+w-radius+i-1, y+radius-j)
						renderer.DrawPoint(x+radius-i, y+h-radius+j-1)
						renderer.DrawPoint(x+w-radius+i-1, y+h-radius+j-1)
					}
				}
			}
		}
	}
}