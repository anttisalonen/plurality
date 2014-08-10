package plurality

var ComponentNameMap = make(map[string]func() Componenter)

type Component struct {
	Object *GameObject
	gameApp *GameApp
	Graphics *Graphics
	Input *Input
	Time *Time
}

func (c *Component) SetObject(obj *GameObject) {
	c.Object = obj
	c.gameApp = obj.gameApp
	c.Graphics = &c.gameApp.graphics
	c.Input = &c.gameApp.input
	c.Time = &c.gameApp.time
}

func (c *Component) GetObjectByName(objname string) *GameObject {
	var obj = c.gameApp.objMap[objname]
	return obj
}

func (c *Component) Instantiate(objtype string, pos Vector2) *GameObject {
	return c.gameApp.Instantiate(objtype, pos)
}

func (c *Component) internalInit() {
}

func (c *Component) Start() {
}

func (c *Component) PreUpdate() {
}

func (c *Component) PostUpdate() {
}

func (c *Component) Update() {
}

func (c *Component) GetTransform() *TransformComponent {
	return c.Object.GetTransform()
}

type Componenter interface {
	Named
	GetTransform() *TransformComponent
	internalInit()
	Start()
	PreUpdate()
	PostUpdate()
	Update()
	SetObject(obj *GameObject)
}

