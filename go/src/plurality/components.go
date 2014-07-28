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
}

func (c *Component) GetObjectByName(objname string) *GameObject {
	var obj = c.gameApp.objMap[objname]
	return &obj
}

func (c *Component) SetGame(g *GameApp) {
	c.gameApp = g
	c.Graphics = &g.graphics
	c.Input = &g.input
	c.Time = &g.time
}

func (c *Component) InternalInit(game *GameApp) {
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
	InternalInit(game *GameApp)
	Start()
	PreUpdate()
	PostUpdate()
	Update()
	SetObject(obj *GameObject)
	SetGame(g *GameApp)
}

