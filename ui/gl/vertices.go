package gl

/*
#define GL_GLEXT_PROTOTYPES
#include <GL/gl.h>
#include <stdlib.h>

#cgo pkg-config: gl
*/
import "C"

import (
	"reflect"
	"unsafe"
)

// A DataType defines the type of data elements.
type DataType C.GLenum

const (
	Float         DataType = C.GL_FLOAT
	Byte          DataType = C.GL_BYTE
	Short         DataType = C.GL_SHORT
	UnsignedByte  DataType = C.GL_UNSIGNED_BYTE
	UnsignedShort DataType = C.GL_UNSIGNED_SHORT
)

var typeEnums = map[reflect.Kind]DataType{
	reflect.Float32: Float,
	reflect.Int8:    Byte,
	reflect.Int16:   Short,
	reflect.Uint8:   UnsignedByte,
	reflect.Uint16:  UnsignedShort,
}

func rawData(i interface{}) (tipe DataType, sz int, ptr unsafe.Pointer) {
	v := reflect.ValueOf(i)
	if k := v.Kind(); k != reflect.Array && k != reflect.Slice {
		panic("Data must be an array or a slice")
	}

	tipe, ok := typeEnums[v.Type().Elem().Kind()]
	if !ok {
		panic("Invalid data element type: " + v.Type().Elem().Kind().String())
	}

	nElms := v.Len()
	if nElms > 0 {
		ptr = unsafe.Pointer(v.Index(0).UnsafeAddr())
	}

	elmSz := v.Type().Elem().Size()
	return tipe, nElms * int(elmSz), ptr
}

type VertexAttribArray C.GLuint

// Pointer defines an array of generic vertex attribute data.
func (l VertexAttribArray) Pointer(sz int, t DataType, norm bool, stride, offs int) {
	n := 0
	if norm {
		n = 1
	}
	C.glVertexAttribPointer(C.GLuint(l), C.GLint(sz), C.GLenum(t), C.GLboolean(n), C.GLsizei(stride), unsafe.Pointer(uintptr(offs)))
}

// Enable enables a generic vertex attribute array.
func (l VertexAttribArray) Enable() {
	C.glEnableVertexAttribArray(C.GLuint(l))
}

// Disable disables a generic vertex attribute array.
func (l VertexAttribArray) Disable() {
	C.glDisableVertexAttribArray(C.GLuint(l))
}

// DrawMode specifies what to draw.
type DrawMode C.GLenum

const (
	Points        DrawMode = C.GL_POINTS
	LineStrip     DrawMode = C.GL_LINE_STRIP
	LineLoop      DrawMode = C.GL_LINE_LOOP
	Lines         DrawMode = C.GL_LINES
	TriangleStrip DrawMode = C.GL_TRIANGLE_STRIP
	TriangleFan   DrawMode = C.GL_TRIANGLE_FAN
	Triangles     DrawMode = C.GL_TRIANGLES
)

// DrawArrays renders primitives from array data.
func DrawArrays(mode DrawMode, first, count int) {
	C.glDrawArrays(C.GLenum(mode), C.GLint(first), C.GLsizei(count))
}
