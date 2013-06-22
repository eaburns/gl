package ui

import (
	"image"
	"image/color"
	"strings"

	"github.com/eaburns/gui/thread0"
	"github.com/eaburns/gui/ui/gl"
)

// A Canvas provides a high-level interface for drawing 2d graphics.
type Canvas struct {
	solid *gl.Program
	tex   *gl.Program
	rect  *gl.Buffer
}

// NewCanvas returns a new canvas.
func NewCanvas() *Canvas {
	c := new(Canvas)
	thread0.Do(func() {
		gl.Enable(gl.Texture2D)
		gl.Enable(gl.Blend)
		gl.BlendFunc(gl.SrcAlpha, gl.OneMinusSrcAlpha)
		c.solid = loadSolidShader()
		c.tex = loadTextureShader()
		c.rect = makeRectBuffer()
	})
	return c
}

func makeRectBuffer() *gl.Buffer {
	buf := gl.NewArrayBuffer()
	buf.SetData(
		gl.StaticDraw,
		0, 0, 0, 0,
		1, 0, 1, 0,
		0, 1, 0, 1,
		1, 1, 1, 1,
	)
	return buf
}

func loadSolidShader() *gl.Program {
	v := strings.NewReader(rectVertShader)
	f := strings.NewReader(solidFragShader)
	p, err := gl.NewProgram(v, f)
	if err != nil {
		panic(err)
	}
	return p
}

func loadTextureShader() *gl.Program {
	v := strings.NewReader(rectVertShader)
	f := strings.NewReader(texFragShader)
	p, err := gl.NewProgram(v, f)
	if err != nil {
		panic(err)
	}
	return p
}

// Close releases the resources for the canvas.
func (c *Canvas) Close() {
	thread0.Do(func() {
		c.solid.Delete()
		c.tex.Delete()
		c.rect.Delete()
	})
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
		c.solid.SetVertexAttributeData("vert", 4, 0, 0)
		c.solid.SetUniform("loc", x, y)
		c.solid.SetUniform("size", w, h)
		c.solid.DrawArrays(gl.TriangleStrip, 0, 4)
		if err := gl.CheckError(); err != nil {
			panic(err)
		}
	})
}

// An Image can be drawn to a canvas.
type Image struct {
	tex           gl.Texture
	Width, Height float32
}

// NewImage returns a new *Image from a normalized RGBA image.
func NewImage(img *image.NRGBA) *Image {
	var t gl.Texture
	thread0.Do(func() { t = gl.MakeImageTexture(img) })
	return &Image{
		tex:    t,
		Width:  float32(img.Bounds().Dx()),
		Height: float32(img.Bounds().Dy()),
	}
}

// DrawImage draws an image at the given location on the canvas.
// X and y specify the upper-left corner of the image.
func (c *Canvas) DrawImage(x, y float32, img *Image) {
	thread0.Do(func() {
		img.tex.Bind(0)
		c.tex.SetUniform("tex", 0)

		c.rect.Bind()
		c.tex.SetVertexAttributeData("vert", 4, 0, 0)
		c.tex.SetUniform("loc", x, y)
		c.tex.SetUniform("size", img.Width, img.Height)
		c.tex.DrawArrays(gl.TriangleStrip, 0, 4)
		if err := gl.CheckError(); err != nil {
			panic(err)
		}
	})
}

var (
	rectVertShader = `
		#version 120
		attribute vec4 vert;
		uniform vec2 loc;
		uniform vec2 size;
		varying vec2 texCoord;
		void main(){
			vec2 p = vec2(vert.x*size.x, vert.y*size.y);
			texCoord = vert.zw;
			gl_Position = gl_ModelViewProjectionMatrix * vec4(p + loc, 0, 1);
		}`

	solidFragShader = `
		#version 120
		uniform vec4 color;
		varying vec2 texCoord;
		void main()
		{
			gl_FragColor = color;
		}`

	texFragShader = `
		#version 120
		uniform sampler2D tex;
		varying vec2 texCoord;
		void main()
		{
			gl_FragColor = texture2D(tex, texCoord);
		}`
)
