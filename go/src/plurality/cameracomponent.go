package main

import (
	"github.com/banthar/Go-SDL/sdl"
)

var cameraComponentName string = "CameraComponent"

type CameraComponent struct {
	Component
	ScreenWidth int
	ScreenHeight int
}

func (c *CameraComponent) Name() string {
	return cameraComponentName
}

func (c *CameraComponent) Start() {
	sdl.Init(sdl.INIT_EVERYTHING)
	_ = sdl.SetVideoMode(c.ScreenWidth, c.ScreenHeight, 32, 0)
}

func init() {
	ComponentNameMap[cameraComponentName] = func() Componenter { return &CameraComponent{} }
}

