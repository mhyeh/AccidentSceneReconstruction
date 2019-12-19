package view

import (
	"fmt"
	"strings"
	"encoding/binary"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl64"
)

const (
	vertexShaderSource = `
			#version 410

			layout(location = 0) in vec3 vp;
			layout(location = 1) in vec3 c;
			
			uniform mat4 ProjectionMatrix;
			uniform mat4 ModelViewMatrix;
			
			out vec4 color;

			void main() {
				gl_PointSize = 2;
				gl_Position = ProjectionMatrix * ModelViewMatrix * vec4(vp, 1.0);
				color = vec4(c, 1.0);
			}
		` + "\x00"

	fragmentShaderSource = `
			#version 410

			in vec4 color;

			out vec4 frag_colour;

			void main() {
				frag_colour = color;
			}
		` + "\x00"
)

type ModelData struct {
	vertices []float32
	normals []float32
	colors []float32

	vao uint32
	vvbo uint32
	cvbo uint32

	prog uint32
	projLoc int32
	viewLoc int32
}

func (m *ModelData) Init() {
	m.InitShader(vertexShaderSource, fragmentShaderSource)
	m.InitVAO()
	m.InitVBO()
}

func (m *ModelData) InitShader(vertexShaderSource string, fragmentShaderSource string) {
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

		log := strings.Repeat("\x00", int(logLength + 1))
		gl.GetProgramInfoLog(m.prog, logLength, nil, gl.Str(log))
		fmt.Errorf("failed to link program: %v", log)
		return;
	}

	m.projLoc = gl.GetUniformLocation(m.prog, gl.Str("ProjectionMatrix\x00"))
	m.viewLoc = gl.GetUniformLocation(m.prog, gl.Str("ModelViewMatrix\x00"))

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
	gl.BufferData(gl.ARRAY_BUFFER, binary.Size(m.vertices), gl.Ptr(&m.vertices[0]), gl.STATIC_DRAW)

	gl.GenBuffers(1, &m.cvbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.cvbo)
	gl.BufferData(gl.ARRAY_BUFFER, binary.Size(m.colors), gl.Ptr(&m.colors[0]), gl.STATIC_DRAW)
}

func (m *ModelData) Draw(proj mgl64.Mat4, view mgl64.Mat4) {
	gl.UseProgram(m.prog)
	gl.Enable(gl.VERTEX_PROGRAM_POINT_SIZE)
	gl.Enable(gl.DEPTH)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.FRONT_AND_BACK)

	gl.BindVertexArray(m.vao)
	gl.UniformMatrix4fv(m.projLoc, 1, false, MatArray(proj))

	gl.UniformMatrix4fv(m.viewLoc, 1, false, MatArray(view))

	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vvbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.cvbo)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)

	gl.DrawArrays(gl.POINTS, 0, int32(len(m.vertices) / 3))

	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.Disable(gl.VERTEX_PROGRAM_POINT_SIZE)
	gl.Disable(gl.DEPTH)
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.FRONT_AND_BACK)
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

		log := strings.Repeat("\x00", int(logLength + 1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}