package main

import (
	"fmt"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/go-gl/gl"
)

var cameraComponentName string = "CameraComponent"

type CameraComponent struct {
	Component
	ScreenWidth int
	ScreenHeight int
	program gl.Program
	vertexBuffer gl.Buffer
}

func (c *CameraComponent) Name() string {
	return cameraComponentName
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

func (c *CameraComponent) initGL() {
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

func (c *CameraComponent) Start() {
	sdl.Init(sdl.INIT_EVERYTHING)
	screen := sdl.SetVideoMode(c.ScreenWidth, c.ScreenHeight, 32, sdl.OPENGL)

	if screen == nil {
		sdl.Quit()
		panic("SDL SetVideoMode: " + sdl.GetError() + "\n")
	}

	if gl.Init() != 0 {
		panic("GL init error")
	}

	sdl.WM_SetCaption("Plurality", "plurality")

	c.initGL()
	fmt.Println("GL vendor:", gl.GetString(gl.VENDOR))
	fmt.Println("GL renderer:", gl.GetString(gl.RENDERER))
	fmt.Println("GL version:", gl.GetString(gl.VERSION))
	fmt.Println("GL shading language version:", gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	c.vertexBuffer = gl.GenBuffer()
	c.vertexBuffer.Bind(gl.ARRAY_BUFFER)
	var vertexBufferData = []float32{0.0, 0.5, 0.0,   -0.5, -0.5, 0.0,   0.5, -0.5, 0.0}
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexBufferData) * 4, vertexBufferData, gl.STATIC_DRAW)
}

func (c *CameraComponent) Update() {
	gl.Viewport(0, 0, c.ScreenWidth, c.ScreenHeight)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	c.program.Use()
	var vertexAttrib gl.AttribLocation = 0
	vertexAttrib.EnableArray()
	c.vertexBuffer.Bind(gl.ARRAY_BUFFER)
	vertexAttrib.AttribPointer(3, gl.FLOAT, false, 0, nil)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	sdl.GL_SwapBuffers()
}

func init() {
	ComponentNameMap[cameraComponentName] = func() Componenter { return &CameraComponent{} }
}

