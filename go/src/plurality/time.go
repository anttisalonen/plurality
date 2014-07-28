package plurality

import (
	"time"
)

type Time struct {
	lastTime time.Time
	deltaTime float64 // seconds
}

func (c *Time) GetDeltaTime() float64 {
	return c.deltaTime
}

func (c *Time) Init() {
	c.lastTime = time.Now()
}

func (c *Time) Update() bool {
	var prevTime = c.lastTime
	c.lastTime = time.Now()
	c.deltaTime = c.lastTime.Sub(prevTime).Seconds()
	return true
}
