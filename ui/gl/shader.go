package gl

/*
#define GL_GLEXT_PROTOTYPES
#include <GL/gl.h>
#include <stdlib.h>

#cgo pkg-config: gl
*/
import "C"

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"unsafe"
)

// A Program is a set of linked shaders.
type Program struct {
	prog     C.GLuint
	attrLocs map[string]C.GLint
	unifLocs map[string]C.GLint
}

// NewProgram complies and links a vertex and fragment shader and returns a new *Program.
func NewProgram(vert, frag io.Reader) (*Program, error) {
	v, err := compileShader(C.GL_VERTEX_SHADER, vert)
	if err != nil {
		return nil, err
	}
	defer C.glDeleteShader(v)

	f, err := compileShader(C.GL_FRAGMENT_SHADER, frag)
	if err != nil {
		return nil, err
	}
	defer C.glDeleteShader(f)

	p, err := linkProgram(v, f)
	if err != nil {
		return nil, err
	}

	return &Program{
		prog:     p,
		attrLocs: make(map[string]C.GLint, 2),
		unifLocs: make(map[string]C.GLint, 2),
	}, nil
}

func linkProgram(vert, frag C.GLuint) (C.GLuint, error) {
	p := C.glCreateProgram()
	if p == 0 {
		return 0, makeError("Failed to create program")
	}

	C.glAttachShader(p, vert)
	C.glAttachShader(p, frag)
	C.glLinkProgram(p)

	var ok C.GLint
	C.glGetProgramiv(p, C.GL_LINK_STATUS, &ok)
	if ok == 0 {
		err := errors.New("Failed to link program: " + linkErrMsg(p))
		C.glDeleteProgram(p)
		return 0, err
	}
	return p, nil
}

func linkErrMsg(p C.GLuint) string {
	var msgSize C.GLint
	C.glGetProgramiv(p, C.GL_INFO_LOG_LENGTH, &msgSize)

	errMsg := "<no message>"
	if msgSize > 0 {
		buf := glbuf(msgSize)
		defer C.free(unsafe.Pointer(buf))
		C.glGetProgramInfoLog(p, C.GLsizei(msgSize), nil, buf)
		errMsg = C.GoString((*C.char)(buf))
	}
	return errMsg
}

var shaderKindNames = map[C.GLenum]string{
	C.GL_VERTEX_SHADER:   "GL_VERTEX_SHADER",
	C.GL_FRAGMENT_SHADER: "GL_FRAGMENT_SHADER",
}

func compileShader(kind C.GLenum, r io.Reader) (C.GLuint, error) {
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}

	kindName := shaderKindNames[kind]
	if kindName == "" {
		kindName = "unknown shader type"
	}

	sh := C.glCreateShader(kind)
	if sh == 0 {
		return 0, makeError("Failed to create " + kindName)
	}

	csrc := glstring(string(src))
	defer C.free(unsafe.Pointer(csrc))
	C.glShaderSource(sh, 1, &csrc, nil)

	C.glCompileShader(sh)

	var ok C.GLint
	C.glGetShaderiv(sh, C.GL_COMPILE_STATUS, &ok)
	if ok == 0 {
		err = errors.New("Failed to compile " + kindName + ": " + compileErrMsg(sh))
		C.glDeleteShader(sh)
		return 0, err
	}
	return sh, nil
}

func compileErrMsg(sh C.GLuint) string {
	var msgSize C.GLint
	C.glGetShaderiv(sh, C.GL_INFO_LOG_LENGTH, &msgSize)

	errMsg := "<no message>"
	if msgSize > 0 {
		buf := glbuf(msgSize)
		defer C.free(unsafe.Pointer(buf))
		C.glGetShaderInfoLog(sh, C.GLsizei(msgSize), nil, buf)
		errMsg = C.GoString((*C.char)(buf))
	}
	return errMsg
}

// Delete deletes the program, freeing its resources.
func (p *Program) Delete() {
	C.glDeleteProgram(p.prog)
}

