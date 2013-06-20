// Package gl provides small, high-level, binding to OpenGL.
// Functions in this package can only be called safely by the main go routine.
package gl

/*
#include <GL/gl.h>

#cgo pkg-config: gl
*/
import "C"

import (
	"errors"
	"image/color"
)

// ClearColor sets the color used to clear the color buffer.
func ClearColor(col color.Color) {
	r, g, b, a := col.RGBA()
	C.glClearColor(
		C.GLclampf(r)/255.0,
		C.GLclampf(g)/255.0,
		C.GLclampf(b)/255.0,
		C.GLclampf(a)/255.0)
}

// ClearFlags is a bitset type for the flags to Clear.
type ClearFlags C.GLbitfield

const (
	// ColorBufferBit is a bit flag for Clear that specifies that the color buffer.
	ColorBufferBit ClearFlags = C.GL_COLOR_BUFFER_BIT

	// DepthBufferBit is a bit flag for Clear that specifies that the depth buffer.
	DepthBufferBit ClearFlags = C.GL_DEPTH_BUFFER_BIT
)

// Clear Clears the buffers specified by the bits.
func Clear(bits ClearFlags) {
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

// CheckError returns an error if on has occurred or else nil.
func CheckError() error {
	e := C.glGetError()
	if e == C.GL_NO_ERROR {
		return nil
	}
	return errors.New(errorStrings[e])
}

func makeError(p string) error {
	return errors.New(p + ": " + errorStrings[C.glGetError()])
}
