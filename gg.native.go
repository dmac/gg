// +build !js

package gg

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	mgl "github.com/go-gl/mathgl/mgl32"
)

var program uint32

func Init() error {
	var err error
	program, err = linkProgram(vertexShader, fragmentShader)
	if err != nil {
		return err
	}
	gl.UseProgram(program)

	//dims := make([]float32, 4)
	//gl.GetFloatv(gl.VIEWPORT, &dims[0]);

	//proj := mgl.Ortho2D(dims[0], dims[1], dims[2], dims[3])
	//projUniform := gl.GetUniformLocation(program, gl.Str("proj\x00"))
	//gl.UniformMatrix4fv(projUniform, 1, false, &proj[0])

	//view := viewMatrix()
	//viewUniform := gl.GetUniformLocation(program, gl.Str("view\x00"))
	//gl.UniformMatrix4fv(viewUniform, 1, false, &view[0])

	//fmt.Println(proj)
	//fmt.Println(view)

	return nil
}

type polyPlatformData struct {
	vao uint32
}

func (p *Poly) init() {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(p.vertices)*int(unsafe.Sizeof(Vec3{})),
		gl.Ptr(p.vertices),
		gl.STATIC_DRAW,
	)

	gl.GenVertexArrays(1, &p.vao)
	gl.BindVertexArray(p.vao)

	vattrib := uint32(gl.GetAttribLocation(program, gl.Str("vertex_position\x00")))
	gl.EnableVertexAttribArray(vattrib)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(vattrib, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
}

func (p *Poly) Draw() {
	//model := p.modelMatrix()
	//modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	//gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
	gl.UseProgram(program)
	gl.BindVertexArray(p.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}

func (p *Poly) modelMatrix() mgl.Mat3 {
	var scale float32 = 1
	S := mgl.Ident2().Mul(scale).Mat3()
	T := mgl.Translate2D(0, 0).Mul(scale)
	R := mgl.Rotate2D(0).Mat3()
	return T.Mul3(R).Mul3(S)
}

func linkProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vshader, err := compileShader(vertexShader, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	fshader, err := compileShader(fragmentShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	program := gl.CreateProgram()
	gl.AttachShader(program, vshader)
	gl.AttachShader(program, fshader)
	gl.LinkProgram(program)
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.TRUE {
		return program, nil
	}
	var logLength int32
	gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
	return 0, fmt.Errorf("link program: %s", log)
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	csource := gl.Str(string(append([]byte(source), 0)))
	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.TRUE {
		return shader, nil
	}
	var logLength int32
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
	return 0, fmt.Errorf("compile shader: %s%s", source, log)
}

func viewMatrix() mgl.Mat3 {
	R := mgl.Ident3()
	T := mgl.Ident3()
	return R.Mul3(T)
}
