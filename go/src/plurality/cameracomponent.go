package plurality

import (
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

func (c *CameraComponent) PreUpdate() {
	gl.Viewport(0, 0, c.graphics.ScreenWidth, c.graphics.ScreenHeight)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	var uLoc = c.graphics.program.GetUniformLocation("uCamera")
	var pos = c.GetTransform().Position
	uLoc.Uniform2f(float32(pos.X), float32(pos.Y))

	var zoom float32 = 1.0

	uLoc = c.graphics.program.GetUniformLocation("uRight")
	uLoc.Uniform1f(float32(c.graphics.ScreenWidth) * zoom * 0.5)

	uLoc = c.graphics.program.GetUniformLocation("uTop")
	uLoc.Uniform1f(float32(c.graphics.ScreenHeight) * zoom * 0.5)

}

func init() {
	ComponentNameMap[cameraComponentName] = func() Componenter { return &CameraComponent{} }
}

