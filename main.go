package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"runtime/pprof"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	minemath "github.com/wmattei/minceraft/math"
	"github.com/wmattei/minceraft/pkg/engine"

	_ "net/http/pprof"
)

const WIDTH = 1200
const HEIGHT = 780

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	window := engine.InitializeWindow(WIDTH, HEIGHT)
	program := engine.InitOpenGL()
	gl.UseProgram(program)

	world := NewWorld(8)

	// return
	// world := NewSingleChunkWorld()
	// world := NewSingleBlockWorld()

	cam := engine.NewPerspectiveCamera(
		[3]float32{0, 89, 0},
		[3]float32{0, 1, 0},
		0,
		0,
		math.Pi/2,
		float32(WIDTH)/float32(HEIGHT),
		0.01,
		1000,
	)

	frustum := engine.NewFrustum(cam)

	viewLoc := gl.GetUniformLocation(program, gl.Str("view\x00"))
	projLoc := gl.GetUniformLocation(program, gl.Str("projection\x00"))

	SetupControls(window, cam)

	glfw.SwapInterval(0)

	lastTime := time.Now()
	fpsTime := lastTime
	frameCount := 0
	fps := 0

	for !window.ShouldClose() {
		currentTime := time.Now()

		HandleInput(window, cam, float32(currentTime.Sub(lastTime).Seconds()))

		world.CheckCollisions(cam)
		world.Update(cam)

		lastTime = currentTime
		frameCount++

		if currentTime.Sub(fpsTime).Seconds() >= 1.0 {
			fps = frameCount
			frameCount = 0
			fpsTime = currentTime
			window.SetTitle(fmt.Sprintf("FPS: %d", fps))
		}

		view := cam.GetViewMatrix()
		viewFlatten := view.Flatten()
		projection := cam.GetProjectionMatrix()
		projectionFlatten := projection.Flatten()
		frustum.UpdateFrustum(minemath.MultiplyMatrices(projection, view))

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(program)
		gl.UniformMatrix4fv(viewLoc, 1, false, &viewFlatten[0])
		gl.UniformMatrix4fv(projLoc, 1, false, &projectionFlatten[0])

		world.Render(program, frustum)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
