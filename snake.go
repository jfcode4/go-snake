package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"math/rand"
	"time"
)

type Scene struct {
	width    int
	height   int
	scale    int
	snake    []rl.Vector2
	snakeAbs []rl.Vector2
	velocity rl.Vector2
	apple    rl.Vector2
	textures struct {
		apple  rl.Texture2D
		head   rl.Texture2D
		body   rl.Texture2D
		corner rl.Texture2D
		tail   rl.Texture2D
	}
}

func NewScene(width, height, scale int) Scene {
	var s Scene
	s.width = width
	s.height = height
	s.scale = scale
	// textures
	s.textures.apple = rl.LoadTexture("assets/apple.png")
	s.textures.head = rl.LoadTexture("assets/head_right.png")
	s.textures.body = rl.LoadTexture("assets/body_horizontal.png")
	s.textures.corner = rl.LoadTexture("assets/body_bottomright.png")
	s.textures.tail = rl.LoadTexture("assets/tail_right.png")
	s.init()
	return s
}

func (s *Scene) init() {
	s.snake = []rl.Vector2{{X: 10, Y: 10}, {X: 9, Y: 10}}
	s.snakeAbs = []rl.Vector2{{X: 10, Y: 10}, {X: 9, Y: 10}}
	s.velocity = rl.Vector2{X: 1, Y: 0}
	s.placeApple()
	rl.SetTargetFPS(16)
}

func (s *Scene) placeApple() {
	for found := false; !found; {
		random := rand.Intn(s.width * s.height)
		s.apple.X = float32(random % s.width)
		s.apple.Y = float32(random / s.width)
		found = true
		for _, v := range s.snake {
			if v == s.apple {
				found = false
			}
		}
	}
}

func (s *Scene) Update() {
	// key presses
	velocity := s.velocity
	for {
		key := rl.GetKeyPressed()
		if key == 0 {
			break
		}
		if velocity.X == 0 {
			if key == rl.KeyA || key == rl.KeyLeft {
				s.velocity = rl.Vector2{X: -1, Y: 0}
			}
			if key == rl.KeyD || key == rl.KeyRight {
				s.velocity = rl.Vector2{X: 1, Y: 0}
			}
		}
		if velocity.Y == 0 {
			if key == rl.KeyW || key == rl.KeyUp {
				s.velocity = rl.Vector2{X: 0, Y: -1}
			}
			if key == rl.KeyS || key == rl.KeyDown {
				s.velocity = rl.Vector2{X: 0, Y: 1}
			}

		}

	}

	// move snake
	for i := len(s.snake) - 1; i >= 1; i-- {
		s.snake[i] = s.snake[i-1]
		s.snakeAbs[i] = s.snakeAbs[i-1]
	}
	s.snake[0].X += s.velocity.X
	s.snake[0].Y += s.velocity.Y
	s.snakeAbs[0].X += s.velocity.X
	s.snakeAbs[0].Y += s.velocity.Y

	// wrap around sides
	width := float32(s.width)
	height := float32(s.height)
	if s.snake[0].X == width {
		s.snake[0].X = 0
	}
	if s.snake[0].X == -1 {
		s.snake[0].X = width - 1
	}
	if s.snake[0].Y == height {
		s.snake[0].Y = 0
	}
	if s.snake[0].Y == -1 {
		s.snake[0].Y = height - 1
	}

	// snake eats the apple
	if rl.Vector2Equals(s.snake[0], s.apple) {
		s.snake = append(s.snake, s.snake[len(s.snake)-1])
		s.snakeAbs = append(s.snakeAbs, s.snake[len(s.snake)-1])
		s.placeApple()
		rl.SetTargetFPS(int32(16 + len(s.snake)/4))
	}

	// snake bumps into itself
	for i := 1; i < len(s.snake); i++ {
		if rl.Vector2Equals(s.snake[0], s.snake[i]) {
			// restart game
			s.gameOver()
			s.init()
		}
	}
}

