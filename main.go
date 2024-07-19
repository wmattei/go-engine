package main

import (
	"fmt"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/wmattei/minceraft/pkg/engine"
)

const RENDER_DISTANCE = 1
const WIDTH = 1200
const HEIGHT = 780

func main() {
	window := engine.InitializeWindow(WIDTH, HEIGHT)
	program := engine.InitOpenGL()
	gl.UseProgram(program)

	world := NewSingleChunkWorld()
	// world := NewSingleBlockWorld()
	cam := engine.NewPerspectiveCamera(
		[3]float32{0, 40, 0},
		[3]float32{0, 1, 0},
		0,
		0,
		90,
		WIDTH/HEIGHT,
		0.1,
		1000.0,
	)

	modelLoc := gl.GetUniformLocation(program, gl.Str("model\x00"))
	viewLoc := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projLoc := gl.GetUniformLocation(program, gl.Str("projection\x00"))

	SetupControls(window, cam)

	lastTime := time.Now()
	fpsTime := lastTime
	frameCount := 0
	fps := 0

	for !window.ShouldClose() {
		currentTime := time.Now()

		lastTime = currentTime
		frameCount++

		if currentTime.Sub(fpsTime).Seconds() >= 1.0 {
			fps = frameCount
			frameCount = 0
			fpsTime = currentTime
			window.SetTitle(fmt.Sprintf("FPS: %d", fps))
		}

		HandleInput(window, cam)

		world.Update()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		view := cam.GetViewMatrix()
		viewFlatten := view.Flatten()

		projection := cam.GetProjectionMatrix()
		projectionFlatten := projection.Flatten()

		for _, chunk := range world.chunks {
			model := chunk.GetModelMatrix()
			flattenModel := model.Flatten()

			gl.UseProgram(program)
			gl.UniformMatrix4fv(modelLoc, 1, false, &flattenModel[0])
			gl.UniformMatrix4fv(viewLoc, 1, false, &viewFlatten[0])
			gl.UniformMatrix4fv(projLoc, 1, false, &projectionFlatten[0])

			chunk.Render()
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