// SetUniform sets the value(s) for a uniform.
// The variadic parameter must contain 1, 2, 3, or 4 values.
func (p *Program) SetUniform(name string, vls ...interface{}) error {
	C.glUseProgram(p.prog)
	defer C.glUseProgram(0)
	l, err := p.uniformLocation(name)
	if err != nil {
		return err
	}

	if len(vls) == 0 || len(vls) > 4 {
		panic("Uniform requires 1, 2, 3, or 4 values")
	}

	switch vls[0].(type) {
	case float32:
		switch len(vls) {
		case 1:
			C.glUniform1f(l, C.GLfloat(vls[0].(float32)))
		case 2:
			C.glUniform2f(l,
				C.GLfloat(vls[0].(float32)),
				C.GLfloat(vls[1].(float32)),
			)
		case 3:
			C.glUniform3f(l,
				C.GLfloat(vls[0].(float32)),
				C.GLfloat(vls[1].(float32)),
				C.GLfloat(vls[2].(float32)),
			)
		case 4:
			C.glUniform4f(l,
				C.GLfloat(vls[0].(float32)),
				C.GLfloat(vls[1].(float32)),
				C.GLfloat(vls[2].(float32)),
				C.GLfloat(vls[3].(float32)),
			)
		}
	case int:
		switch len(vls) {
		case 1:
			C.glUniform1i(l, C.GLint(vls[0].(int)))
		case 2:
			C.glUniform2i(l,
				C.GLint(vls[0].(int)),
				C.GLint(vls[1].(int)),
			)
		case 3:
			C.glUniform3i(l,
				C.GLint(vls[0].(int)),
				C.GLint(vls[1].(int)),
				C.GLint(vls[2].(int)),
			)
		case 4:
			C.glUniform4i(l,
				C.GLint(vls[0].(int)),
				C.GLint(vls[1].(int)),
				C.GLint(vls[2].(int)),
				C.GLint(vls[3].(int)),
			)
		}
	default:
		panic(fmt.Sprintf("Unknown type in SetUniform: %T", vls[0]))
	}
	return nil
}

func (p *Program) uniformLocation(name string) (C.GLint, error) {
	l, ok := p.unifLocs[name]
	if ok {
		return l, nil
	}

	cstr := glstring(name)
	defer C.free(unsafe.Pointer(cstr))
	if l = C.glGetUniformLocation(p.prog, cstr); l < 0 {
		return -1, errors.New("Failed to get uniform location for " + name)
	}

	p.unifLocs[name] = l
	return l, nil
}

func uniformf(l C.GLint, f ...float32) {
	switch len(f) {
	case 1:
		C.glUniform1f(l, C.GLfloat(f[0]))
	case 2:
		C.glUniform2f(l, C.GLfloat(f[0]), C.GLfloat(f[1]))
	case 3:
		C.glUniform3f(l, C.GLfloat(f[0]), C.GLfloat(f[1]), C.GLfloat(f[2]))
	case 4:
		C.glUniform4f(l, C.GLfloat(f[0]), C.GLfloat(f[1]), C.GLfloat(f[2]), C.GLfloat(f[3]))
	default:
		panic("Uniform requires 1, 2, 3, or 4 values")
	}
}

// SetVertexAttributeData sets vertex data for the named attribute to that of the
// buffer currently bound to GL_ARRAY_BUFFER.
func (p *Program) SetVertexAttributeData(name string, size, stride, offs int) error {
	C.glUseProgram(p.prog)
	defer C.glUseProgram(0)

	l, err := p.attributeLocation(name)
	if err != nil {
		return err
	}

	C.glVertexAttribPointer(
		C.GLuint(l),
		C.GLint(size),
		C.GL_FLOAT,
		C.GL_FALSE,
		C.GLsizei(stride),
		unsafe.Pointer(uintptr(offs)))
	return nil
}

func (p *Program) attributeLocation(name string) (C.GLint, error) {
	l, ok := p.attrLocs[name]
	if ok {
		return l, nil
	}

	cstr := glstring(name)
	defer C.free(unsafe.Pointer(cstr))
	if l = C.glGetAttribLocation(p.prog, cstr); l < 0 {
		return -1, errors.New("Failed to get attribute location for " + name)
	}

	p.attrLocs[name] = l
	return l, nil
}

