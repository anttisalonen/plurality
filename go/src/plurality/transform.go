package main

var transformComponentName string = "TransformComponent"

type TransformComponent struct {
	Component
	Position Vector2
}

func (c *TransformComponent) Name() string {
	return transformComponentName
}

func init() {
	ComponentNameMap[transformComponentName] = func() Componenter { return &TransformComponent{} }
}

