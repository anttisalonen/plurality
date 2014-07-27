package plurality

import (
	"io/ioutil"
	"github.com/go-gl/gl"
)

var cameraComponentName string = "CameraComponent"

type CameraComponent struct {
	Component
	graphics *Graphics
}

func (c *CameraComponent) Name() string {
	return cameraComponentName
}

func (c *CameraComponent) InternalInit(game *GameApp) {
	c.graphics = &game.graphics
}

func loadShader(typ gl.GLenum, sourcefilename string) gl.Shader {
	source, err := ioutil.ReadFile(sourcefilename)
	if err != nil {
		panic(err)
	}

	var shader = gl.CreateShader(typ)
	shader.Source(string(source))
	shader.Compile()
	var compiled = shader.Get(gl.COMPILE_STATUS)
	if compiled == 0 {
		panic("Shader compilation: " + shader.GetInfoLog())
	}
	return shader
}

func initShader() gl.Program {
	var vs = loadShader(gl.VERTEX_SHADER, "../share/shader.vert")
	var fs = loadShader(gl.FRAGMENT_SHADER, "../share/shader.frag")
	var prog = gl.CreateProgram()
	prog.AttachShader(vs)
	prog.AttachShader(fs)
	prog.BindAttribLocation(0, "aPosition")
	prog.BindAttribLocation(1, "aTexcoord")
	prog.Link()
	var linked = prog.Get(gl.LINK_STATUS)
	if linked == 0 {
		panic("Shader linking: " + prog.GetInfoLog())
	}

	return prog
}

func (c *CameraComponent) PreUpdate() {
	gl.Viewport(0, 0, c.graphics.screenWidth, c.graphics.screenHeight)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	var uLoc = c.graphics.program.GetUniformLocation("uCamera")
	uLoc.Uniform2f(0.0, 0.0)

	var zoom float32 = 1.0

	uLoc = c.graphics.program.GetUniformLocation("uTop")
	uLoc.Uniform1f(1.0 * zoom)

	uLoc = c.graphics.program.GetUniformLocation("uRight")
	uLoc.Uniform1f(1.0 * zoom)

}

func init() {
	ComponentNameMap[cameraComponentName] = func() Componenter { return &CameraComponent{} }
}

