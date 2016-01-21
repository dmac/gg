package gg

type Backend interface {
	Enable(Enum)
	DepthFunc(f Enum)
	BlendFunc(sfactor, dfactor Enum)
	Clear(mask Enum)
	ClearColor(r, g, b, a float32)
	CreateBuffer() *Buffer
	BindBuffer(typ Enum, b *Buffer)
	BufferData(typ Enum, src []byte, usage Enum)
	CreateShader(src []byte, typ Enum) (*Shader, error)
	AttachShader(*Program, *Shader)
	CreateProgram() *Program
	LinkProgram(*Program) error
	UseProgram(*Program)
	GetUniformLocation(*Program, string) (*Uniform, error)
	Uniform1f(*Uniform, float32)
	Uniform1i(*Uniform, int)
	Uniform4f(*Uniform, float32, float32, float32, float32)
	UniformMatrix4fv(*Uniform, []float32)
	GetAttribLocation(*Program, string) (*Attribute, error)
	EnableVertexAttribArray(*Attribute)
	VertexAttribPointer(
		a *Attribute, size int, typ Enum, normalized bool,
		stride, offset int,
	)
	CreateTexture() *Texture
	ActiveTexture(tex Enum)
	BindTexture(target Enum, texture *Texture)
	TexImage2D(
		target Enum, level int, internalFormat Enum,
		width, height, border int,
		format, typ Enum,
		data interface{},
	)
	TexParameteri(target Enum, pname Enum, param Enum)
	DrawArrays(mode Enum, first, count int)
}

type Buffer struct {
	Value interface{}
}

type Program struct {
	Value interface{}
}

type Shader struct {
	Value interface{}
}

type Uniform struct {
	Value interface{}
}

type Attribute struct {
	Value interface{}
}

type Texture struct {
	Value interface{}
}

type Enum uint32

var backend Backend

func Register(b Backend) {
	if b == nil {
		panic("gg: Register with nil backend")
	}
	if backend != nil {
		panic("gg: Register called twice")
	}
	backend = b
}

func Enable(c Enum) {
	backend.Enable(c)
}

func DepthFunc(f Enum) {
	backend.DepthFunc(f)
}

func BlendFunc(sfactor, dfactor Enum) {
	backend.BlendFunc(sfactor, dfactor)
}

func Clear(mask Enum) {
	backend.Clear(mask)
}

func ClearColor(r, g, b, a float32) {
	backend.ClearColor(r, g, b, a)
}

func CreateBuffer() *Buffer {
	return backend.CreateBuffer()
}

func BindBuffer(typ Enum, b *Buffer) {
	backend.BindBuffer(typ, b)
}

func BufferData(typ Enum, src []byte, usage Enum) {
	backend.BufferData(typ, src, usage)
}

func CreateShader(src []byte, typ Enum) (*Shader, error) {
	return backend.CreateShader(src, typ)
}

func CreateProgram() *Program {
	return backend.CreateProgram()
}

func AttachShader(p *Program, s *Shader) {
	backend.AttachShader(p, s)
}

func LinkProgram(p *Program) error {
	return backend.LinkProgram(p)
}

func UseProgram(p *Program) {
	backend.UseProgram(p)
}

func GetUniformLocation(p *Program, name string) (*Uniform, error) {
	return backend.GetUniformLocation(p, name)
}

func Uniform1f(u *Uniform, v0 float32) {
	backend.Uniform1f(u, v0)
}

func Uniform4f(u *Uniform, v0, v1, v2, v3 float32) {
	backend.Uniform4f(u, v0, v1, v2, v3)
}

func Uniform1i(u *Uniform, v0 int) {
	backend.Uniform1i(u, v0)
}

func UniformMatrix4fv(u *Uniform, value []float32) {
	backend.UniformMatrix4fv(u, value)
}

func GetAttribLocation(p *Program, name string) (*Attribute, error) {
	return backend.GetAttribLocation(p, name)
}

func EnableVertexAttribArray(a *Attribute) {
	backend.EnableVertexAttribArray(a)
}

func VertexAttribPointer(a *Attribute, size int, typ Enum, normalized bool, stride, offset int) {
	backend.VertexAttribPointer(a, size, typ, normalized, stride, offset)
}

func CreateTexture() *Texture {
	return backend.CreateTexture()
}

func ActiveTexture(tex Enum) {
	backend.ActiveTexture(tex)
}

func BindTexture(target Enum, texture *Texture) {
	backend.BindTexture(target, texture)
}

func TexImage2D(
	target Enum, level int, internalFormat Enum,
	width, height, border int,
	format, typ Enum,
	data interface{},
) {
	backend.TexImage2D(
		target, level, internalFormat,
		width, height, border,
		format, typ,
		data,
	)
}

func TexParameteri(target Enum, pname Enum, param Enum) {
	backend.TexParameteri(target, pname, param)
}

func DrawArrays(mode Enum, first, count int) {
	backend.DrawArrays(mode, first, count)
}
