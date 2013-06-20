package ui

import (
	"image/color"
	"strings"

	"github.com/eaburns/gui/thread0"
	"github.com/eaburns/gui/ui/gl"
)

// A Canvas provides a high-level interface for drawing 2d graphics.
type Canvas struct {
	solid *gl.Program
	rect  *gl.Buffer
}

// NewCanvas returns a new canvas.
func NewCanvas() *Canvas {
	c := new(Canvas)
	thread0.Do(func() {
		c.solid = loadSolidShader()
		c.rect = makeRectBuffer()
	})
	return c
}

func makeRectBuffer() *gl.Buffer {
	buf := gl.NewArrayBuffer()
	buf.SetData(
		gl.StaticDraw,
		0, 0,
		1, 0,
		0, 1,
		1, 1,
	)
	return buf
}

func loadSolidShader() *gl.Program {
	v := strings.NewReader(rectShader)
	f := strings.NewReader(solidFragShader)
	solid, err := gl.NewProgram(v, f)
	if err != nil {
		panic(err)
	}
	return solid
}

// Close releases the resources for the canvas.
func (c *Canvas) Close() {
	c.solid.Delete()
	c.rect.Delete()
}

// Clear clears the canvas with the given color.
func (c *Canvas) Clear(col color.Color) {
	thread0.Do(func() {
		gl.ClearColor(col)
		gl.Clear(gl.ColorBufferBit)
	})
}

// FillRect fills a rectangle with the given color.
// X and y specify the upper-right corner of the rectangle.
func (c *Canvas) FillRect(x, y, w, h float32, col color.Color) {
	thread0.Do(func() {
		r, g, b, a := col.RGBA()
		c.solid.SetUniform("color",
			float32(r)/0xFFFF,
			float32(g)/0xFFFF,
			float32(b)/0xFFFF,
			float32(a)/0xFFFF,
		)

		c.rect.Bind()
		c.solid.SetVertexAttributeData("vert", 2, 0, 0)
		c.solid.SetUniform("loc", x, y)
		c.solid.SetUniform("size", w, h)
		c.solid.DrawArrays(gl.TriangleStrip, 0, 4)
		if err := gl.CheckError(); err != nil {
			panic(err)
		}
	})
}

var (
	rectShader = `
		#version 120
		attribute vec2 vert;
		uniform vec2 loc;
		uniform vec2 size;
		void main(){
			vec2 p = vec2(vert.x*size.x, vert.y*size.y);
			gl_Position = gl_ModelViewProjectionMatrix * vec4(p + loc, 0, 1);
		}`

	solidFragShader = `
		#version 120
		uniform vec4 color;
		void main()
		{
			gl_FragColor = color;
		}`
)
