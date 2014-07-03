package main

import (
	"fmt"
	"encoding/json"
	"os"
	"reflect"
)

type Named interface {
	Name() string
}

type GameObject struct {
	name string
	components []*Componenter
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: %s <game JSON file>", os.Args[0])
		os.Exit(1)
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
}
