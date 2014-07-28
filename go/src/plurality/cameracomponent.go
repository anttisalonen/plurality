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

