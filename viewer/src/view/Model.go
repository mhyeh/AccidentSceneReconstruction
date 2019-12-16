package view

import (
	"fmt"
	// "log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type ModelData struct {
	vertices []float32
	normals []float32
	colors []float32

	vao uint32
	vvbo uint32
	cvbo uint32

	prog uint32
}

func (m *ModelData) Init(vertexShaderSource string, fragmentShaderSource string) {
	m.InitShader(vertexShaderSource, fragmentShaderSource)
	m.InitVAO()
	m.InitVBO()
}

func (m *ModelData) InitShader(vertexShaderSource string, fragmentShaderSource string) {
	m.prog = gl.CreateProgram()

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	m.prog = gl.CreateProgram()
	gl.AttachShader(m.prog, vertexShader)
	gl.AttachShader(m.prog, fragmentShader)

	gl.LinkProgram(m.prog)

	var status int32
	gl.GetProgramiv(m.prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(m.prog, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(m.prog, logLength, nil, gl.Str(log))
		fmt.Errorf("failed to link program: %v", log)
		return;
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
}

func (m *ModelData) InitVAO() {
	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)
}

func (m *ModelData) InitVBO() {
	gl.GenBuffers(1, &m.vvbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vvbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4 * len(m.vertices), gl.Ptr(m.vertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &m.cvbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.cvbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4 * len(m.colors), gl.Ptr(m.colors), gl.STATIC_DRAW)
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

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}