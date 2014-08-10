package plurality

import (
	"fmt"
	"encoding/json"
	"os"
	"reflect"
	"io/ioutil"
	"runtime"
	"strconv"
)

type Named interface {
	Name() string
}

type GameObject struct {
	name string
	components []Componenter
	gameApp *GameApp
}

func (g *GameObject) GetComponent(comptype string) Componenter {
	for _, c := range g.components {
		if c.Name() == comptype {
			return c
		}
	}
	return nil
}

func (o *GameObject) GetTransform() *TransformComponent {
	var t = o.components[0]
	if t.Name() != "TransformComponent" {
		panic("First component must be transform, but it is " + t.Name())
	}
	var tt = t.(*TransformComponent)
	return tt
}

type GameApp struct {
	objMap map[string]*GameObject
	graphics Graphics
	input Input
	time Time
	prefabMap map[string]interface{}
	nextPrefabIndex int
}

func (g *GameApp) Instantiate(objtype string, pos Vector2) *GameObject {
	var objinst = prefabToObject(g, g.prefabMap[objtype].(map[string]interface{}))

	objinst.GetTransform().Position = pos
	objinst.name = "__fab" + strconv.Itoa(g.nextPrefabIndex)
	g.objMap[objinst.name] = objinst
	g.nextPrefabIndex++

	for _, comp := range objinst.components {
		comp.SetObject(objinst)
		comp.internalInit()
		comp.Start()
	}
	return objinst
}

type Vector2 struct {
	X float64
	Y float64
}

func (v Vector2) Add(v2 Vector2) Vector2 {
	return Vector2{v.X + v2.X, v.Y + v2.Y}
}

func (v Vector2) Multiplied(a float64) Vector2 {
	return Vector2{v.X * a, v.Y * a}
}

func Main() {
	// OpenGL needs to be locked on this thread - see http://stackoverflow.com/questions/21010854/golang-fmt-println-causes-game-crash
	// app's main function should do this but have it here as well just in case
	runtime.LockOSThread()
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <game JSON file>\n", os.Args[0])
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

	var game GameApp
	game.Prepare()
	game.Run(jsonData)
}

func (game *GameApp) Prepare() {
	game.graphics.Init(800, 600)
	game.input.Init()
	game.time.Init()
}

func (game *GameApp) Run(jsonData map[string]interface{}) {
	var objects, prefabMap = loadGame(game, jsonData)
	game.objMap = make(map[string]*GameObject)
	for i, _ := range objects {
		var on = objects[i].name
		game.objMap[on] = objects[i]
	}

	game.prefabMap = prefabMap

	runGame(game)
}

func loadGame(game* GameApp, jsonData map[string]interface{}) ([]*GameObject, map[string]interface{}) {
	var objs = loadObjects(game, jsonData["objects"].([]interface{}))
	var prefabMap = loadPrefabs(jsonData["prefabs"].([]interface{}))
	return objs, prefabMap
}

func loadPrefabs(objs []interface{}) map[string]interface{} {
	var prefabMap = make(map[string]interface{})
	for _, jobj := range objs {
		objmap := jobj.(map[string]interface{})
		var prefab GameObject
		prefab.name = objmap["name"].(string)
		prefabMap[objmap["name"].(string)] = objmap
	}
	return prefabMap
}

func loadObjects(game* GameApp, objs []interface{}) []*GameObject {
	var objects = []*GameObject{}
	for _, jobj := range objs {
		objmap := jobj.(map[string]interface{})
		var obj = prefabToObject(game, objmap)
		objects = append(objects, obj)
	}

	return objects
}

func prefabToObject(game *GameApp, objmap map[string]interface{}) *GameObject {
	var obj GameObject
	obj.name = objmap["name"].(string)
	obj.gameApp = game
	components := objmap["components"].([]interface{})
	for _, jcomp := range components {
		comp := jcomp.(map[string]interface{})
		typeName := comp["type"]
		compInst := ComponentNameMap[typeName.(string)]()
		compInst.SetObject(&obj)
		obj.components = append(obj.components, compInst)

		for jvaluename, jvaluedata := range comp["values"].(map[string]interface{}) {
			compValue := reflect.ValueOf(compInst).Elem()
			fieldValue := compValue.FieldByName(jvaluename)
			typ := fieldValue.Kind()
			switch typ {
			case reflect.Bool:
				fieldValue.SetBool(jvaluedata.(bool))
			case reflect.Float64:
				v, err := jvaluedata.(json.Number).Float64()
				if err != nil {
					fmt.Printf("Error on field %s: %s\n", jvaluename, err)
				}
				fieldValue.SetFloat(v)
			case reflect.Int:
				v, err := jvaluedata.(json.Number).Int64()
				if err != nil {
					fmt.Printf("Error on field %s: %s\n", jvaluename, err)
				}
				fieldValue.SetInt(v)
			case reflect.String:
				fieldValue.SetString(jvaluedata.(string))
			case reflect.Struct:
				readVector2(&fieldValue, jvaluedata)
			default:
				fmt.Println("Unknown type", typ, "for", jvaluename, "at", typeName, "in", obj.name)
			}
		}
	}
	return &obj
}

func readVector2(fieldValue *reflect.Value, jvaluedata interface{}) {
	var vec2 Vector2
	if fieldValue.Type() == reflect.TypeOf(vec2) {
		var jv = jvaluedata.([]interface{})
		v, err := jv[0].(json.Number).Float64()
		if err != nil {
			panic(err)
		}
		fieldValue.FieldByName("X").SetFloat(v)

		v, err = jv[1].(json.Number).Float64()
		if err != nil {
			panic(err)
		}
		fieldValue.FieldByName("Y").SetFloat(v)
	}
}

func runGame(game *GameApp) {
	var objects = game.objMap
	for _, obj := range objects {
		for _, comp := range obj.components {
			comp.internalInit()
		}
	}

	for _, obj := range objects {
		for _, comp := range obj.components {
			comp.Start()
		}
	}

	for {
		for _, obj := range objects {
			for _, comp := range obj.components {
				comp.PreUpdate()
			}
		}

		for _, obj := range objects {
			for _, comp := range obj.components {
				comp.Update()
			}
		}

		for _, obj := range objects {
			for _, comp := range obj.components {
				comp.PostUpdate()
			}
		}

		game.graphics.Update()
		var running = game.input.Update()
		if !running {
			break
		}

		running = game.time.Update()
		if !running {
			break
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
			case reflect.Float64:
				str = "float64"
			case reflect.Struct:
				str = "Vector2"
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
