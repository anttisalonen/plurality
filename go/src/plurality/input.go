package plurality

import (
	"github.com/banthar/Go-SDL/sdl"
)

type Input struct {
	axis map[string]float64
}

func (c *Input) GetAxis(axisname string) float64 {
	return c.axis[axisname]
}

func (c *Input) Init() {
	c.axis = make(map[string]float64)
	c.axis["Vertical"] = 0.0
	c.axis["Horizontal"] = 0.0
}

func (c *Input) Update() bool {
	for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
		switch e := ev.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			switch e.Type {
			case sdl.KEYDOWN:
				switch e.Keysym.Sym {
				case sdl.K_ESCAPE: // escape
					return false
				case sdl.K_UP:
					c.axis["Vertical"] = 1.0
				case sdl.K_DOWN:
					c.axis["Vertical"] = -1.0
				case sdl.K_RIGHT:
					c.axis["Horizontal"] = 1.0
				case sdl.K_LEFT:
					c.axis["Horizontal"] = -1.0
				}
			case sdl.KEYUP:
				switch e.Keysym.Sym {
				case sdl.K_UP:
					c.axis["Vertical"] = 0.0
				case sdl.K_DOWN:
					c.axis["Vertical"] = 0.0
				case sdl.K_RIGHT:
					c.axis["Horizontal"] = 0.0
				case sdl.K_LEFT:
					c.axis["Horizontal"] = 0.0
				}
			}
		}
	}

	return true
}
