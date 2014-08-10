package plurality

import (
	"fmt"
	"github.com/go-gl/gl"
)

var meshComponentName string = "MeshComponent"

type MeshComponent struct {
	Component
	vertexBuffer gl.Buffer
	texcoordBuffer gl.Buffer
	Scale Vector2
}

func (c *MeshComponent) Name() string {
	return meshComponentName
}

func init() {
	ComponentNameMap[meshComponentName] = func() Componenter { return &MeshComponent{} }
}

func (c *MeshComponent) SetVertices(vertexData []float32) {
	c.vertexBuffer.Delete()
	c.vertexBuffer = gl.GenBuffer()
	c.vertexBuffer.Bind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, len(vertexData) * 4, vertexData, gl.STATIC_DRAW)
}

func (c *MeshComponent) SetTextureCoordinates(texcoordData []float32) {
	c.texcoordBuffer.Delete()
	c.texcoordBuffer = gl.GenBuffer()
	c.texcoordBuffer.Bind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, len(texcoordData) * 4, texcoordData, gl.STATIC_DRAW)
}

func (c *MeshComponent) Update() {
	var vertexAttrib gl.AttribLocation = c.Graphics.program.GetAttribLocation("aPosition")
	vertexAttrib.EnableArray()
	c.vertexBuffer.Bind(gl.ARRAY_BUFFER)
	vertexAttrib.AttribPointer(3, gl.FLOAT, false, 0, nil)

	var texcoordAttrib gl.AttribLocation = c.Graphics.program.GetAttribLocation("aTexcoord")
	texcoordAttrib.EnableArray()
	c.texcoordBuffer.Bind(gl.ARRAY_BUFFER)
	texcoordAttrib.AttribPointer(4, gl.FLOAT, false, 0, nil)

	var uLoc = c.Graphics.program.GetUniformLocation("uTextured")
	uLoc.Uniform1i(0)

	uLoc = c.Graphics.program.GetUniformLocation("uPosition")
	var pos = c.GetTransform().Position
	uLoc.Uniform2f(float32(pos.X), float32(pos.Y))

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
	var err = gl.GetError()
	if err != 0 {
		fmt.Println("GL error", err)
	}
}


