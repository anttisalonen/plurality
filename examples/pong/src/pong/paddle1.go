package main

import (
	"plurality"
	"math"
)

var paddle1ComponentName string = "paddle1"

type paddle1 struct {
	plurality.Component
	Speed float64
}

func (c *paddle1) Name() string {
	return paddle1ComponentName
}

func clamp(t, mn, mx float64) float64 {
	return math.Max(mn, math.Min(t, mx))
}

func clampPos(posy, scaley float64, screenheight int) float64 {
	return clamp(posy, float64(-screenheight) * 0.5 + scaley * 0.5,
	float64(screenheight) * 0.5 - scaley * 0.5)
}

func (c *paddle1) Start() {
}

func (c *paddle1) Update() {
	var inp = c.Input.GetAxis("Vertical")
	if inp != 0.0 {
		c.GetTransform().Position.Y += inp * c.Time.GetDeltaTime() * c.Speed
	}

	var scale = c.Object.GetComponent("RectangleComponent").(*plurality.RectangleComponent).Scale
	c.GetTransform().Position.Y = clampPos(c.GetTransform().Position.Y, scale.Y, c.Graphics.ScreenHeight)
}

func init() {
	plurality.ComponentNameMap[paddle1ComponentName] = func() plurality.Componenter { return &paddle1{} }
}

