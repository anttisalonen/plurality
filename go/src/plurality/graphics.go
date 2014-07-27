package plurality

import (
	"fmt"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/go-gl/gl"
)

type Graphics struct {
	program gl.Program
	screenWidth int
	screenHeight int
}

func (c *Graphics) Init(width int, height int) {
	sdl.Init(sdl.INIT_EVERYTHING)
	c.screenWidth = width
	c.screenHeight = height
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
	c.program = initShader()
	c.program.Use()
}

func (c *Graphics) Update() bool {
	for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
		switch e := ev.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			if e.Keysym.Sym == 27 { // escape
				return false
			}
		}
	}

	sdl.GL_SwapBuffers()
	return true
}
