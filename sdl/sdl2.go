package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth int = 800
const winHeight int = 600

type color struct {
	r, g, b byte
}

// setPixel will set the pixel inside the pixel byte slice, 4 bytes per pixel, 3 are currently used
func setPixel(x, y int, c color, pixels []byte) error {
	index := (y*int(winWidth) + x) * 4

	if index > len(pixels) || index < 0 {
		// simple string-based error
		return fmt.Errorf("the pixel index is not valid, index:  %d", index)
	}

	if index < len(pixels) && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}

	return nil
}

func exampleEq(x float64) float64 {
	// 2x- 3y = 1
	// 2x - 1 = 3y
	// 2x - 1 / 3 = y

	// example x=4, (4*2 -1) / 3

	y := (x*2 - 1) / 10

	return y
}

func main() {

	window, err := sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING,
		int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	for y := 0; y < winHeight; y++ {
		for x := 0; x < winWidth; x++ {
			var colored color
			eqresult := exampleEq(float64(x)) * 21

			if eqresult > float64(y) {
				colored = color{0, 255, 0}
			} else {
				colored = color{255, 0, 0}
			}

			_ = setPixel(x, y, colored, pixels) // ignore the error for now
		}
	}

	tex.Update(nil, pixels, winWidth*4)
	renderer.Copy(tex, nil, nil)
	renderer.Present()

	sdl.Delay(2000)

}
