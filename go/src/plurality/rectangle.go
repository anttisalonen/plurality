package plurality

var rectangleComponentName string = "RectangleComponent"

type RectangleComponent struct {
	Component
	Scale Vector2
}

func (c *RectangleComponent) Name() string {
	return rectangleComponentName
}

func (c *RectangleComponent) Start() {
	var mesh = c.Object.GetComponent("MeshComponent")
	if mesh == nil {
		mesh = c.Object.AddComponent("MeshComponent")
	}
	var meshc = mesh.(*MeshComponent)
	var vertexBufferData = []float32{
		-0.5 * float32(c.Scale.X), -0.5 * float32(c.Scale.Y), 0.0,
		-0.5 * float32(c.Scale.X),  0.5 * float32(c.Scale.Y), 0.0,
		 0.5 * float32(c.Scale.X),  0.5 * float32(c.Scale.Y), 0.0,
		 0.5 * float32(c.Scale.X), -0.5 * float32(c.Scale.Y), 0.0}
	meshc.SetVertices(vertexBufferData)
	var indexBufferData = []int16{0, 2, 1, 0, 3, 2}
	meshc.SetIndices(indexBufferData)
}

func init() {
	ComponentNameMap[rectangleComponentName] = func() Componenter { return &RectangleComponent{} }
}

