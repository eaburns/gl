/*Package gl provides a set of bindings for OpenGL ES.

The intent of this package is to provide a very simple one-to-one set
of bindings, but with some small modification to be more Go-friendly,
and to take advantage of some of the niceties available in Go (extra
type safety, reflection, methods, color.Color, etc.). With very few
exceptions, each function or method in this package corresponds
directly to a single OpenGL ES function. The exceptions to this rule
are noted explicitly in their documentation.

A Note on Naming

Enum names are the same as their corresponding name in C, but they are
in CamlCase and have the "GL" prefix removed. For example:
GL_COLOR_BUFFER_BIT is ColorBufferBit. In a very small number of cases
the names needed an extra prefix to disambiguate constants of
different types that would otherwise have the same name. For example,
GL_DELETE_STATUS is both ShaderDeleteStatus or ProgramDeleteStatus.
One can be used when getting properties from shaders, and the other
when getting properties from programs. This provides better type
safety and is not expected to lead to any confusion. Additionally,
godoc should make it trivial to find the avaliable enum values for
differenc functions since they are associated with parameters via
their types.

Function names are the same as their corresponding name in C but with
the "gl" prefix and all size and type suffixes removed. Functions
whose C equivalents operate on objects (such as buffers, shaders,
programs, textures, etc.) have been converted to methods.

Method names are the same as their corresponding function names in C
but with the "gl" prefix, all size and type suffixes, and the object
name removed where possible. In cases where these removals results in
an empty identifier (e.g., glUniform) then the object name is not
removed. For example:
	glShaderCompile is Shader.Compile
	glGetProgram  is Program.Get
	glEnableVertexAttribArray is VertexAttribArray.Enable()
	glUniform is Uniform.Uniform
This should save typing in many places and be, overall, more pleasent.

It should be noted that these rules make for some non-idiomatic names.
For example, getter methods have a "Get" prefix (Shader.GetInfoLog)
and setter methods have no prefix (Shader.Source); Go prefers the
opposite: no prefix for getters, and a "Set" prefix for setters. While
annoying, this wasn't reconsiled on purpose. This package doesn't
provide detailed documentation, and the ability to easily go between
these names iand their C equivalents in order to find the
documentation is more important than idioms (even if Go's idiom is
superior ☺).

An apology

I must apologize in advance for the poor godoc comments. The comments
for functions and methods are straight out of the man pages and are,
in most cases, completely unhelpful. On the other hand, OpenGL is
fairly complex, and I am in no way qualified to write useful
documentation for it. I am also not qualified to write documentation
geared toward those who have never used OpenGL before, so this
package's documentation assumes a user who is familiar with OpenGL
from other sources. I found this tutorial helpful:
http://www.arcsynthesis.org/gltut.
*/
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
