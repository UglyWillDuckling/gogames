package main

// Homework
// GameOver State - win/lose
// make gameplay more interesting
// 2 Player vs Playing computer
// AI needs to be more imperfect
// Handling resizing of the window
// load bitmaps for our paddles/ball

import (
	"fmt"
	"math"
	"time"
	"vlado/game/noise"

	"github.com/veandco/go-sdl2/sdl"
)

// --
type gameState int

const (
	start gameState = iota
	play
)

var state = start

// --

var nums = [][]byte{
	{
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1},
	{
		1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1},
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1},
	{
		1, 1, 1,
		0, 0, 1,
		0, 1, 1,
		0, 0, 1,
		1, 1, 1},
	{
		1, 0, 0,
		1, 0, 0,
		1, 1, 1,
		0, 0, 1,
		0, 0, 1},
}

const winWidth int = 800
const winHeight int = 600

const PADDLEWIDTH = 20
const PADDLEHEIGHT = 100

func drawNumber(pos pos, color color, size int, num int, pixels []byte) {
	startX := int(pos.x) - (size*3)/2
	startY := int(pos.y) - (size*5)/2

	for i, v := range nums[num] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += size
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}

func flerp(b1 byte, b2 byte, pct float32) byte {
	return byte(float32(b1) + pct*(float32(b2)-float32(b1)))
}

func colorLerp(c1, c2 color, pct float32) color {
	return color{flerp(c1.r, c2.r, pct), flerp(c1.g, c2.g, pct), flerp(c1.b, c2.b, pct)}
}

func getGradient(c1, c2 color) []color {
	result := make([]color, 256)

	for i := range result {
		pct := float32(i) / float32(255)
		result[i] = colorLerp(c1, c2, pct)
	}

	return result
}

func getDualGradient(c1, c2, c3, c4 color) []color {
	result := make([]color, 256)

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

func rescaleAndDraw(noise []float32, min, max float32, gradient []color, w, h int) []byte {

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

func lerp(a float32, b float32, pct float32) float32 {
	return a + pct*(b-a)
}

func getCenter() pos {
	return pos{float32(winWidth / 2), float32(winHeight / 2)}
}

type gameobject interface {
	draw(pixels []byte)
}

type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius float32
	xv     float32
	yv     float32
	color  color
}

func (ball *ball) draw(pixels []byte) {
	//YAGNI - Ya Aint gonna need it

	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x+x), int(ball.y+y), ball.color, pixels)
			}
		}
	}
}

func (ball *ball) update(leftpaddle *paddle, rightpaddle *paddle, elapsedTime float32) {
	ball.x += ball.xv * elapsedTime
	ball.y += ball.yv * elapsedTime

	// reverse the ball when it hits the edge
	if (ball.y)-ball.radius < 0 || (ball.y)+ball.radius > float32(winHeight) {
		ball.yv = -ball.yv
	}

	// reset the ball
	if ball.x < 0 {
		rightpaddle.score++
		ball.pos = getCenter()
		state = start
		fmt.Println("restart")
	} else if ball.x+ball.radius > float32(winWidth) {
		leftpaddle.score++
		ball.pos = getCenter()
		state = start
	}

	ballLeftX := (ball.x) - ball.radius
	ballRightX := (ball.x) + ball.radius
	ballTopY := (ball.y) + ball.radius
	ballBottomY := (ball.y) - ball.radius

	leftPaddleREdge := (leftpaddle.x) + float32(leftpaddle.w/2)
	leftPaddleTop := (leftpaddle.y) - float32(leftpaddle.h/2)
	leftPaddleBottom := (leftpaddle.y) + float32(leftpaddle.w/2)

	if ballLeftX <= leftPaddleREdge {
		if ballTopY >= leftPaddleTop && ballBottomY <= leftPaddleBottom {

			// calculate the velocity based on the x difference
			// if the diff is 0 just bounce the ball back
			// if it isn't calculate the velocity in steps of 5

			xDiff := math.Abs(float64(ballLeftX - leftPaddleREdge))
			step := float64(8)  // hardcode to 8 for now
			parts := float64(8) // number of parts for the horizontal values

			if xDiff <= step {
				ball.xv = -ball.xv
			} else {
				currentPart := xDiff / step

				fmt.Println("new x velocity", float32(math.Abs(float64(ball.xv-float32(currentPart/parts)*ball.xv/float32(parts)))))

				ball.xv = float32(math.Abs(float64(ball.xv - float32(currentPart/parts)*ball.xv/float32(parts))))
				ball.yv = -ball.yv
			}

			ball.x += leftpaddle.w/2.0 + ball.radius
		}
	}

	if ballRightX > rightpaddle.x-rightpaddle.w/2 {
		if ballBottomY > rightpaddle.y-rightpaddle.h/2 && ball.y < rightpaddle.y+rightpaddle.h/2 {
			ball.xv = -ball.xv

			fmt.Println(ball.x, rightpaddle.w/2.0, ball.radius, ball.x-rightpaddle.w/2.0-ball.radius)
			ball.x = ball.x - rightpaddle.w/2.0 - ball.radius
			fmt.Println(ball.x)
		}
	}
}

