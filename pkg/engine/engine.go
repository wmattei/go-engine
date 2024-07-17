package engine

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

type Engine struct {
	hasControls bool
	game        Game
	window      *glfw.Window
	program     uint32
}

func NewEngine(game Game) *Engine {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	width, height := game.GetScreenSize()

	window, err := glfw.CreateWindow(width, height, "Minceraft", nil, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize gl:", err)
	}

	program := initOpenGL()
	gl.UseProgram(program)

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.Enable(gl.DEBUG_OUTPUT_SYNCHRONOUS)

	return &Engine{
		game:    game,
		window:  window,
		program: program,
	}
}

func (e *Engine) Start() {

	e.window.MakeContextCurrent()

	projectionLoc := gl.GetUniformLocation(e.program, gl.Str("projection\x00"))
	viewLoc := gl.GetUniformLocation(e.program, gl.Str("view\x00"))
	modelLoc := gl.GetUniformLocation(e.program, gl.Str("model\x00"))

	lastTime := time.Now()
	fpsTime := lastTime
	frameCount := 0
	fps := 0

	for !e.window.ShouldClose() {
		currentTime := time.Now()

		lastTime = currentTime
		frameCount++

		if currentTime.Sub(fpsTime).Seconds() >= 1.0 {
			fps = frameCount
			frameCount = 0
			fpsTime = currentTime
			e.window.SetTitle(fmt.Sprintf("FPS: %d", fps))
		}

		scene := NewScene()

		e.game.Update(float32(0))
		e.game.Render(&scene)

		HandleInput(e.window, scene.camera)

		if !e.hasControls {
			SetupControls(e.window, scene.camera)
			e.hasControls = true
		}

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		projectionMatrix := scene.camera.GetProjectionMatrix()
		viewMatrix := scene.camera.GetViewMatrix()
		gl.UniformMatrix4fv(projectionLoc, 1, false, &projectionMatrix[0][0])
		gl.UniformMatrix4fv(viewLoc, 1, false, &viewMatrix[0][0])

		for _, model := range scene.models {
			if model.NeedsRender {
				if model.vao == 0 {
					model.SetVao()
				}

				modelMatrix := model.GetModelMatrix()
				flat := modelMatrix.Flatten()
				gl.UniformMatrix4fv(modelLoc, 1, false, &flat[0])

				gl.BindVertexArray(model.vao)
				gl.DrawElements(gl.TRIANGLES, int32(model.trianglesCount), gl.UNSIGNED_INT, nil)
				gl.BindVertexArray(0)
			}
		}

		e.window.SwapBuffers()
		glfw.PollEvents()

	}
}

func initOpenGL() uint32 {
	program := gl.CreateProgram()

	initShaders(program)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength)
		gl.GetProgramInfoLog(program, logLength, nil, &log[0])
		panic("failed to link program:" + string(log))
	}

	destroyShaders()

	return program
}

func Destroy() {
	glfw.Terminate()
}
