package view

import (
	"fmt"
	"log"
	"time"
	"image"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl64"
	"gocv.io/x/gocv"

	// "github.com/nareix/joy4/av"
	"github.com/nareix/joy4/format/rtmp"
)

func init() {
	runtime.LockOSThread()
}

type Renderer struct {
	Camera *Camera
	Model  *ModelData
	Prog   uint32
	fbo    uint32
	isStop bool
	window *glfw.Window
}

const (
	vertexShaderSource = `
			#version 410

			layout(location = 0) in vec3 vp;
			layout(location = 1) in vec3 c;
			
			uniform mat4 ProjectionMatrix;
			uniform mat4 ModelViewMatrix;
			
			out vec4 color;

			void main() {
				gl_PointSize = 5;
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

var drawCnt int


func MatToPrt(mat mgl64.Mat4) *float64 {
	var res []float64
	for i := 0; i < 16; i++ {
		res = append(res, mat[i])
	}
	return &res[0]
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

	fmt.Println("init renderer")
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	r.Model.Init(vertexShaderSource, fragmentShaderSource)

	gl.GenFramebuffers(1, &r.fbo)
	drawCnt = 0
}

func (r *Renderer) GetFrame(w int, h int) gocv.Mat {
	r.Draw()
	// buffer := make([]byte, w*h*3)
	im := image.NewNRGBA(image.Rect(0, 0, w, h))
	gl.ReadPixels(0, 0, int32(w), int32(h), gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(im.Pix))
	mat, _ := ToRGB8(im)
	return mat
}

func ToRGB8(img image.Image) (gocv.Mat, error) {
	bounds := img.Bounds()
	x := bounds.Dx()
	y := bounds.Dy()
	bytes := make([]byte, 0, x*y*3)

	//don't get surprised of reversed order everywhere below
	for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
		for i := bounds.Min.X; i < bounds.Max.X; i++ {
			r, g, b, _ := img.At(i, j).RGBA()
			bytes = append(bytes, byte(b>>8), byte(g>>8), byte(r>>8))
		}
	}
	return gocv.NewMatFromBytes(y, x, gocv.MatTypeCV8UC3, bytes)
}

func (r *Renderer) Draw() {
	fmt.Println(drawCnt)
	drawCnt++
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.fbo)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(r.Prog)
	gl.Enable(gl.VERTEX_PROGRAM_POINT_SIZE)

	gl.BindVertexArray(r.Model.vao)
	var location int32
	location = gl.GetUniformLocation(r.Model.prog, gl.Str("ProjectionMatrix\x00"))
	gl.UniformMatrix4dv(location, 1, false, MatToPrt(r.Camera.PerspectiveMat))
	location = gl.GetUniformLocation(r.Model.prog, gl.Str("ModelViewMatrix\x00"))
	fmt.Println(location, r.Camera.ViewMat())
	gl.UniformMatrix4dv(location, 1, false, MatToPrt(r.Camera.ViewMat()))

	gl.BindBuffer(gl.ARRAY_BUFFER, r.Model.vvbo)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.BindBuffer(gl.ARRAY_BUFFER, r.Model.cvbo)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)

	gl.DrawArrays(gl.POINTS, 0, int32(len(r.Model.vertices)))

	glfw.PollEvents()
	r.window.SwapBuffers()

	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.Disable(gl.VERTEX_PROGRAM_POINT_SIZE)
}

func (r *Renderer) Streaming(w int, h int, conn *rtmp.Conn) {
	r.isStop = false
	id := 0
	for !r.isStop {
		// mt := r.GetFrame(w, h)
		// pkt := av.Packet{}
		// pkt.Data = mt.ToBytes()
		// if id % 30 == 0 {
		// 	pkt.IsKeyFrame = true
		// }
		// conn.WritePacket(pkt)
		r.Draw()
		time.Sleep(time.Duration(33) * time.Millisecond)
		id++
	}

}

func (r *Renderer) Stop() {
	r.isStop = true
}
