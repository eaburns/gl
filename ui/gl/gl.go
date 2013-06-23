package gl

/*
#define GL_GLEXT_PROTOTYPES
#include <GL/gl.h>

#cgo pkg-config: gl
*/
import "C"

import (
	"errors"
	"image/color"
)

// ClearFlags is a bitset type for the flags to Clear.
type ClearFlags C.GLbitfield

const (
	ColorBufferBit   ClearFlags = C.GL_COLOR_BUFFER_BIT
	DepthBufferBit   ClearFlags = C.GL_DEPTH_BUFFER_BIT
	StencilBufferBit ClearFlags = C.GL_STENCIL_BUFFER_BIT
)

// Clear clears buffers to preset values.
func Clear(bits ClearFlags) {
	C.glClear(C.GLbitfield(bits))
}

// ClearColor specifies clear values for the color buffers.
func ClearColor(col color.Color) {
	r, g, b, a := col.RGBA()
	C.glClearColor(
		C.GLclampf(r)/255.0,
		C.GLclampf(g)/255.0,
		C.GLclampf(b)/255.0,
		C.GLclampf(a)/255.0)
}

// ClearDepth specifies the clear value for the depth buffer.
func ClearDepth(d float32) {
	C.glClearDepthf(C.GLfloat(d))
}

// ClearStencil specifies the clear value for the stencil buffer.
func ClearStencil(s int) {
	C.glClearStencil(C.GLint(s))
}

// Flush forces execution of GL commands in finite time.
func Flush() {
	C.glFlush()
}

// Finish blocks until all GL execution is complete.
func Finish() {
	C.glFinish()
}

// A Capability is a feature of OpenGL that can be enabled or disabled.
type Capability C.GLenum

const (
	Blend  Capability = C.GL_BLEND
	Dither Capability = C.GL_DITHER
)

// Enable enables OpenGL capabilities.
func Enable(c Capability) {
	C.glEnable(C.GLenum(c))
}

// Disable disables OpenGL capabilities.
func Disable(c Capability) {
	C.glDisable(C.GLenum(c))
}

// BlendFunction specifies how either source or destination colors are blended.
type BlendFunction C.GLenum

const (
	Zero                  BlendFunction = C.GL_ZERO
	One                   BlendFunction = C.GL_ONE
	SrcColor              BlendFunction = C.GL_SRC_COLOR
	OneMinusSrcColor      BlendFunction = C.GL_ONE_MINUS_SRC_COLOR
	SrcAlpha              BlendFunction = C.GL_SRC_ALPHA
	OneMinusSrcAlpha      BlendFunction = C.GL_ONE_MINUS_SRC_ALPHA
	DstColor              BlendFunction = C.GL_DST_COLOR
	OneMinusDstColor      BlendFunction = C.GL_ONE_MINUS_DST_COLOR
	DstAlpha              BlendFunction = C.GL_DST_ALPHA
	OneMinusDstAlpha      BlendFunction = C.GL_ONE_MINUS_DST_ALPHA
	ConstantColor         BlendFunction = C.GL_CONSTANT_COLOR
	OneMinusConstantColor BlendFunction = C.GL_ONE_MINUS_CONSTANT_COLOR
	ConstantAlpha         BlendFunction = C.GL_CONSTANT_ALPHA
	OneMinusConstantAlpha BlendFunction = C.GL_ONE_MINUS_CONSTANT_ALPHA
)

// BlendFunc specifies pixel arithmetic
func BlendFunc(srcFactor, dstFactor BlendFunction) {
	C.glBlendFunc(C.GLenum(srcFactor), C.GLenum(dstFactor))
}

// LineWidth specifies the width of rasterized lines.
func LineWidth(w float32) {
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

// GetError returns error information: an error if one occurred or nil.
func GetError() error {
	e := C.glGetError()
	if e == C.GL_NO_ERROR {
		return nil
	}
	return errors.New(errorStrings[e])
}
