// Gl is a very hacky and lame set of OpenGL bindings simply used to verify that
// OpenGL can draw to a window.
package gl

/*
#include <GL/gl.h>
#include <GL/glu.h>

#cgo pkg-config: gl glu
*/
import "C"

import (
	"errors"
	"image/color"
)

func Init(left, right, bottom, top float64) {
	ClearColor(color.Black)
	C.glMatrixMode(C.GL_PROJECTION)
	C.gluOrtho2D(C.GLdouble(left), C.GLdouble(right), C.GLdouble(bottom), C.GLdouble(top))
	C.glMatrixMode(C.GL_MODELVIEW)
	C.glLoadIdentity()
}

func BeginQuads() {
	C.glBegin(C.GL_POLYGON)
}

func End() {
	C.glEnd()
}

func Vertex2(x, y float64) {
	C.glVertex2d(C.GLdouble(x), C.GLdouble(y))
}

func Color(col color.Color) {
	r, g, b, a := col.RGBA()
	C.glColor4f(
		C.GLfloat(r)/255.0,
		C.GLfloat(g)/255.0,
		C.GLfloat(b)/255.0,
		C.GLfloat(a)/255.0)
}

// ClearColor sets the color used to clear the color buffer.
func ClearColor(col color.Color) {
	r, g, b, a := col.RGBA()
	C.glClearColor(
		C.GLclampf(r)/255.0,
		C.GLclampf(g)/255.0,
		C.GLclampf(b)/255.0,
		C.GLclampf(a)/255.0)
}

const (
	// ColorBufferBit is a bit flag for Clear that specifies that the color buffer.
	ColorBufferBit = C.GL_COLOR_BUFFER_BIT

	// DepthBufferBit is a bit flag for Clear that specifies that the depth buffer.
	DepthBufferBit = C.GL_DEPTH_BUFFER_BIT
)

// Clear Clears the buffers specified by the bits.
func Clear(bits int) {
	C.glClear(C.GLbitfield(bits))
}

var errorStrings = map[C.GLenum]string{
	C.GL_NO_ERROR:                      "GL_NO_ERROR",
	C.GL_INVALID_ENUM:                  "GL_INVALID_ENUM",
	C.GL_INVALID_VALUE:                 "GL_INVALID_VALUE",
	C.GL_INVALID_OPERATION:             "GL_INVALID_OPERATION",
	C.GL_INVALID_FRAMEBUFFER_OPERATION: "GL_INVALID_FRAMEBUFFER_OPERATION",
	C.GL_OUT_OF_MEMORY:                 "GL_OUT_OF_MEMORY",
	C.GL_STACK_UNDERFLOW:               "GL_STACK_UNDERFLOW",
	C.GL_STACK_OVERFLOW:                "GL_STACK_OVERFLOW",
}

// Error returns an error from OpenGL, or nil if there was no error.
func Error() error {
	e := C.glGetError()
	if e == C.GL_NO_ERROR {
		return nil
	}
	return errors.New(errorStrings[e])
}
