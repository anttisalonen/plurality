package main

import (
	"fmt"
	"github.com/go-gl/gl"
)

var rectangleComponentName string = "RectangleComponent"

type RectangleComponent struct {
	Component
	graphics *Graphics
	vertexBuffer gl.Buffer
	texcoordBuffer gl.Buffer
	Position Vector2
}

func (c *RectangleComponent) Name() string {
	return rectangleComponentName
}

func (c *RectangleComponent) InternalInit(game *GameApp) {
	c.graphics = &game.graphics
}

func (c *RectangleComponent) Start() {
	c.vertexBuffer = gl.GenBuffer()
	c.vertexBuffer.Bind(gl.ARRAY_BUFFER)
	var vertexBufferData = []float32{0.0, 0.5, 0.0,   -0.5, -0.5, 0.0,   0.5, -0.5, 0.0}
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexBufferData) * 4, vertexBufferData, gl.STATIC_DRAW)

	c.texcoordBuffer = gl.GenBuffer()
	c.texcoordBuffer.Bind(gl.ARRAY_BUFFER)
	var texcoordBufferData = []float32{0.0, 0.0,   0.0, 1.0,   1.0, 1.0,  1.0, 0.0}
	gl.BufferData(gl.ARRAY_BUFFER, len(texcoordBufferData) * 4, texcoordBufferData, gl.STATIC_DRAW)
}

func (c *RectangleComponent) Update() {
	var vertexAttrib gl.AttribLocation = c.graphics.program.GetAttribLocation("aPosition")
	vertexAttrib.EnableArray()
	c.vertexBuffer.Bind(gl.ARRAY_BUFFER)
	vertexAttrib.AttribPointer(3, gl.FLOAT, false, 0, nil)

	var texcoordAttrib gl.AttribLocation = c.graphics.program.GetAttribLocation("aTexcoord")
	texcoordAttrib.EnableArray()
	c.texcoordBuffer.Bind(gl.ARRAY_BUFFER)
	texcoordAttrib.AttribPointer(4, gl.FLOAT, false, 0, nil)

	var uLoc = c.graphics.program.GetUniformLocation("uTextured")
	uLoc.Uniform1i(0)

	uLoc = c.graphics.program.GetUniformLocation("uPosition")
	uLoc.Uniform2f(float32(c.Position.X), float32(c.Position.Y))

	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	var err = gl.GetError()
	if err != 0 {
		fmt.Println("GL error", err)
	}
}

func init() {
	ComponentNameMap[rectangleComponentName] = func() Componenter { return &RectangleComponent{} }
}

