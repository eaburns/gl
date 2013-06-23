package gles

/*
#define GL_GLEXT_PROTOTYPES
#include <GL/gl.h>
#include <stdlib.h>

#cgo pkg-config: gl
*/
import "C"

// A Texture is the name of a texture object.
type Texture C.GLuint

// GenTextures generates texture names.
func GenTextures(n int) []Texture {
	texs := make([]Texture, n)
	C.glGenTextures(C.GLsizei(n), (*C.GLuint)(&texs[0]))
	return texs
}

// DeleteTextures deletes named textures.
func DeleteTextures(texs []Texture) {
	C.glDeleteTextures(C.GLsizei(len(texs)), (*C.GLuint)(&texs[0]))
}

// Delete deletes a named texture.
func (tex Texture) Delete() {
	C.glDeleteTextures(1, (*C.GLuint)(&tex))
}

// A TextureTarget is a texturing target to which a texture can be bound.
type TextureTarget C.GLenum

const (
	Texture2D TextureTarget = C.GL_TEXTURE_2D
)

// Bind binds a named texture to a texturing target.
func (tex Texture) Bind(targ TextureTarget) {
	C.glBindTexture(C.GLenum(targ), C.GLuint(tex))
}

// ActiveTexture selects active texture unit.
// Unit must be between 0 and the maximum supported texture units, of which there are at least 80.
func ActiveTexture(unit int) {
	C.glActiveTexture(C.GLenum(unit))
}

// A TextureFormat specifies how to interpret texture data.
type TextureFormat C.GLenum

const (
	Alpha          TextureFormat = C.GL_ALPHA
	Luminance      TextureFormat = C.GL_LUMINANCE
	LuminanceAlpha TextureFormat = C.GL_LUMINANCE_ALPHA
	RGB            TextureFormat = C.GL_RGB
	RGBA           TextureFormat = C.GL_RGBA
)

// TexImage2D specifies a two-dimensional texture image.
func TexImage2D(targ TextureTarget, lvl int, ifmt TextureFormat, w, h, border int, fmt TextureFormat, data interface{}) {
	tipe, _, ptr := rawData(data)
	C.glTexImage2D(C.GLenum(targ),
		C.GLint(lvl),
		C.GLint(ifmt),
		C.GLsizei(w),
		C.GLsizei(h),
		C.GLint(border),
		C.GLenum(fmt),
		C.GLenum(tipe),
		ptr)
}
