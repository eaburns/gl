package ui

import (
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/eaburns/gui/thread0"
	"github.com/eaburns/gui/ui/gl"
)

// A Canvas provides a high-level interface for drawing 2d graphics.
type Canvas struct {
	lineProg      gl.Program
	lineVerts     gl.Buffer
	solidRectProg gl.Program
	imgRectProg   gl.Program
	rectVerts     gl.Buffer
}

// NewCanvas returns a new canvas.
func NewCanvas() *Canvas {
	c := new(Canvas)
	thread0.Do(func() {
		gl.Enable(gl.Blend)
		gl.BlendFunc(gl.SrcAlpha, gl.OneMinusSrcAlpha)
		c.lineProg = program(lineVertShader, solidFragShader)
		c.solidRectProg = program(rectVertShader, solidFragShader)
		c.imgRectProg = program(rectVertShader, texFragShader)

		if err := gl.GetError(); err != nil {
			panic(err)
		}

		bs := gl.GenBuffers(2)
		c.lineVerts = bs[0]
		c.rectVerts = bs[1]
		c.rectVerts.Bind(gl.ArrayBuffer)
		gl.BufferData(gl.ArrayBuffer,
			[]int8{
				0, 0, 0, 0,
				1, 0, 1, 0,
				0, 1, 0, 1,
				1, 1, 1, 1,
			},
			gl.StaticDraw)

		if err := gl.GetError(); err != nil {
			panic(err)
		}
	})
	return c
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

		if err := gl.GetError(); err != nil {
			panic(err)
		}
	})
}

// StrokeLine strokes a line of the given color and width.
func (c *Canvas) StrokeLine(col color.Color, width float32, pts ...[2]float32) {
	thread0.Do(func() {
		gl.LineWidth(width)

		c.lineProg.Use()

		r, g, b, a := col.RGBA()
		c.lineProg.GetUniformLocation("color").Uniform(
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
		c.lineVerts.Bind(gl.ArrayBuffer)
		gl.BufferData(gl.ArrayBuffer, data, gl.DynamicDraw)

		if err := gl.GetError(); err != nil {
			panic(err)
		}

		vattr := c.lineProg.GetAttribLocation("vert")
		vattr.Pointer(2, gl.Float, false, 0, 0)
		vattr.Enable()
		gl.DrawArrays(gl.LineStrip, 0, len(pts))
		vattr.Disable()

		if err := gl.GetError(); err != nil {
			panic(err)
		}

		gl.Program(0).Use()
	})
}

// FillRect fills a rectangle with the given color.
// X and y specify the upper-right corner of the rectangle.
func (c *Canvas) FillRect(x, y, w, h float32, col color.Color) {
	thread0.Do(func() {
		c.solidRectProg.Use()

		r, g, b, a := col.RGBA()
		c.solidRectProg.GetUniformLocation("color").Uniform(
			float32(r)/0xFFFF,
			float32(g)/0xFFFF,
			float32(b)/0xFFFF,
			float32(a)/0xFFFF,
		)
		c.solidRectProg.GetUniformLocation("loc").Uniform(x, y)
		c.solidRectProg.GetUniformLocation("size").Uniform(w, h)

		c.rectVerts.Bind(gl.ArrayBuffer)
		vattr := c.solidRectProg.GetAttribLocation("vert")
		vattr.Pointer(4, gl.Byte, false, 0, 0)
		vattr.Enable()
		gl.DrawArrays(gl.TriangleStrip, 0, 4)
		vattr.Disable()
		if err := gl.GetError(); err != nil {
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
	var i Image
	thread0.Do(func() {
		i.tex = gl.GenTextures(1)[0]
		i.tex.Bind(gl.Texture2D)
		w := img.Bounds().Dx()
		h := img.Bounds().Dy()
		gl.TexImage2D(gl.Texture2D, 0, gl.RGBA, w, h, 0, gl.RGBA, img.Pix)
		gl.TexParameter(gl.Texture2D, gl.TextureMagFilter, gl.Linear)
		gl.TexParameter(gl.Texture2D, gl.TextureMinFilter, gl.Linear)
		i.Width = float32(w)
		i.Height = float32(h)
		if err := gl.GetError(); err != nil {
			panic(err)
		}
	})
	return &i
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
		c.imgRectProg.Use()

		gl.ActiveTexture(0)
		img.tex.Bind(gl.Texture2D)
		c.imgRectProg.GetUniformLocation("tex").Uniform(0)
		c.imgRectProg.GetUniformLocation("loc").Uniform(x, y)
		c.imgRectProg.GetUniformLocation("size").Uniform(img.Width, img.Height)

		c.rectVerts.Bind(gl.ArrayBuffer)
		vattr := c.imgRectProg.GetAttribLocation("vert")
		vattr.Pointer(4, gl.Byte, false, 0, 0)
		vattr.Enable()
		gl.DrawArrays(gl.TriangleStrip, 0, 4)
		vattr.Disable()
		if err := gl.GetError(); err != nil {
			panic(err)
		}
	})
}

func program(vert, frag string) gl.Program {
	vsh := shader(gl.VertexShader, vert)
	defer vsh.Delete()
	fsh := shader(gl.FragmentShader, frag)
	defer fsh.Delete()

	p := gl.CreateProgram()
	if p == 0 {
		panic("Failed to create program")
	}
	p.AttachShader(vsh)
	p.AttachShader(fsh)
	p.Link()
	if p.Get(gl.LinkStatus) == 0 {
		panic("Failed to link program: " + p.GetInfoLog())
	}
	return p
}

func shader(kind gl.ShaderKind, src string) gl.Shader {
	sh := gl.CreateShader(kind)
	if sh == 0 {
		panic("Failed to create shader")
	}
	sh.Source(src)
	sh.Compile()
	if sh.Get(gl.CompileStatus) == 0 {
		panic("Failed to compile shader: " + sh.GetInfoLog())
	}
	return sh
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
