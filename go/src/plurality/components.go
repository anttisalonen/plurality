package main

var ComponentNameMap = make(map[string]func() Componenter)

type Component struct {
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

type Componenter interface {
	Named
	InternalInit(game *GameApp)
	Start()
	PreUpdate()
	PostUpdate()
	Update()
}

