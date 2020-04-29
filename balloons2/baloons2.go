package main

// Homeworks

// 2. basic Collision Detection
//  treat  balloons as spheres - only check distance between center of balloons, use vec3, distance is less then the radius of the balloon
//  Minimum translation vector is important

import (
	"fmt"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"
	"vlado/game/noise"
	. "vlado/game/vec3"

	"github.com/veandco/go-sdl2/sdl"
)

func insertionSortBalloonsZ(items []*balloon) []*balloon {
	for i := range items {
		for j := i; j > 0; j-- {
			if items[j-1].pos.Z-items[j].pos.Z > 2.5 {
				v := items[j]
				items[j] = items[j-1]
				items[j-1] = v
			} else {
				continue
			}
		}
	}
	return items
}

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

type balloon struct {
	tex  *sdl.Texture
	pos  Vector3
	dir  Vector3
	w, h int

	exploding         bool
	exploded          bool
	explosionStart    time.Time
	explosionInterval float32
	explosionTexture  *sdl.Texture
}

func newBalloon(tex *sdl.Texture, pos, dir Vector3, explosionTexture *sdl.Texture) *balloon {
	_, _, w, h, err := tex.Query()

	if err != nil {
		panic(err)
	}

	return &balloon{tex, pos, dir, int(w), int(h), false, false, time.Now(), 35, explosionTexture}
}

type balloonArray []*balloon

func (balloons balloonArray) Len() int {
	return len(balloons)
}
func (balloons balloonArray) Swap(i, j int) {
	balloons[i], balloons[j] = balloons[j], balloons[i]
}
func (balloons balloonArray) Less(i, j int) bool {
	diff := balloons[i].pos.Z - balloons[j].pos.Z
	return diff < -2.5
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

func pixelsToTexture(renderer *sdl.Renderer, pixels []byte, w, h int) *sdl.Texture {
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h))

	if err != nil {
		panic(err)
	}

	tex.Update(nil, pixels, w*4) // 4 bytes per pixel

	return tex
}

func (balloon *balloon) getScale() float32 {
	scale := (balloon.pos.Z/200 + 1) / 2

	return scale
}

func (balloon *balloon) getCircle() (x, y, r float32) {
	x = balloon.pos.X
	y = balloon.pos.Y - 30*balloon.getScale()
	r = float32(balloon.w) / 2 * balloon.getScale()

	return x, y, r
}

func updateBaloons(balloons []*balloon, elapsedTime float32,
	currentMouseState, prevMouseState mouseState, audioState *audioState) []*balloon {

	numAnimations := 16
	balloonClicked := false
	balloonsExploded := false
	for i := len(balloons) - 1; i >= 0; i-- {
		balloon := balloons[i]

		if balloon.exploding {
			animationElapsed := float32(time.Since(balloon.explosionStart).Seconds() * 1000)
			animationIndex := numAnimations - 1 - int(animationElapsed/balloon.explosionInterval)
			if animationIndex < 0 {
				balloon.exploding = false
				balloon.exploded = true
				balloonsExploded = true
			}
		}

		if !balloonClicked && !prevMouseState.leftButton && currentMouseState.leftButton {
			x, y, r := balloon.getCircle()

			mouseX := currentMouseState.x
			mouseY := currentMouseState.y

			xDiff := float32(mouseX) - x
			yDiff := float32(mouseY) - y
			dist := float32(math.Sqrt(float64(xDiff*xDiff + yDiff*yDiff)))

			if dist < r {
				balloonClicked = true
				sdl.ClearQueuedAudio(audioState.deviceID)
				sdl.QueueAudio(audioState.deviceID, audioState.explosionBytes)
				sdl.PauseAudioDevice(audioState.deviceID, false)
				balloon.exploding = true
				balloon.explosionStart = time.Now()
			}
		}

		p := Add(balloon.pos, Mult(balloon.dir, elapsedTime))
		// reverse direction
		if p.X < 0 || p.X > float32(winWidth) {
			balloon.dir.X = -balloon.dir.X
		}
		if p.Y < 0 || p.Y > float32(winHeight) {
			balloon.dir.Y = -balloon.dir.Y
		}
		if p.Z < 0 || p.Z > float32(winDepth) {
			balloon.dir.Z = -balloon.dir.Z
		}

		balloon.pos = Add(balloon.pos, Mult(balloon.dir, elapsedTime))
	}

	if balloonsExploded {
		filteredBalloons := balloons[0:0]

		for _, balloon := range balloons {
			if !balloon.exploded {
				filteredBalloons = append(filteredBalloons, balloon)
			}
		}
		balloons = filteredBalloons
	}

	// Collision Detection
	var zIndexCurrent float32
	var zBalloons []*balloon

	for i := len(balloons) - 1; i >= 0; i-- {
		b := balloons[i]
		balloonBounced := false

		if b.pos.Z != zIndexCurrent {
			zIndexCurrent = b.pos.Z
			zBalloons = nil // reset the slice
			zBalloons = []*balloon{b}
			continue
		}

		for _, zBalloon := range zBalloons {
			// compute the distance between vectors
			dist := Distance(b.pos, zBalloon.pos)
			// get the radius
			_, _, r := b.getCircle()

			if dist <= r {
				// bounce the balloons, mind the Minimum translation vector
				translationDistance := float32(math.Abs(float64(dist - r)))

				b.dir.X = -(b.dir.X)
				b.dir.Y = -(b.dir.Y)
				zBalloon.dir.Y = -(zBalloon.dir.Y)
				zBalloon.dir.Y = -(zBalloon.dir.Y)

				// move the balloon immediately
				b.pos.X += -(b.pos.X / float32(math.Abs(float64(b.pos.X)))) * translationDistance
				b.pos.Y += -(b.pos.Y / float32(math.Abs(float64(b.pos.Y)))) * translationDistance

				zBalloon.pos.X += -(zBalloon.pos.X / float32(math.Abs(float64(zBalloon.pos.X)))) * translationDistance
				zBalloon.pos.Y += -(zBalloon.pos.Y / float32(math.Abs(float64(zBalloon.pos.Y)))) * translationDistance

				balloonBounced = true
				break
			}
		}

		if !balloonBounced {
			zBalloons = append(zBalloons, b)
		}
	}

	return balloons
}

