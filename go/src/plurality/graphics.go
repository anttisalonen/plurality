package plurality

import (
	"fmt"
	"io/ioutil"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/go-gl/gl"
)

type Graphics struct {
	program gl.Program
	ScreenWidth int
	ScreenHeight int
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

func (c *Graphics) Init(width int, height int) {
	sdl.Init(sdl.INIT_EVERYTHING)
	c.ScreenWidth = width
	c.ScreenHeight = height
	screen := sdl.SetVideoMode(width, height, 32, sdl.OPENGL)

	if screen == nil {
		sdl.Quit()
		panic("SDL SetVideoMode: " + sdl.GetError() + "\n")
	}

	if gl.Init() != 0 {
		panic("GL init error")
	}

	sdl.WM_SetCaption("Plurality", "plurality")

	fmt.Println("GL vendor:", gl.GetString(gl.VENDOR))
	fmt.Println("GL renderer:", gl.GetString(gl.RENDERER))
	fmt.Println("GL version:", gl.GetString(gl.VERSION))
	fmt.Println("GL shading language version:", gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	gl.ClearColor(0, 0, 0, 0)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.TEXTURE_2D)
	c.program = initShader()
	c.program.Use()
}

func (c *Graphics) Update() {
	sdl.GL_SwapBuffers()
}
