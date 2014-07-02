package main

var ComponentNameMap = make(map[string]func() Componenter)

type Componenter interface {
	Named
	Start()
}

