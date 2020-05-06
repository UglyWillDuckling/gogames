package main

import (
	"fmt"
	"math/rand"
	"time"
	. "vlado/game/evolvingpictures/apt"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight, winDepth int = 800, 600, 100

type pos struct {
	x, y float32
}

type audioState struct {
	explosionBytes []byte
	deviceID       sdl.AudioDeviceID
	audioSpec      *sdl.AudioSpec
}

type mouseState struct {
	leftButton  bool
	rightButton bool
	x, y        int
}

func getMouseState() mouseState {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()

	var result mouseState
	result.x = int(mouseX)
	result.y = int(mouseY)

	result.leftButton = !(leftButton == 0)
	result.rightButton = !(rightButton == 0)

	return result
}

type rgba struct {
	r, g, b byte
}
type aColorKey struct {
	src, dest, a int
}

// setPixel will set the pixel inside the pixel byte slice, 4 bytes per pixel, 3 are currently used
func setPixel(x, y int, c rgba, pixels []byte) error {
	index := (y*int(winWidth) + x) * 4

	if index < len(pixels) && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}

	return nil
}

func pixelsToTexture(renderer *sdl.Renderer, pixels []byte, w, h int) *sdl.Texture {
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h))
	if err != nil {
		panic(err)
	}

	tex.Update(nil, pixels, w*4) // 4 bytes per pixel
	return tex
}

// APTToTexture does this
func APTToTexture(node1, node2, node3 Node, w, h int, renderer *sdl.Renderer) *sdl.Texture {
	// -1.0 and 1.0
	scale := float32(255 / 2)
	offset := float32(-1.0 * scale)

	pixels := make([]byte, w*h*4)
	pixelIndex := 0

	fmt.Println("scale", scale, "offset", offset)

	for yi := 0; yi < h; yi++ {
		y := float32(yi)/float32(h)*2 - 1

		for xi := 0; xi < w; xi++ {
			x := float32(xi)/float32(w)*2 - 1

			c := node1.Eval(x, y)
			c2 := node2.Eval(x, y)
			c3 := node3.Eval(x, y)

			// color := 360
			green := c*scale - offset
			red := c2*scale - offset
			blue := c3*scale - offset

			pixels[pixelIndex] = byte(blue)
			pixelIndex++
			pixels[pixelIndex] = byte(green) // byte(c) // green plus
			pixelIndex++
			pixels[pixelIndex] = byte(red) // red sine  byte(c*scale - offset)
			pixelIndex++

			pixelIndex++ // alpha

			// fmt.Println("x", x, "y", y, "c", c, "green color", c*scale-offset, "red color", c2*scale-offset)
		}
	}

	return pixelsToTexture(renderer, pixels, w, h)
}

func main() {
	/* !!! required for audio !!! */
	sdl.InitSubSystem(sdl.INIT_AUDIO)
	/* !!! required for audio !!! */

	window, err := sdl.CreateWindow("Evolving Pictures", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
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

	// var audioSpec sdl.AudioSpec
	// explosionBytes, _ := sdl.LoadWAV("explode.wav")
	// audioID, err := sdl.OpenAudioDevice("", false, &audioSpec, nil, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// defer sdl.FreeWAV(explosionBytes)
	// audioState := audioState{explosionBytes, audioID, &audioSpec}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	var elapsedTime float32
	currentMouseState := getMouseState()

	rand.Seed(time.Now().UTC().UnixNano())

	aptR := GetRandomNode()
	aptG := GetRandomNode()
	aptB := GetRandomNode()

	num := rand.Intn(20)
	for i := 0; i < num; i++ {
		aptR.AddRandom(GetRandomNode())
	}

	num = rand.Intn(20)
	for i := 0; i < num; i++ {
		aptG.AddRandom(GetRandomNode())
	}

	num = rand.Intn(20)
	for i := 0; i < num; i++ {
		aptB.AddRandom(GetRandomNode())
	}

	for {
		_, nilCount := aptR.NodeCounts()
		if nilCount == 0 {
			break
		}
		aptR.AddRandom(GetRandomLeaf())
	}
	for {
		_, nilCount := aptG.NodeCounts()
		if nilCount == 0 {
			break
		}
		aptG.AddRandom(GetRandomLeaf())
	}
	for {
		_, nilCount := aptB.NodeCounts()
		if nilCount == 0 {
			break
		}
		aptB.AddRandom(GetRandomLeaf())
	}

	tex := APTToTexture(aptR, aptG, aptB, winWidth, winHeight, renderer)

	fmt.Println("R: ", aptR)
	fmt.Println("G: ", aptG)
	fmt.Println("B: ", aptB)

	for {
		frameStart := time.Now()
		currentMouseState = getMouseState()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return

			case *sdl.TouchFingerEvent:
				if e.Type == sdl.FINGERDOWN {

					touchX := int(e.X * float32(winWidth))
					touchY := int(e.Y * float32(winHeight))
					currentMouseState.x = touchX
					currentMouseState.y = touchY
					currentMouseState.leftButton = true
				}
			}
		}

		renderer.Copy(tex, nil, nil)
		renderer.Present()
		/***** DRAW *****/

		elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		}
	}
}
