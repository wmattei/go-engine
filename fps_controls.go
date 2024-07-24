package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/wmattei/minceraft/pkg/engine"
)

func SetupControls(window *glfw.Window, camera *engine.PerspectiveCamera) {
	var lastX, lastY float64
	var firstMouse bool = true

	mouseCallback := func(window *glfw.Window, xpos float64, ypos float64) {
		if firstMouse {
			lastX = xpos
			lastY = ypos
			firstMouse = false
		}

		xoffset := xpos - lastX
		yoffset := lastY - ypos // Reversed since y-coordinates go from bottom to top

		lastX = xpos
		lastY = ypos

		// Adjust sensitivity as needed
		sensitivity := 0.1
		xoffset *= sensitivity
		yoffset *= sensitivity

		camera.Rotate(float32(yoffset), float32(xoffset))
	}

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPosCallback(mouseCallback)
}

func HandleInput(window *glfw.Window, camera *engine.PerspectiveCamera, dt float32) {
	if window.GetKey(glfw.KeyW) == glfw.Press {
		camera.ProcessKeyboard("FORWARD", dt)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		camera.ProcessKeyboard("BACKWARD", dt)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		camera.ProcessKeyboard("LEFT", dt)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		camera.ProcessKeyboard("RIGHT", dt)
	}
	if window.GetKey(glfw.KeySpace) == glfw.Press {
		camera.ProcessKeyboard("UP", dt)
	}
	if window.GetKey(glfw.KeyLeftShift) == glfw.Press {
		camera.ProcessKeyboard("DOWN", dt)
	}

	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}

}
