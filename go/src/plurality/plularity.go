package main

import (
	"fmt"
	"encoding/json"
	"os"
	"reflect"
	"io/ioutil"
	"runtime"
)

type Named interface {
	Name() string
}

type GameObject struct {
	name string
	components []*Componenter
}

func main() {
	// OpenGL needs to be locked on this thread - see http://stackoverflow.com/questions/21010854/golang-fmt-println-causes-game-crash
	runtime.LockOSThread()
	if len(os.Args) < 2 {
		fmt.Println("Usage: %s <game JSON file>", os.Args[0])
		os.Exit(1)
	}

	if os.Args[1] == "-o" && len(os.Args) == 3 {
		outputInterfaceAndExit(os.Args[2])
	}

	filename := os.Args[1]
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var jsonData map[string]interface{}
	dec := json.NewDecoder(f)
	dec.UseNumber()
	if err := dec.Decode(&jsonData); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	loadGame(jsonData)
}

func loadGame(jsonData map[string]interface{}) {
	var objects = []GameObject{}
	objs := jsonData["objects"].([]interface{})
	for _, jobj := range objs {
		objmap := jobj.(map[string]interface{})
		var obj GameObject
		obj.name = objmap["name"].(string)
		components := objmap["components"].([]interface{})
		for _, jcomp := range components {
			comp := jcomp.(map[string]interface{})
			typeName := comp["type"]
			compInst := ComponentNameMap[typeName.(string)]()
			obj.components = append(obj.components, &compInst)

			for jvaluename, jvaluedata := range comp["values"].(map[string]interface{}) {
				compValue := reflect.ValueOf(compInst).Elem()
				fieldValue := compValue.FieldByName(jvaluename)
				typ := fieldValue.Kind()
				switch typ {
				case reflect.Bool:
					fieldValue.SetBool(jvaluedata.(bool))
				case reflect.Int:
					v, err := jvaluedata.(json.Number).Int64()
					if err != nil {
						fmt.Println("Error on field %s: %s", jvaluename, err)
					}
					fieldValue.SetInt(v)
				case reflect.String:
					fieldValue.SetString(jvaluedata.(string))
				default:
					fmt.Println("Unknown type", typ)
				}
			}
		}
		objects = append(objects, obj)
	}

	runGame(objects)
}

func runGame(objects []GameObject) {
	for _, obj := range objects {
		for _, comp := range obj.components {
			(*comp).Start()
		}
	}

	for i := 0; i < 50; i++ {
		for _, obj := range objects {
			for _, comp := range obj.components {
				(*comp).PreUpdate()
			}
		}

		for _, obj := range objects {
			for _, comp := range obj.components {
				(*comp).Update()
			}
		}

		for _, obj := range objects {
			for _, comp := range obj.components {
				(*comp).PostUpdate()
			}
		}
	}
}

func outputInterfaceAndExit(filename string) {
	v := make(map[string]interface{})
	var components []interface{}
	for k, v := range ComponentNameMap {
		comp := make(map[string]interface{})
		comp["name"] = k
		values := make(map[string]string)
		compInst := v()
		compType := reflect.TypeOf(compInst).Elem()
		for i := 0; i < compType.NumField(); i++ {
			p := compType.Field(i)
			if p.Anonymous || p.PkgPath != "" {
				continue
			}
			t := p.Type.Kind()
			var str string
			switch t {
			case reflect.Int:
				str = "int"
			case reflect.String:
				str = "string"
			case reflect.Bool:
				str = "bool"
			default:
				fmt.Println("Unknown type", t, "for", p.Name)
			}
			if len(str) > 0 {
				values[p.Name] = str
			}
		}
		comp["values"] = values
		components = append(components, comp)
	}
	v["components"] = components
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
