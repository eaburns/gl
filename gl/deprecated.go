package gl

// Deprecated junk that we should work to ditch.

/*
#include <GL/gl.h>
#include <GL/glu.h>

#cgo pkg-config: gl glu
*/
import "C"

import (
	"image/color"
)

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
