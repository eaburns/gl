package gl

/*
#define GL_GLEXT_PROTOTYPES
#include <GL/gl.h>
#include <stdlib.h>

#cgo pkg-config: gl
*/
import "C"

// A Buffer is an OpenGL buffer object.
type Buffer C.GLuint

// GenBuffers generates and returns n named buffer objects.
func GenBuffers(n int) []Buffer {
	bufs := make([]Buffer, n)
	C.glGenBuffers(C.GLsizei(n), (*C.GLuint)(&bufs[0]))
	return bufs
}

// DeleteBuffers deletes named buffer objects.
func DeleteBuffers(bufs []Buffer) {
	C.glDeleteBuffers(C.GLsizei(len(bufs)), (*C.GLuint)(&bufs[0]))
}

// Delete deletes the named buffer object.
func (b Buffer) Delete() {
	C.glDeleteBuffers(1, (*C.GLuint)(&b))
}

type BufferTarget C.GLenum

const (
	ArrayBuffer        BufferTarget = C.GL_ARRAY_BUFFER
	ElementArrayBuffer BufferTarget = C.GL_ELEMENT_ARRAY_BUFFER
)

// Bind binds a named buffer object.
func (b Buffer) Bind(targ BufferTarget) {
	C.glBindBuffer(C.GLenum(targ), C.GLuint(b))
}

type BufferDataUsage C.GLenum

const (
	StaticDraw  BufferDataUsage = C.GL_STATIC_DRAW
	StreamDraw  BufferDataUsage = C.GL_STREAM_DRAW
	DynamicDraw BufferDataUsage = C.GL_DYNAMIC_DRAW
)

// Data creates and initializes a buffer object's data store.
func BufferData(targ BufferTarget, data interface{}, usage BufferDataUsage) {
	_, sz, ptr := rawData(data)
	C.glBufferData(C.GLenum(targ), C.GLsizeiptr(sz), ptr, C.GLenum(usage))
}

// SubData updates a subset of a buffer object's data store.
func BufferSubData(targ BufferTarget, offs int, data interface{}) {
	_, sz, ptr := rawData(data)
	C.glBufferSubData(C.GLenum(targ), C.GLintptr(offs), C.GLsizeiptr(sz), ptr)
}
