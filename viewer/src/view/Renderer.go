package view

import (
	"fmt"
	"log"
	"image"
	"image/jpeg"
	"bytes"
	"encoding/base64"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

type Renderer struct {
	Camera *Camera
	Model  *ModelData
	Prog   uint32
	fbo    uint32
	isStop bool
	window *glfw.Window
	Frame *image.RGBA
	Count int
}

const (
	VertexShaderSource = `
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

	FragmentShaderSource = `
			#version 410

			in vec4 color;

			out vec4 frag_colour;

			void main() {
				frag_colour = color;
			}
		` + "\x00"
)

func MatArray(d mgl64.Mat4) *float32 {
	n := len(d)
	f := make([]float32, n)
	for i := 0; i < n; i++ {
		f[i] = float32(d[i])
	}
	return &f[0]
}

func (r *Renderer) Init(w int, h int) {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(w, h, "test", nil, nil)
	if err != nil {
		panic(err)
	}
	r.window = window
	r.window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}
	
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
	
	if r.Model != nil {
		r.Model.Init(VertexShaderSource, FragmentShaderSource)
	}

	r.Frame = image.NewRGBA(image.Rect(0, 0, w, h))
	r.Count = 0
	
	fmt.Println("init renderer")
}

func (r *Renderer) GetFrame(w int, h int) {
	gl.ReadPixels(0, 0, int32(w), int32(h), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(r.Frame.Pix))
}

func (r *Renderer) Frame2Base64() string {
	buf := bytes.NewBuffer(nil)   
	jpeg.Encode(buf, r.Frame, &jpeg.Options{Quality: 80})
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func (r *Renderer) Draw(w int, h int) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	if r.Model != nil {
		r.Model.Draw(r.Camera.PerspectiveMat, r.Camera.ViewMat())
	}
	if r.Count % 10 == 0 {
		r.GetFrame(w, h)
	}

	glfw.PollEvents()
	r.window.SwapBuffers()
}

func (r *Renderer) Render(w int, h int) {
	for !r.window.ShouldClose() {
		r.Draw(w, h)
		r.Count++
	}
}

func (r *Renderer) Stop() {
	r.window.SetShouldClose(true)
	r.isStop = true
}
