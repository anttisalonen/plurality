package plurality

import (
	"fmt"
	"os"
	"github.com/go-gl/gl"
	"image"
	"image/draw"
	"image/png"
)

var meshComponentName string = "MeshComponent"

type MeshComponent struct {
	Component
	vertexBuffer gl.Buffer
	indexBuffer gl.Buffer
	numIndices int16
	texcoordBuffer gl.Buffer
	showTexture bool
	texture gl.Texture
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

func (c *MeshComponent) SetIndices(indexData []int16) {
	c.indexBuffer.Delete()
	c.indexBuffer = gl.GenBuffer()
	c.indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)

	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indexData) * 2, indexData, gl.STATIC_DRAW)
	c.numIndices = int16(len(indexData))
}

func (c *MeshComponent) SetTextureCoordinates(texcoordData []float32) {
	c.texcoordBuffer.Delete()
	c.texcoordBuffer = gl.GenBuffer()
	c.texcoordBuffer.Bind(gl.ARRAY_BUFFER)

	gl.BufferData(gl.ARRAY_BUFFER, len(texcoordData) * 4, texcoordData, gl.STATIC_DRAW)
}

func (c *MeshComponent) ShowTexture(show bool) {
	c.showTexture = show
}

func (c *MeshComponent) SetTexture(filename string) {
	fi, err := os.Open("share/" + filename)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.texture.Delete()
	c.texture = gl.GenTexture()
	c.texture.Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	img, err := png.Decode(fi)
	if err != nil {
		panic(err)
	}
	copyimg := image.NewRGBA(img.Bounds())
	draw.Draw(copyimg, img.Bounds(), img, image.Pt(0, 0), draw.Src)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, copyimg.Bounds().Dx(), copyimg.Bounds().Dy(),
			0, gl.RGBA, gl.UNSIGNED_BYTE, copyimg.Pix)
}

func (c *MeshComponent) Update() {
	if c.numIndices == 0 {
		return
	}

	var vertexAttrib gl.AttribLocation = c.Graphics.program.GetAttribLocation("aPosition")
	vertexAttrib.EnableArray()
	c.vertexBuffer.Bind(gl.ARRAY_BUFFER)
	vertexAttrib.AttribPointer(3, gl.FLOAT, false, 0, nil)

	var texcoordAttrib gl.AttribLocation = c.Graphics.program.GetAttribLocation("aTexcoord")
	texcoordAttrib.EnableArray()
	c.texcoordBuffer.Bind(gl.ARRAY_BUFFER)
	texcoordAttrib.AttribPointer(2, gl.FLOAT, false, 0, nil)

	var uLoc = c.Graphics.program.GetUniformLocation("uTextured")
	if c.showTexture {
		uLoc.Uniform1i(1)
		uLoc = c.Graphics.program.GetUniformLocation("sTexture")
		uLoc.Uniform1i(0)
		gl.ActiveTexture(gl.TEXTURE0)
		c.texture.Bind(gl.TEXTURE_2D)
	} else {
		uLoc.Uniform1i(0)
	}

	uLoc = c.Graphics.program.GetUniformLocation("uPosition")
	var pos = c.GetTransform().Position
	uLoc.Uniform2f(float32(pos.X), float32(pos.Y))

	c.indexBuffer.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.DrawElements(gl.TRIANGLES, int(c.numIndices), gl.UNSIGNED_SHORT, nil)

	var err = gl.GetError()
	if err != 0 {
		fmt.Println("GL error", err)
	}
}


