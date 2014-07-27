package main

import (
	"plurality"
)

var paddle1ComponentName string = "Paddle1Component"

type Paddle1Component struct {
	plurality.Component
	Position plurality.Vector2
}

func (c *Paddle1Component) Name() string {
	return paddle1ComponentName
}

func init() {
	plurality.ComponentNameMap[paddle1ComponentName] = func() plurality.Componenter { return &Paddle1Component{} }
}

