package ui

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

	"github.com/eaburns/gui/thread0"
	"github.com/eaburns/gui/ui/gl"
)

// A Canvas provides a high-level interface for drawing 2d graphics.
type Canvas struct {
	lineProg      *gl.Program
	lineVerts     *gl.Buffer
	solidRectProg *gl.Program
	imgRectProg   *gl.Program
	rectVerts     *gl.Buffer
}

// NewCanvas returns a new canvas.
func NewCanvas() *Canvas {
	c := new(Canvas)
	thread0.Do(func() {
		gl.Enable(gl.Texture2D)
		gl.Enable(gl.Blend)
		gl.BlendFunc(gl.SrcAlpha, gl.OneMinusSrcAlpha)
		c.lineProg = loadLineProg()
		c.lineVerts = gl.NewArrayBuffer()
		c.solidRectProg = loadSolidRectProg()
		c.imgRectProg = loadImgRectProg()
		c.rectVerts = makeRectBuffer()
	})
	return c
}

func loadLineProg() *gl.Program {
	v := strings.NewReader(lineVertShader)
	f := strings.NewReader(solidFragShader)
	p, err := gl.NewProgram(v, f)
	if err != nil {
		panic(err)
	}
	return p
}

func loadSolidRectProg() *gl.Program {
	v := strings.NewReader(rectVertShader)
	f := strings.NewReader(solidFragShader)
	p, err := gl.NewProgram(v, f)
	if err != nil {
		panic(err)
	}
	return p
}

func loadImgRectProg() *gl.Program {
	v := strings.NewReader(rectVertShader)
	f := strings.NewReader(texFragShader)
	p, err := gl.NewProgram(v, f)
	if err != nil {
		panic(err)
	}
	return p
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

// Close releases the resources for the canvas.
func (c *Canvas) Close() {
	thread0.Do(func() {
		c.lineProg.Delete()
		c.lineVerts.Delete()
		c.solidRectProg.Delete()
		c.imgRectProg.Delete()
		c.rectVerts.Delete()
	})
}

// Clear clears the canvas with the given color.
func (c *Canvas) Clear(col color.Color) {
	thread0.Do(func() {
		gl.ClearColor(col)
		gl.Clear(gl.ColorBufferBit)
	})
}

// StrokeLine strokes a line of the given color and width.
func (c *Canvas) StrokeLine(col color.Color, width float32, pts ...[2]float32) {
	thread0.Do(func() {
		gl.LineWidth(width)

		r, g, b, a := col.RGBA()
		c.lineProg.SetUniform("color",
			float32(r)/0xFFFF,
			float32(g)/0xFFFF,
			float32(b)/0xFFFF,
			float32(a)/0xFFFF,
		)

		data := make([]float32, len(pts)*2)
		for i, p := range pts {
			data[2*i] = p[0]
			data[(2*i)+1] = p[1]
		}
		c.lineVerts.SetData(gl.DynamicDraw, data...)

		c.lineVerts.Bind()
		c.lineProg.SetVertexAttributeData("vert", 2, 0, 0)
		c.lineProg.DrawArrays(gl.LineStrip, 0, len(pts))
		if err := gl.CheckError(); err != nil {
			panic(err)
		}
	})
}

// FillRect fills a rectangle with the given color.
// X and y specify the upper-right corner of the rectangle.
func (c *Canvas) FillRect(x, y, w, h float32, col color.Color) {
	thread0.Do(func() {
		r, g, b, a := col.RGBA()
		c.solidRectProg.SetUniform("color",
			float32(r)/0xFFFF,
			float32(g)/0xFFFF,
			float32(b)/0xFFFF,
			float32(a)/0xFFFF,
		)

		c.rectVerts.Bind()
		c.solidRectProg.SetVertexAttributeData("vert", 4, 0, 0)
		c.solidRectProg.SetUniform("loc", x, y)
		c.solidRectProg.SetUniform("size", w, h)
		c.solidRectProg.DrawArrays(gl.TriangleStrip, 0, 4)
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

// LoadPng returns an *Image loaded from a PNG file.
func LoadPng(path string) (*Image, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	img, err := png.Decode(r)
	if err != nil {
		return nil, err
	}
	return NewImage(img.(*image.NRGBA)), nil
}

// DrawImage draws an image at the given location on the canvas.
// X and y specify the upper-left corner of the image.
func (c *Canvas) DrawImage(x, y float32, img *Image) {
	thread0.Do(func() {
		img.tex.Bind(0)
		c.imgRectProg.SetUniform("tex", 0)

		c.rectVerts.Bind()
		c.imgRectProg.SetVertexAttributeData("vert", 4, 0, 0)
		c.imgRectProg.SetUniform("loc", x, y)
		c.imgRectProg.SetUniform("size", img.Width, img.Height)
		c.imgRectProg.DrawArrays(gl.TriangleStrip, 0, 4)
		if err := gl.CheckError(); err != nil {
			panic(err)
		}
	})
}

var (
	lineVertShader = `
		#version 120
		attribute vec2 vert;
		void main(){
			gl_Position = gl_ModelViewProjectionMatrix * vec4(vert, 0, 1);
		}`

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
