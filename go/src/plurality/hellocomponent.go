package main

import (
	"fmt"
)

var helloComponentName string = "HelloComponent"

type HelloComponent struct {
	Greetee string
	NumGreets int
	Component
}

func (c *HelloComponent) Name() string {
	return helloComponentName
}

func (c *HelloComponent) Start() {
	for i := 0; i < c.NumGreets; i++ {
		fmt.Printf("Hello %s!\n", c.Greetee)
	}
}

func init() {
	ComponentNameMap[helloComponentName] = func() Componenter { return &HelloComponent{} }
}

