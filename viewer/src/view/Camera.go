package view

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

type Camera struct {
	Position       mgl64.Vec3
	Rotation       mgl64.Vec3
	Fov            float64
	PerspectiveMat mgl64.Mat4

	Start mgl64.Quat
	Now   mgl64.Quat

	Down mgl64.Vec2

	Mode int

	Pan mgl64.Vec2
}

const (
	ROTATE = -1
	NONE   = 0
	PAN    = 1
)

func NewCamera(aspect float64) *Camera {
	camera := &Camera{Position: mgl64.Vec3{0, 0, 250}, Fov: 40}
	camera.PerspectiveMat = mgl64.Perspective(camera.Fov, aspect, 0.001, 1000)
	camera.Reset()
	camera.Spin(mgl64.Vec3{0.2, 0.4, 0})
	return camera
}

func MatArray(d mgl64.Mat4) *float32 {
	n := len(d)
	f := make([]float32, n)
	for i := 0; i < n; i++ {
		f[i] = float32(d[i])
	}
	return &f[0]
}

func (c *Camera) GetMouseNDC(w int, h int, mx int, my int) (x float64, y float64) {
	x = float64(mx) / float64(w) * 2 - 1
	y = float64(my) / float64(h) * 2 - 1
	return
}


func (c *Camera) ViewMat() mgl64.Mat4 {
	mat := mgl64.Ident4()
	mat = mat.Mul4(mgl64.Translate3D(-c.Position.X(), -c.Position.Y(), -c.Position.Z()))
	return mat.Mul4(c.GetMatrix())
}

func (c *Camera) Spin(vec mgl64.Vec3) {
	c.Start = c.Now.Mul(c.Start)
	c.Now = mgl64.QuatIdent()

	iw := vec.Len()
	if iw < 1 {
		iw = math.Sqrt(1 - iw)
	} else {
		iw = 0
	}

	newQ := mgl64.Quat{iw, vec}
	c.Start = newQ.Normalize().Mul(c.Start).Normalize()
}

func (c *Camera) GetMatrix() mgl64.Mat4 {
	qAll := c.Now.Mul(c.Start)
	qAll = qAll.Conjugate()
	return qAll.Mat4()
}

func (c *Camera) ComputeNow(now mgl64.Vec2) {
	if c.Mode == ROTATE {
		d := onUnitSphere(c.Down)
		m := onUnitSphere(now)
		c.Now.V = d.Cross(m)
		c.Now.W = d.Dot(m)
		c.Now = c.Now.Normalize()
	} else if c.Mode == PAN {
		d := now.Sub(c.Down).Mul(c.Position.Z())
		c.Position = c.Position.Add(c.Pan.Sub(d).Vec3(0))
		c.Pan = d
	}
}

func (c *Camera) MouseDown(point mgl64.Vec2) {
	c.Start = c.Now.Mul(c.Start)
	c.Now = mgl64.QuatIdent()

	c.Down = point
	c.Pan = mgl64.Vec2{0, 0}
}

func (c *Camera) Reset() {
	c.Start = mgl64.QuatIdent()
	c.Now = mgl64.QuatIdent()
}

func onUnitSphere(vec mgl64.Vec2) mgl64.Vec3 {
	mag := vec.Len()
	if mag > 1 {
		vec = vec.Normalize()
		return vec.Vec3(0)
	}
	return vec.Vec3(math.Sqrt(1 - mag))
}
