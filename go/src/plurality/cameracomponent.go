package main

import (
	"github.com/banthar/Go-SDL/sdl"
)

var cameraComponentName string = "CameraComponent"

type CameraComponent struct {
	Component
}

func (c *CameraComponent) Name() string {
	return cameraComponentName
}

func (c *CameraComponent) Start() {
	sdl.Init(sdl.INIT_EVERYTHING)
	//_ = sdl.CreateRGBSurface(0, 640, 480, 32, 0, 0, 0, 0)
	_ = sdl.SetVideoMode(640, 480, 32, 0)
}

func init() {
	ComponentNameMap[cameraComponentName] = func() Componenter { return &CameraComponent{} }
}

