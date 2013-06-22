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

// A Capability is a feature of OpenGL that can be enabled or disabled.
type Capability C.GLenum

const (
	// Texture2D is the 2D texture target.
	Texture2D Capability = C.GL_TEXTURE_2D

	// Blend is the alpha blending capability.
	Blend Capability = C.GL_BLEND
)

// Enable enables OpenGL capabilities.
func Enable(c Capability) {
	C.glEnable(C.GLenum(c))
}

// Disable disables OpenGL capabilities.
func Disable(c Capability) {
	C.glDisable(C.GLenum(c))
}

// BlendFactor specifies how either source or destination colors are blended.
type BlendFactor C.GLenum

const (
	// SrcAlpha uses the source color's alpha value.
	SrcAlpha BlendFactor = C.GL_SRC_ALPHA

	// OneMinusSrcAlpha uses 1 minus the source color's alpha value.
	OneMinusSrcAlpha BlendFactor = C.GL_ONE_MINUS_SRC_ALPHA
)

// BlendFunc sets the way that colors are blended.
func BlendFunc(srcFactor, dstFactor BlendFactor) {
	C.glBlendFunc(C.GLenum(srcFactor), C.GLenum(dstFactor))
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

// LineWidth sets the current line width.  If w non-positive then it is set to 1.
func LineWidth(w float32) {
	if w <= 0 {
		w = 1
	}
	C.glLineWidth(C.GLfloat(w))
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
