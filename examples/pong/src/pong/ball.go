package main

import (
	"plurality"
	"fmt"
)

var ballComponentName string = "ball"

type ball struct {
	plurality.Component
	Speed float64
	velocity plurality.Vector2
	score1 int
	score2 int
}

func (c *ball) Name() string {
	return ballComponentName
}

func (c *ball) Start() {
	c.velocity.X = c.Speed
	c.velocity.Y = c.Speed
}

func (c *ball) Update() {
	var pos = c.GetTransform().Position
	pos = c.GetTransform().Position.Add(c.velocity.Multiplied(c.Time.GetDeltaTime()))
	if pos.X > float64(c.Graphics.ScreenWidth) * 0.5 {
		pos.X = 0
		c.velocity.X = -c.velocity.X
		c.score1++
		fmt.Printf("Score: %d - %d\n", c.score1, c.score2)
	}
	if pos.X < float64(-c.Graphics.ScreenWidth) * 0.5 {
		pos.X = 0
		c.velocity.X = -c.velocity.X
		c.score2++
		fmt.Printf("Score: %d - %d\n", c.score1, c.score2)
	}
	if pos.Y > float64(c.Graphics.ScreenHeight) * 0.5 && c.velocity.Y > 0.0 {
		c.velocity.Y = -c.velocity.Y
	}
	if pos.Y < float64(-c.Graphics.ScreenHeight) * 0.5 && c.velocity.Y < 0.0 {
		c.velocity.Y = -c.velocity.Y
	}

	var paddle1 = c.GetObjectByName("Paddle1")
	var paddle1pos = paddle1.GetTransform().Position
	var paddle1size = paddle1.GetComponent("RectangleComponent").(*plurality.RectangleComponent).Scale
	if c.velocity.X < 0.0 && pos.X < paddle1pos.X + paddle1size.X * 0.5 &&
	pos.X > paddle1pos.X - paddle1size.X * 0.5 &&
	pos.Y < paddle1pos.Y + paddle1size.Y * 0.5 &&
	pos.Y > paddle1pos.Y - paddle1size.Y * 0.5 {
		c.velocity.X = -c.velocity.X
	}

	var paddle2 = c.GetObjectByName("Paddle2")
	var paddle2pos = paddle2.GetTransform().Position
	var paddle2size = paddle2.GetComponent("RectangleComponent").(*plurality.RectangleComponent).Scale
	if c.velocity.X > 0.0 && pos.X > paddle2pos.X - paddle2size.X * 0.5 &&
	pos.X < paddle2pos.X + paddle2size.X * 0.5 &&
	pos.Y < paddle2pos.Y + paddle2size.Y * 0.5 &&
	pos.Y > paddle2pos.Y - paddle2size.Y * 0.5 {
		c.velocity.X = -c.velocity.X
	}

	c.GetTransform().Position = pos
}

func init() {
	plurality.ComponentNameMap[ballComponentName] = func() plurality.Componenter { return &ball{} }
}

