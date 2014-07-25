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
}

func (c *CameraComponent) Name() string {
	return cameraComponentName
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

	fmt.Println("GL vendor:", gl.GetString(gl.VENDOR))
	fmt.Println("GL renderer:", gl.GetString(gl.RENDERER))
	fmt.Println("GL version:", gl.GetString(gl.VERSION))
	fmt.Println("GL shading language version:", gl.GetString(gl.SHADING_LANGUAGE_VERSION))
}

func (c *CameraComponent) PreUpdate() {
	gl.Viewport(0, 0, c.ScreenWidth, c.ScreenHeight)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (c *CameraComponent) PostUpdate() {
	sdl.GL_SwapBuffers()
}

func init() {
	ComponentNameMap[cameraComponentName] = func() Componenter { return &CameraComponent{} }
}

