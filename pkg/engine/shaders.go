package engine

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader uint32

var SIMPLE_VERTEX_SHADER Shader
var SIMPLE_FRAGMENT_SHADER Shader

const simpleVertexShdr = `
    #version 410 core

    layout(location = 0) in vec3 inPosition;
    layout(location = 1) in vec4 inColor;
    layout(location = 2) in vec2 inTexCoord;
    layout(location = 3) in float inTexIndex;

    uniform mat4 model;
    uniform mat4 view;
    uniform mat4 projection;

    out vec4 color;
    out vec2 texCoord;
    out float texIndex;

    void main() {
        color = inColor;
        texCoord = inTexCoord;
        texIndex = inTexIndex;
        gl_Position = projection * view * model * vec4(inPosition, 1.0);
    }
` + "\x00"

const simpleFragmentShdr = `
    #version 410

    in vec4 color;
    in vec2 texCoord;
    in float texIndex;

    out vec4 frag_color;

    uniform sampler2D textures[6]; // Adjust size as needed

    void main() {
        int idx = int(texIndex);
        vec4 texColor = texture(textures[idx], texCoord);
        
        // If color alpha is 0.0, just use the texture color; otherwise, multiply by color
        if (color.a == 0.0) {
            frag_color = texColor;
        } else {
            frag_color = texColor * color;
        }
    }
` + "\x00"

func initShaders(program uint32) {
	vertexShader, err := compileShader(simpleVertexShdr, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragShader, err := compileShader(simpleFragmentShdr, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	SIMPLE_VERTEX_SHADER = Shader(vertexShader)
	SIMPLE_FRAGMENT_SHADER = Shader(fragShader)

	gl.AttachShader(program, uint32(SIMPLE_VERTEX_SHADER))
	gl.AttachShader(program, uint32(SIMPLE_FRAGMENT_SHADER))
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := make([]byte, logLength)
		gl.GetShaderInfoLog(shader, logLength, nil, &log[0])

		return 0, fmt.Errorf("failed to compile %v: %v", source, string(log))
	}

	return shader, nil
}

func destroyShaders() {
	gl.DeleteShader(uint32(SIMPLE_VERTEX_SHADER))
	gl.DeleteShader(uint32(SIMPLE_FRAGMENT_SHADER))
}
