package main

// Loading other kinds of images
// Add this to pong, use images for the paddles
// See if you can speedup alpha blending

import (
	"fmt"
	"image/png"
	"os"
	"time"
	"vlado/game/noise"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth int = 800
const winHeight int = 600

type pos struct {
	x, y float32
}

type texture struct {
	pos
	pixels      []byte
	w, h, pitch int
	scale       float32
}

type rgba struct {
	r, g, b byte
}

// setPixel will set the pixel inside the pixel byte slice, 4 bytes per pixel, 3 are currently used
func setPixel(x, y int, c rgba, pixels []byte) error {
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

func (tex *texture) drawBilinearScaled(scaleX, scaleY float32, pixels []byte) {

	newWidth := int(float32(tex.w) * scaleX)
	newHeight := int(float32(tex.h) * scaleY)
	texW4 := tex.w * 4

	for y := 0; y < newHeight; y++ {

		fy := float32(y) / float32(newHeight) * float32(tex.h-1)
		fyi := int(fy)

		if y == 50 {
			fmt.Println(y, newHeight, tex.h, fy, fyi)
		}

		screenY := int(fy*scaleY) + int(tex.y)
		screenIndex := screenY*winWidth*4 + int(tex.x)*4
		ty := fy - float32(fyi)

		for x := 0; x < newWidth; x++ {
			fx := float32(x) / float32(newWidth) * float32(tex.w-1)
			fxi := int(fx)

			screenX := int(fx*scaleX) + int(tex.x)
			if screenX >= 0 && screenX < winWidth && screenY >= 0 && screenY < winHeight {

				c00i := fyi*texW4 + fxi*4
				c10i := fyi*texW4 + (fxi+1)*4
				c01i := (fyi+1)*texW4 + fxi*4
				c11i := (fyi+1)*texW4 + (fxi+1)*4

				tx := fx - float32(fxi)

				for i := 0; i < 4; i++ {
					c00 := float32(tex.pixels[c00i+i])
					c10 := float32(tex.pixels[c10i+i])
					c01 := float32(tex.pixels[c01i+i])
					c11 := float32(tex.pixels[c11i+i])

					pixels[screenIndex] = byte(blerp(c00, c10, c01, c11, tx, ty))
					screenIndex++
				}
			}
		}
	}
}

func flerp2(a, b, pct float32) float32 {
	return a + (b-a)*pct
}
func blerp(c00, c10, c01, c11, tx, ty float32) float32 {
	return flerp2(flerp2(c00, c10, tx), flerp2(c01, c11, tx), ty)
}

func (tex *texture) drawScaled(scaleX, scaleY float32, pixels []byte) {

	newWidth := int(float32(tex.w) * scaleX)
	newHeight := int(float32(tex.h) * scaleY)
	texW4 := tex.w * 4

	for y := 0; y < newHeight; y++ {

		fy := float32(y) / float32(newHeight) * float32(tex.h-1)
		fyi := int(fy)

		if y == 50 {
			fmt.Println(y, newHeight, tex.h, fy, fyi)
		}

		screenY := int(fy*scaleY) + int(tex.y)
		screenIndex := screenY*winWidth*4 + int(tex.x)*4

		for x := 0; x < newWidth; x++ {
			fx := float32(x) / float32(newWidth) * float32(tex.w-1)

			screenX := int(fx*scaleX) + int(tex.x)

			if screenX >= 0 && screenX < winWidth && screenY >= 0 && screenY < winHeight {
				fxi4 := int(fx) * 4

				pixels[screenIndex] = tex.pixels[fyi*texW4+fxi4]
				screenIndex++
				pixels[screenIndex] = tex.pixels[fyi*texW4+fxi4+1]
				screenIndex++
				pixels[screenIndex] = tex.pixels[fyi*texW4+fxi4+2]
				screenIndex++
				screenIndex++
			}
		}
	}

}

func (tex *texture) draw(pixels []byte) {
	for y := 0; y < tex.h; y++ {
		for x := 0; x < tex.w; x++ {
			screenY := y + int(tex.y)
			screenX := x + int(tex.x)

			if screenX >= 0 && screenX < winWidth && screenY >= 0 && screenY < winHeight {
				texIndex := y*tex.pitch + x*4
				screenIndex := screenY*winWidth*4 + screenX*4

				pixels[screenIndex] = tex.pixels[texIndex]
				pixels[screenIndex+1] = tex.pixels[texIndex+1]
				pixels[screenIndex+2] = tex.pixels[texIndex+2]
				pixels[screenIndex+3] = tex.pixels[texIndex+3]
			}
		}
	}
}

func flerp(b1 byte, b2 byte, pct float32) byte {
	return byte(float32(b1) + pct*(float32(b2)-float32(b1)))
}

func colorLerp(c1, c2 rgba, pct float32) rgba {
	return rgba{flerp(c1.r, c2.r, pct), flerp(c1.g, c2.g, pct), flerp(c1.b, c2.b, pct)}
}

func getGradient(c1, c2 rgba) []rgba {
	result := make([]rgba, 256)

	for i := range result {
		pct := float32(i) / float32(255)
		result[i] = colorLerp(c1, c2, pct)
	}

	return result
}

func getDualGradient(c1, c2, c3, c4 rgba) []rgba {
	result := make([]rgba, 256)

	for i := range result {
		pct := float32(i) / float32(255)

		if pct < 0.5 {
			result[i] = colorLerp(c1, c2, pct*float32(2))
		} else {
			result[i] = colorLerp(c3, c4, pct*float32(1.5)-float32(0.5))
		}
	}

	return result
}

func clamp(min, max, v int) int {
	if v < min {
		v = min
	} else if v > max {
		v = max
	}

	return v
}

func rescaleAndDraw(noise []float32, min, max float32, gradient []rgba, w, h int) []byte {

	result := make([]byte, w*h*4)
	scale := 255.0 / (max - min)
	offset := min * scale

	for i := range noise {
		noise[i] = noise[i]*scale - offset

		c := gradient[clamp(0, 255, int(noise[i]))]
		pidx := i * 4

		result[pidx] = c.r
		result[pidx+1] = c.g
		result[pidx+2] = c.b
	}

	return result
}

type aColorKey struct {
	src, dest, a int
}

// will draw with alpha blending
func (tex *texture) drawAlpha(pixels []byte) {

	// colorMap := make(map[aColorKey]int)

	for y := 0; y < tex.h; y++ {
		screenY := y + int(tex.y)
		for x := 0; x < tex.w; x++ {
			screenX := x + int(tex.x)

			if screenX >= 0 && screenX < winWidth && screenY >= 0 && screenY < winHeight {
				texIndex := y*tex.pitch + x*4
				screenIndex := screenY*winWidth*4 + screenX*4

				srcR := int(tex.pixels[texIndex])
				srcG := int(tex.pixels[texIndex+1])
				srcB := int(tex.pixels[texIndex+2])
				srcA := int(tex.pixels[texIndex+3])

				dstr := int(pixels[screenIndex])
				dstg := int(pixels[screenIndex+1])
				dstb := int(pixels[screenIndex+2])

				var rstR, rstG, rstB int

				// rKey := aColorKey{srcR, dstr, srcA}
				// gKey := aColorKey{srcR, dstr, srcA}
				// bKey := aColorKey{srcR, dstr, srcA}

				rstR = (srcR*255 + dstr*(255-srcA)) / 255
				rstG = (srcG*255 + dstg*(255-srcA)) / 255
				rstB = (srcB*255 + dstb*(255-srcA)) / 255

				// if colorMap[rKey] != 0 {
				// 	rstR = colorMap[rKey]
				// } else {
				// 	rstR = (srcR*255 + dstr*(255-srcA)) / 255
				// }
				// if colorMap[gKey] != 0 {
				// 	rstG = colorMap[gKey]
				// } else {
				// 	rstG = (srcG*255 + dstg*(255-srcA)) / 255
				// }
				// if colorMap[bKey] != 0 {
				// 	rstB = colorMap[bKey]
				// } else {
				// 	rstB = (srcB*255 + dstb*(255-srcA)) / 255
				// }

				pixels[screenIndex] = byte(rstR)
				pixels[screenIndex+1] = byte(rstG)
				pixels[screenIndex+2] = byte(rstB)
				// pixels[screenIndex+3] = tex.pixels[texIndex+3]
			}
		}
	}
}

func loadBallons() []texture {

	balloonStrs := []string{"balloon_red.png", "balloon_blue.png", "balloon_green.png"}
	balloonTextures := make([]texture, len(balloonStrs))

	for i, bstr := range balloonStrs {

		infile, err := os.Open(bstr)

		if err != nil {
			panic(err)
		}
		defer infile.Close()

		img, err := png.Decode(infile)
		if err != nil {
			panic(err)
		}

		w := img.Bounds().Max.X
		h := img.Bounds().Max.Y

		balloonPixels := make([]byte, w*h*4)
		bIndex := 0

		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r, g, b, a := img.At(x, y).RGBA()

				balloonPixels[bIndex] = byte(r / 256)
				bIndex++
				balloonPixels[bIndex] = byte(g / 256)
				bIndex++
				balloonPixels[bIndex] = byte(b / 256)
				bIndex++
				balloonPixels[bIndex] = byte(a / 256)
				bIndex++
			}
		}

		balloonTextures[i] = texture{pos{float32(i * 60), float32(i * 60)}, balloonPixels, w, h, w * 4, float32(1 + i)}
	}

	return balloonTextures
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

// main
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

	cloudNoise, min, max := noise.MakeNoise(noise.FMB, .009, .5, 3, 3, winWidth, winHeight)
	cloudGradient := getGradient(rgba{0, 0, 255}, rgba{255, 255, 255})
	cloudPixels := rescaleAndDraw(cloudNoise, min, max, cloudGradient, winWidth, winHeight)
	cloudTexture := texture{pos{0, 0}, cloudPixels, winWidth, winHeight, winWidth * 4, 1}

	pixels := make([]byte, winWidth*winHeight*4)
	balloonTextures := loadBallons()
	dir := 1

	for {
		frameStart := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		/*** Draw ***/

		// clear(pixels)
		cloudTexture.draw(pixels)

		for _, tex := range balloonTextures {
			tex.drawBilinearScaled(tex.scale, tex.scale, pixels)
		}

		balloonTextures[1].x += float32(1 * dir)
		if balloonTextures[1].x > 180 || balloonTextures[1].x < 0 {
			dir *= -1
		}

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		/*** Draw ***/

		elapsedTime := float32(time.Since(frameStart).Seconds() * 1000)
		// fmt.Println("ms per frame:", elapsedTime)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}

		sdl.Delay(16)
	}
}