func (balloon *balloon) draw(renderer *sdl.Renderer) {
	scale := balloon.getScale()
	newW := int32(float32(balloon.w) * scale)
	newH := int32(float32(balloon.h) * scale)
	x := int32(balloon.pos.X - float32(newW)/2)
	y := int32(balloon.pos.Y - float32(newH)/2)

	rect := &sdl.Rect{x, y, newW, newH}

	renderer.Copy(balloon.tex, nil, rect)

	if balloon.exploding {
		numAnimations := 16
		animationElapsed := float32(time.Since(balloon.explosionStart).Seconds() * 1000)
		animationIndex := numAnimations - 1 - int(animationElapsed/balloon.explosionInterval)
		/*	15 14 13 12
			11 10 9 8
			7 6 5 4
			3 2 1 0
		*/
		animationX := animationIndex % 4
		animationY := 64 * (animationIndex - animationX/4)
		animationX *= 64

		animationRect := &sdl.Rect{int32(animationX), int32(animationY), 64, 64}
		rect.X -= rect.W / 2
		rect.Y -= rect.H / 2
		rect.W *= 2
		rect.H *= 2
		renderer.Copy(balloon.explosionTexture, animationRect, rect)
	}
}

func imgFileToTexture(renderer *sdl.Renderer, filename string) *sdl.Texture {
	infile, err := os.Open(filename)

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

	pixels := make([]byte, w*h*4)
	bIndex := 0

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}
	tex := pixelsToTexture(renderer, pixels, w, h)

	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic(err)
	}

	return tex
}

func loadBallons(renderer *sdl.Renderer, numBalloons int) []*balloon {
	explosionTexture := imgFileToTexture(renderer, "explosion.png")

	balloonStrs := []string{"balloon_red.png", "balloon_blue.png", "balloon_green.png"}

	balloonTextures := make([]*sdl.Texture, len(balloonStrs))
	for i, bstr := range balloonStrs {
		balloonTextures[i] = imgFileToTexture(renderer, bstr)
	}

	balloons := make([]*balloon, numBalloons)
	for i := range balloons {
		tex := balloonTextures[i%3]
		pos := Vector3{rand.Float32() * float32(winWidth), rand.Float32() * float32(winHeight), 1}
		dir := Vector3{rand.Float32()*0.5 - .25, rand.Float32()*0.5 - .25, 0}

		balloons[i] = newBalloon(tex, pos, dir, explosionTexture)
	}

	return balloons
}

// main
func main() {
	/* !!! required for audio !!! */
	sdl.InitSubSystem(sdl.INIT_AUDIO)
	/* !!! required for audio !!! */

	window, err := sdl.CreateWindow("Exploding balloons", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
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

	var audioSpec sdl.AudioSpec

	explosionBytes, _ := sdl.LoadWAV("explode.wav")
	audioID, err := sdl.OpenAudioDevice("", false, &audioSpec, nil, 0)
	if err != nil {
		panic(err)
	}
	defer sdl.FreeWAV(explosionBytes)

	audioState := audioState{explosionBytes, audioID, &audioSpec}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	cloudNoise, min, max := noise.MakeNoise(noise.FMB, .009, .5, 3, 3, winWidth, winHeight)
	cloudGradient := getGradient(rgba{255, 0, 0}, rgba{255, 255, 255})
	cloudPixels := rescaleAndDraw(cloudNoise, min, max, cloudGradient, winWidth, winHeight)
	cloudTexture := pixelsToTexture(renderer, cloudPixels, winWidth, winHeight)

	balloons := loadBallons(renderer, 15)
	var elapsedTime float32
	currentMouseState := getMouseState()
	prevMouseState := currentMouseState

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

		balloons = updateBaloons(balloons, elapsedTime, currentMouseState, prevMouseState, &audioState)
		/***** DRAW *****/
		renderer.Copy(cloudTexture, nil, nil)

		balloons = insertionSortBalloonsZ(balloons)
		for _, balloon := range balloons {
			balloon.draw(renderer)
		}
		// if time.Since(startTime).Seconds() > 10 {
		// 	totalSortT := float64(0)
		// 	for _, v := range sortResult {
		// 		totalSortT += v
		// 	}
		// 	fmt.Println("total time for sort", totalSortT)
		// 	fmt.Println("average sort time", totalSortT/float64(len(sortResult)))
		// 	return
		// }
		renderer.Present()
		/***** DRAW *****/

		elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		}

		prevMouseState = currentMouseState
	}
}

/*
custom sort
total time for sort 0.15528486499999994
average sort time 5.489037292329443e-05
built in sort
0.16
*/
