package main

import (
	"plurality"
)

var paddle1ComponentName string = "Paddle1Component"

type Paddle1Component struct {
	plurality.Component
	Speed float64
}

func (c *Paddle1Component) Name() string {
	return paddle1ComponentName
}

func (c *Paddle1Component) Start() {
}

func (c *Paddle1Component) Update() {
	var inp = c.Input.GetAxis("Vertical")
	if inp != 0.0 {
		c.GetTransform().Position.Y += inp * c.Time.GetDeltaTime() * c.Speed
	}
}

func init() {
	plurality.ComponentNameMap[paddle1ComponentName] = func() plurality.Componenter { return &Paddle1Component{} }
}

