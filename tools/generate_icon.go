package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	// Create a 256x256 icon
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))

	// Blue background circle
	blue := color.RGBA{33, 150, 243, 255}
	green := color.RGBA{76, 175, 80, 255}
	darkBlue := color.RGBA{25, 118, 210, 255}
	white := color.RGBA{255, 255, 255, 255}

	// Draw background circle
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			dx, dy := x-128, y-128
			if dx*dx+dy*dy < 120*120 {
				img.Set(x, y, blue)
			}
		}
	}

	// Draw computer monitor (white rectangle)
	for y := 80; y < 160; y++ {
		for x := 70; x < 186; x++ {
			img.Set(x, y, white)
		}
	}

	// Draw screen (dark blue)
	for y := 85; y < 145; y++ {
		for x := 75; x < 181; x++ {
			img.Set(x, y, darkBlue)
		}
	}

	// Draw monitor stand
	for y := 160; y < 180; y++ {
		for x := 115; x < 141; x++ {
			img.Set(x, y, white)
		}
	}

	// Draw monitor base
	for y := 180; y < 188; y++ {
		for x := 100; x < 156; x++ {
			img.Set(x, y, white)
		}
	}

	// Draw connection nodes (small circles)
	drawCircle(img, 40, 60, 12, green)
	drawCircle(img, 216, 60, 12, green)
	drawCircle(img, 40, 196, 12, green)
	drawCircle(img, 216, 196, 12, green)

	// Draw connection lines
	drawLine(img, 52, 60, 70, 90, green)
	drawLine(img, 204, 60, 186, 90, green)
	drawLine(img, 52, 196, 70, 160, green)
	drawLine(img, 204, 196, 186, 160, green)

	// Draw cursor on screen
	for y := 95; y < 107; y++ {
		for x := 85; x < 89; x++ {
			img.Set(x, y, green)
		}
	}

	// Save as PNG
	f, err := os.Create("Icon.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}

func drawCircle(img *image.RGBA, cx, cy, r int, c color.RGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			dx, dy := x-cx, y-cy
			if dx*dx+dy*dy <= r*r {
				img.Set(x, y, c)
			}
		}
	}
}

func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.RGBA) {
	dx := x2 - x1
	dy := y2 - y1
	steps := max(abs(dx), abs(dy))

	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := x1 + int(float64(dx)*t)
		y := y1 + int(float64(dy)*t)

		// Draw thicker line
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				img.Set(x+dx, y+dy, c)
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
