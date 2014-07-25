package main

import (
	"fmt"
	"github.com/go-gl/gl"
)

var rectangleComponentName string = "RectangleComponent"

type RectangleComponent struct {
	Component
	program gl.Program
	vertexBuffer gl.Buffer
}

func (c *RectangleComponent) Name() string {
	return rectangleComponentName
}

var vertShader string =
`attribute vec4 vPosition;
void main() {
    gl_Position = vPosition;
}`

var fragShader string =
`#ifdef GL_ES
precision mediump float;
#endif

void main() {
    gl_FragColor = vec4(1.0, 0.0, 0.0, 1.0);
}`

func loadShader(typ gl.GLenum, source string) gl.Shader {
	var shader = gl.CreateShader(typ)
	shader.Source(source)
	shader.Compile()
	var compiled = shader.Get(gl.COMPILE_STATUS)
	if compiled == 0 {
		panic("Shader compilation: " + shader.GetInfoLog())
	}
	return shader
}

func (c *RectangleComponent) initShader() {
	var vs = loadShader(gl.VERTEX_SHADER, vertShader)
	var fs = loadShader(gl.FRAGMENT_SHADER, fragShader)
	var prog = gl.CreateProgram()
	prog.AttachShader(vs)
	prog.AttachShader(fs)
	prog.BindAttribLocation(0, "vPosition")
	prog.Link()
	var linked = prog.Get(gl.LINK_STATUS)
	if linked == 0 {
		panic("Shader linking: " + prog.GetInfoLog())
	}

	gl.ClearColor(0, 0, 0, 0)
	c.program = prog
}

func (c *RectangleComponent) Start() {
	c.initShader()
	c.vertexBuffer = gl.GenBuffer()
	c.vertexBuffer.Bind(gl.ARRAY_BUFFER)
	var vertexBufferData = []float32{0.0, 0.5, 0.0,   -0.5, -0.5, 0.0,   0.5, -0.5, 0.0}
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexBufferData) * 4, vertexBufferData, gl.STATIC_DRAW)
}

func (c *RectangleComponent) Update() {
	c.program.Use()
	var vertexAttrib gl.AttribLocation = 0
	vertexAttrib.EnableArray()
	c.vertexBuffer.Bind(gl.ARRAY_BUFFER)
	vertexAttrib.AttribPointer(3, gl.FLOAT, false, 0, nil)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	var err = gl.GetError()
	if err != 0 {
		fmt.Println("GL error", err)
	}
}

func init() {
	ComponentNameMap[rectangleComponentName] = func() Componenter { return &RectangleComponent{} }
}