type paddle struct {
	pos
	w     float32
	h     float32
	speed float32
	color color
	score int
}

func (paddle *paddle) update(keyState []uint8, controllerAxis int16, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= paddle.speed * elapsedTime
	}
	if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += paddle.speed * elapsedTime
	}

	if math.Abs(float64(controllerAxis)) > 1500 {
		pct := float32(controllerAxis) / 32767.0
		paddle.y += paddle.speed * pct * elapsedTime
	}
}

func (paddle *paddle) aiUpdate(ball *ball, elapsedTime float32) {
	paddle.y = ball.y
}

func (paddle *paddle) draw(pixels []byte) {
	startX := paddle.x - paddle.w/2
	startY := paddle.y - paddle.h/2

	for y := 0; y < int(paddle.h); y++ {
		for x := 0; x < int(paddle.w); x++ {
			setPixel(int(startX)+x, int(startY)+y, paddle.color, pixels)
		}
	}

	numX := lerp(paddle.x, getCenter().x, 0.2)

	drawNumber(pos{numX, 35}, paddle.color, 10, paddle.score, pixels)
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
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

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

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

	var controllerHandlers []*sdl.GameController
	for i := 0; i < sdl.NumJoysticks(); i++ {
		controllerHandlers = append(controllerHandlers, sdl.GameControllerOpen(i))
		defer controllerHandlers[i].Close()
	}

	pixels := make([]byte, winWidth*winHeight*4)

	player1 := paddle{pos{50, 100}, PADDLEWIDTH, PADDLEHEIGHT, 500, color{255, 255, 255}, 0}
	player2 := paddle{pos{float32(winWidth) - 50, 100}, PADDLEWIDTH, PADDLEHEIGHT, 300, color{255, 255, 255}, 0}

	ball := ball{pos{300, 300}, 20, 400, 400, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	noise, min, max := noise.MakeNoise(noise.FMB, .01, 0.5, 2, 3, winWidth, winHeight)
	gradient := getGradient(color{255, 0, 0}, color{0, 0, 0})

	noisePixels := rescaleAndDraw(noise, min, max, gradient, winWidth, winHeight)

	var frameStart time.Time
	var elapsedTime float32
	var controllerAxis int16

	for {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		for _, controller := range controllerHandlers {
			if controller != nil {
				controllerAxis = controller.Axis(sdl.CONTROLLER_AXIS_LEFTY)
			}
		}

		if state == play {
			player1.update(keyState, controllerAxis, elapsedTime)
			player2.aiUpdate(&ball, elapsedTime)
			ball.update(&player1, &player2, elapsedTime)
		} else if state == start {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 3 || player2.score == 3 {

					player1.score = 0
					player2.score = 0

				}
				state = play
			}
		}

		for i := range noisePixels {
			pixels[i] = noisePixels[i]
		}

		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())

		if elapsedTime < .005 {
			sdl.Delay(5 - uint32(elapsedTime*1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
	}

}
