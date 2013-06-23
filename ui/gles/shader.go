package gles

/*
#define GL_GLEXT_PROTOTYPES
#include <GL/gl.h>
#include <stdlib.h>

#cgo pkg-config: gl
*/
import "C"

import (
	"unsafe"
)

// A ShaderKind specifies the type of shader.
type ShaderKind C.GLenum

const (
	VertexShader   ShaderKind = C.GL_VERTEX_SHADER
	FragmentShader ShaderKind = C.GL_FRAGMENT_SHADER
)

type Shader C.GLuint

// CreateShader creates a shader object.
func CreateShader(kind ShaderKind) Shader {
	sh := C.glCreateShader(C.GLenum(kind))
	return Shader(sh)
}

// Delete deletes a shader.
func (s Shader) Delete() {
	C.glDeleteShader(C.GLuint(s))
}

// Source replaces the source code in a shader object.
func (s Shader) Source(src ...string) {
	if len(src) == 0 {
		return
	}
	csrc := make([]*C.GLchar, len(src))
	for i, s := range src {
		csrc[i] = (*C.GLchar)(C.CString(s))
		defer C.free(unsafe.Pointer(csrc[i]))
	}
	C.glShaderSource(C.GLuint(s), C.GLsizei(len(src)), &csrc[0], nil)
}

// Compile compiles a shader object.
func (s Shader) Compile() {
	C.glCompileShader(C.GLuint(s))
}

// A ShaderParameter names a gettable parameter of a shader.
type ShaderParameter C.GLenum

const (
	CompileStatus       ShaderParameter = C.GL_COMPILE_STATUS
	ShaderInfoLogLength ShaderParameter = C.GL_INFO_LOG_LENGTH
	ShaderSourceLength  ShaderParameter = C.GL_SHADER_SOURCE_LENGTH
	ShaderType          ShaderParameter = C.GL_SHADER_TYPE
	ShaderDeleteStatus  ShaderParameter = C.GL_DELETE_STATUS
)

// Get returns a parameter from a shader object
func (s Shader) Get(parm ShaderParameter) int {
	var vl C.GLint
	C.glGetShaderiv(C.GLuint(s), C.GLenum(parm), &vl)
	return int(vl)
}

// GetInfoLog returns the information log for a shader object.
func (s Shader) GetInfoLog() string {
	sz := s.Get(ShaderInfoLogLength)
	cstr := (*C.char)(C.malloc(C.size_t(sz + 1)))
	defer C.free(unsafe.Pointer(cstr))
	C.glGetShaderInfoLog(C.GLuint(s), C.GLsizei(sz), nil, (*C.GLchar)(cstr))
	return C.GoString(cstr)
}

// A Program is a set of linked shaders that can be loaded onto the graphics card.
type Program C.GLuint

// CreateProgram creates a program object.
func CreateProgram() Program {
	return Program(C.glCreateProgram())
}

// Delete deletes a program object.
func (p Program) Delete() {
	C.glDeleteProgram(C.GLuint(p))
}

// AttachShader attaches a shader object to a program object.
func (p Program) AttachShader(s Shader) {
	C.glAttachShader(C.GLuint(p), C.GLuint(s))
}

// DetachShader detaches a shader object from a program object to which it is attached.
func (p Program) DetachShader(s Shader) {
	C.glDetachShader(C.GLuint(p), C.GLuint(s))
}

// Link links a program object.
func (p Program) Link() {
	C.glLinkProgram(C.GLuint(p))
}

// Use installs a program object as part of current rendering state.
func (p Program) Use() {
	C.glUseProgram(C.GLuint(p))
}

// A ProgramParameter is a gettable parameter of a program.
type ProgramParameter C.GLenum

const (
	ProgramDeleteStatus      ProgramParameter = C.GL_DELETE_STATUS
	LinkStatus               ProgramParameter = C.GL_LINK_STATUS
	ValidateStatus           ProgramParameter = C.GL_VALIDATE_STATUS
	ProgramInfoLogLength     ProgramParameter = C.GL_INFO_LOG_LENGTH
	AttachedShaders          ProgramParameter = C.GL_ATTACHED_SHADERS
	ActiveAttributes         ProgramParameter = C.GL_ACTIVE_ATTRIBUTES
	ActiveAttributeMaxLength ProgramParameter = C.GL_ACTIVE_ATTRIBUTE_MAX_LENGTH
	ActiveUniforms           ProgramParameter = C.GL_ACTIVE_UNIFORMS
	ActiveUniformMaxLength   ProgramParameter = C.GL_ACTIVE_UNIFORM_MAX_LENGTH
)

// Get returns a parameter from a program object.
func (p Program) Get(parm ProgramParameter) int {
	var vl C.GLint
	C.glGetProgramiv(C.GLuint(p), C.GLenum(parm), &vl)
	return int(vl)
}

// GetInfoLog returns the information log for a program object.
func (p Program) GetInfoLog() string {
	sz := p.Get(ProgramInfoLogLength)
	cstr := (*C.char)(C.malloc(C.size_t(sz + 1)))
	defer C.free(unsafe.Pointer(cstr))
	C.glGetProgramInfoLog(C.GLuint(p), C.GLsizei(sz), nil, (*C.GLchar)(cstr))
	return C.GoString(cstr)
}

type UniformLocation C.GLint

// GetUniformLocation returns the location of a uniform variable.
func (p Program) GetUniformLocation(name string) UniformLocation {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	l := C.glGetUniformLocation(C.GLuint(p), (*C.GLchar)(cstr))
	return UniformLocation(l)
}

// GetAttribLocation returns the location of an attribute variable.
func (p Program) GetAttribLocation(name string) AttributeLocation {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	l := C.glGetAttribLocation(C.GLuint(p), (*C.GLchar)(cstr))
	return AttributeLocation(l)
}

// Uniform specifies the value of a uniform variable for the current program object.
func Uniform(l UniformLocation, vls ...interface{}) {
	if len(vls) == 0 || len(vls) > 4 {
		panic("Uniform requires 1, 2, 3, or 4 values")
	}
	switch vls[0].(type) {
	case float32:
		switch len(vls) {
		case 1:
			C.glUniform1f(C.GLint(l), C.GLfloat(vls[0].(float32)))
		case 2:
			C.glUniform2f(C.GLint(l), C.GLfloat(vls[0].(float32)), C.GLfloat(vls[1].(float32)))
		case 3:
			C.glUniform3f(C.GLint(l), C.GLfloat(vls[0].(float32)), C.GLfloat(vls[1].(float32)), C.GLfloat(vls[2].(float32)))
		case 4:
			C.glUniform4f(C.GLint(l), C.GLfloat(vls[0].(float32)), C.GLfloat(vls[1].(float32)), C.GLfloat(vls[2].(float32)), C.GLfloat(vls[3].(float32)))
		}
	case int:
		switch len(vls) {
		case 1:
			C.glUniform1i(C.GLint(l), C.GLint(vls[0].(int)))
		case 2:
			C.glUniform2i(C.GLint(l), C.GLint(vls[0].(int)), C.GLint(vls[1].(int)))
		case 3:
			C.glUniform3i(C.GLint(l), C.GLint(vls[0].(int)), C.GLint(vls[1].(int)), C.GLint(vls[2].(int)))
		case 4:
			C.glUniform4i(C.GLint(l), C.GLint(vls[0].(int)), C.GLint(vls[1].(int)), C.GLint(vls[2].(int)), C.GLint(vls[3].(int)))
		}
	default:
		panic("Uniform only accepts int and float32 typed parameters")
	}
}