// A DrawMode specifies how vertices should be interpreted when drawing.
type DrawMode C.GLenum

const (
	// LineStrip mode connects the vertices with a line.
	LineStrip DrawMode = C.GL_LINE_STRIP

	// Triangles mode draws triangles between each set of three vertices, no vertex is reused.
	Triangles DrawMode = C.GL_TRIANGLES

	// TriangleStrip mode draws triangles between each set of three vertices.  If multiple
	// triangles are drawn, then the final two vertices of the previous triangle are reused
	// as the first two vertices of the next.
	TriangleStrip DrawMode = C.GL_TRIANGLE_STRIP
)

// DrawArrays draws using this program an all vertex attribute arrays that that have been set.
func (p *Program) DrawArrays(mode DrawMode, first, count int) {
	C.glUseProgram(p.prog)
	defer C.glUseProgram(0)

	alocs := make([]C.GLuint, len(p.attrLocs))
	i := 0
	for _, l := range p.attrLocs {
		alocs[i] = C.GLuint(l)
		i++
	}

	for _, l := range alocs {
		C.glEnableVertexAttribArray(l)
	}

	C.glDrawArrays(C.GLenum(mode), C.GLint(first), C.GLsizei(count))

	for _, l := range alocs {
		C.glDisableVertexAttribArray(l)
	}
}

// Glstring copies a string into a C string and returns it cast into a *C.GLchar.
// The return value must be freed by the caller using C.free.
func glstring(s string) *C.GLchar {
	return (*C.GLchar)(unsafe.Pointer(C.CString(s)))
}

// Glbuf allocates a buffer of the given number of bytes and returns it cast into a *C.GLchar.
// The return value must be freed by the caller using C.free.
func glbuf(sz C.GLint) *C.GLchar {
	return (*C.GLchar)(C.malloc(C.size_t(sz)))
}

// An Buffer is a named buffer object.
type Buffer struct {
	buf      C.GLuint
	capacity int // in bytes
	target   C.GLenum
	usage    DataUsage
}

// NewArrayBuffer returns new buffer object using the GL_ARRAY_BUFFER target.
func NewArrayBuffer() *Buffer {
	b := Buffer{target: C.GL_ARRAY_BUFFER}
	C.glGenBuffers(1, &b.buf)
	return &b
}

// Delete deletes the buffer.
func (b *Buffer) Delete() {
	C.glDeleteBuffers(1, &b.buf)
}

// Bind binds the buffer.
func (b *Buffer) Bind() {
	C.glBindBuffer(b.target, b.buf)
}

// DataUsage is a hint specifying how buffer data will be used.
type DataUsage C.GLenum

const (
	// StaticDraw hints that:
	// "The data store contents will be modified once and used many times."
	// and "The data store contents are modified by the application, and used as the
	// source for GL drawing and image specification commands."
	StaticDraw DataUsage = C.GL_STATIC_DRAW

	// StreamDraw hints that:
	// "The data store contents will be modified once and used at most a few times."
	// and "The data store contents are modified by the application, and used as the
	// source for GL drawing and image specification commands."
	StreamDraw DataUsage = C.GL_STREAM_DRAW

	// DynamicDraw hints that:
	// "The data store contents will be modified repeatedly and used many times."
	// and "The data store contents are modified by the application, and used as the
	// source for GL drawing and image specification commands."
	DynamicDraw DataUsage = C.GL_DYNAMIC_DRAW
)

// SetData sets the buffer's data to the given floats.
// If the buffer has already been allocated for the same usage, with sufficient
// capacity then it is reused, otherwise it is reallocated.
func (b *Buffer) SetData(usage DataUsage, fs ...float32) {
	b.Bind()

	data := unsafe.Pointer(&fs[0])
	sz := len(fs) * 4
	if sz <= b.capacity && usage == b.usage {
		C.glBufferSubData(b.target, 0, C.GLsizeiptr(sz), data)
	}
	b.capacity = sz
	b.usage = usage
	C.glBufferData(b.target, C.GLsizeiptr(sz), data, C.GLenum(b.usage))
}