func RenderRotated(t rl.Texture2D, position rl.Vector2, direction rl.Vector2) {
	angle := rl.Vector2Angle(rl.Vector2{X: 1, Y: 0}, direction) / math.Pi * 180
	offset := rl.Vector2Scale(rl.Vector2Subtract(rl.Vector2One(), rl.Vector2Rotate(rl.Vector2One(), angle*math.Pi/180)), 0.5)
	position = rl.Vector2Add(position, rl.Vector2Scale(offset, float32(t.Width)))
	rl.DrawTextureEx(t, position, angle, 1, rl.White)

}

func (s *Scene) Render() {
	scale := int32(s.scale)
	// background pattern
	rl.ClearBackground(rl.Color{173, 214, 68, 255})
	for y := 0; y < s.height; y += 1 {
		x := 0
		if y%2 == 1 {
			x = 1
		}
		for ; x < s.width; x += 2 {
			rl.DrawRectangle(int32(x)*scale, int32(y)*scale, scale, scale, rl.Color{166, 209, 60, 255})
		}
	}
	// render snake head
	direction := rl.Vector2Subtract(s.snakeAbs[0], s.snakeAbs[1])
	RenderRotated(s.textures.head, rl.Vector2Scale(s.snake[0], float32(scale)), direction)
	// render snake body
	for i := 1; i < len(s.snake)-1; i++ {
		//x := int32(s.snake[i].X) * scale
		//y := int32(s.snake[i].Y) * scale
		//rl.DrawRectangle(x, y, scale, scale, rl.Color{0, 255, 0, 255})
		if s.snake[i-1].X == s.snake[i+1].X || s.snake[i-1].Y == s.snake[i+1].Y {
			direction := rl.Vector2Subtract(s.snakeAbs[i-1], s.snakeAbs[i])
			RenderRotated(s.textures.body, rl.Vector2Scale(s.snake[i], float32(scale)), direction)
		} else {
			direction := rl.Vector2Subtract(s.snakeAbs[i-1], s.snakeAbs[i])
			if !rl.Vector2Equals(rl.Vector2Add(s.snakeAbs[i], rl.Vector2Rotate(direction, math.Pi/2)), s.snakeAbs[i+1]) {
				direction = rl.Vector2Subtract(s.snakeAbs[i+1], s.snakeAbs[i])
			}
			RenderRotated(s.textures.corner, rl.Vector2Scale(s.snake[i], float32(scale)), direction)
		}
		//rl.DrawTexture(s.textures.body, x, y, rl.White)
	}
	// render snake tail
	last := len(s.snake) - 1
	direction = rl.Vector2Subtract(s.snakeAbs[last], s.snakeAbs[last-1])
	RenderRotated(s.textures.tail, rl.Vector2Scale(s.snake[last], float32(scale)), direction)

	// render apple
	x := int32(s.apple.X) * scale
	y := int32(s.apple.Y) * scale
	//rl.DrawRectangle(x, y, scale, scale, rl.Color{255, 0, 0, 255})
	rl.DrawTexture(s.textures.apple, x, y, rl.White)

	// score
	rl.DrawText(fmt.Sprint("Score: ", len(s.snake)), 5, 5, 20, rl.RayWhite)
	//rl.DrawFPS(920, 5)
	rl.DrawText(fmt.Sprintf("%.0f FPS", 1/rl.GetFrameTime()), 1200, 5, 20, rl.Color{80, 80, 80, 255})
}

func (s Scene) gameOver() {
	rl.BeginDrawing()
	centerX := s.width * s.scale / 2
	centerY := s.height * s.scale / 2
	rl.DrawText("GAME OVER!", int32(centerX-120), int32(centerY-20), 40, rl.White)
	rl.EndDrawing()
	time.Sleep(2 * time.Second)
}

func main() {
	rl.InitWindow(1280, 720, "Snake Game")

	s := NewScene(32, 18, 40)

	for !rl.WindowShouldClose() {
		// update
		s.Update()

		// draw
		rl.BeginDrawing()
		s.Render()
		rl.EndDrawing()
	}
}
