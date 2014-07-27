package main

var ComponentNameMap = make(map[string]func() Componenter)

type Component struct {
	obj *GameObject
}

func (c *Component) GetObject() *GameObject {
	return c.obj
}

func (c *Component) SetObject(obj *GameObject) {
	c.obj = obj
}

func (c *Component) InternalInit(game *GameApp) {
}

func (c *Component) Start() {
}

func (c *Component) PreUpdate() {
}

func (c *Component) PostUpdate() {
}

func (c *Component) Update() {
}

func (c *Component) GetTransform() *TransformComponent {
	var t = c.obj.components[0]
	if (*t).Name() != "TransformComponent" {
		panic("First component must be transform, but it is " + (*t).Name())
	}
	var tt = *t
	return tt.(*TransformComponent)
}

type Componenter interface {
	Named
	GetTransform() *TransformComponent
	InternalInit(game *GameApp)
	Start()
	PreUpdate()
	PostUpdate()
	Update()
	GetObject() *GameObject
	SetObject(obj *GameObject)
}

