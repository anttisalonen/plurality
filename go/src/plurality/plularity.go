package main

import (
	"fmt"
	"encoding/json"
	"os"
	"io/ioutil"
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
	var filecontents []byte
	filecontents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(filecontents, &jsonData); err != nil {
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
				switch typ := jvaluedata.(type) {
				case string:
					reflect.ValueOf(compInst).Elem().FieldByName(jvaluename).SetString(jvaluedata.(string))
				case int:
					reflect.ValueOf(compInst).Elem().FieldByName(jvaluename).SetInt(jvaluedata.(int64))
				case float64:
					reflect.ValueOf(compInst).Elem().FieldByName(jvaluename).SetInt(int64(jvaluedata.(float64)))
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
