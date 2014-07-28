package main

import (
	"plurality"
)

var paddle2ComponentName string = "Paddle2Component"

func init() {
	plurality.ComponentNameMap[paddle2ComponentName] = func() plurality.Componenter { return &Paddle2Component{} }
}

type Paddle2Component struct {
	plurality.Component
	Speed float64
}

func (c *Paddle2Component) Name() string {
	return paddle2ComponentName
}

func (c *Paddle2Component) Update() {
	var ball = c.GetObjectByName("Ball")
	var ballpos = ball.GetTransform().Position.Y
	var ballvelx = ball.GetComponent("BallComponent").(*BallComponent).velocity.X
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


