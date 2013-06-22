package gl

/*
#define GL_GLEXT_PROTOTYPES
#include <GL/gl.h>
#include <stdlib.h>

#cgo pkg-config: gl
*/
import "C"

import (
	"image"
	"unsafe"
)

// A Texture is data accessible via the graphics card.
type Texture struct {
	tex    C.GLuint
	target C.GLenum
}

// MakeImageTexture returns a texture from an image.
func MakeImageTexture(img *image.NRGBA) Texture {
	b := img.Bounds()
	w, h := C.GLsizei(b.Dx()), C.GLsizei(b.Dy())

	var t Texture
	t.target = C.GL_TEXTURE_2D
	C.glGenTextures(1, &t.tex)
	C.glBindTexture(C.GL_TEXTURE_2D, t.tex)
	defer C.glBindTexture(C.GL_TEXTURE_2D, 0)

	C.glTexImage2D(C.GL_TEXTURE_2D, 0, C.GL_RGBA, w, h, 0,
		C.GL_RGBA, C.GL_UNSIGNED_BYTE, unsafe.Pointer(&img.Pix[0]))

	C.glTexParameteri(C.GL_TEXTURE_2D, C.GL_TEXTURE_BASE_LEVEL, 0)
	C.glTexParameteri(C.GL_TEXTURE_2D, C.GL_TEXTURE_MAX_LEVEL, 0)
	C.glTexParameteri(C.GL_TEXTURE_2D, C.GL_TEXTURE_MAG_FILTER, C.GL_LINEAR)
	C.glTexParameteri(C.GL_TEXTURE_2D, C.GL_TEXTURE_MIN_FILTER, C.GL_LINEAR)
	return t
}

// Delete frees the resources allocated for the texture.
func (t Texture) Delete() {
	C.glDeleteTextures(1, &t.tex)
}

// Bind binds the texture to a particular texture image unit for its target.
func (t Texture) Bind(unit int) {
	C.glActiveTexture(C.GLenum(C.GL_TEXTURE0 + unit))
	C.glBindTexture(t.target, t.tex)
}
