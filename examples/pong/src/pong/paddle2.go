package main

import (
	"plurality"
)

var paddle2ComponentName string = "paddle2"

func init() {
	plurality.ComponentNameMap[paddle2ComponentName] = func() plurality.Componenter { return &paddle2{} }
}

type paddle2 struct {
	plurality.Component
	Speed float64
}

func (c *paddle2) Name() string {
	return paddle2ComponentName
}

func (c *paddle2) Update() {
	var b = c.GetObjectByName("Ball")
	var ballpos = b.GetTransform().Position.Y
	var ballvelx = b.GetComponent("ball").(*ball).velocity.X
	var myposy = c.GetTransform().Position.Y
	var scaley = c.Object.GetComponent("RectangleComponent").(*plurality.RectangleComponent).Scale.Y
	if ballvelx > 0.0 && myposy - scaley * 0.4 < ballpos {
		myposy += 1 * c.Time.GetDeltaTime() * c.Speed
	}
	if ballvelx > 0.0 && myposy + scaley * 0.4 > ballpos {
		myposy -= 1 * c.Time.GetDeltaTime() * c.Speed
	}

	c.GetTransform().Position.Y = clampPos(myposy, scaley, c.Graphics.ScreenHeight)
}


