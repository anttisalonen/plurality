package main

import "fmt"

var componentName string = "HelloComponent"

type HelloComponent struct {
	Greetee string
	NumGreets int
}

func (c *HelloComponent) Name() string {
	return componentName
}

func (c *HelloComponent) Start() {
	for i := 0; i < c.NumGreets; i++ {
		fmt.Printf("Hello %s!\n", c.Greetee)
	}
}

func init() {
	ComponentNameMap[componentName] = func() Componenter { return &HelloComponent{} }
}

