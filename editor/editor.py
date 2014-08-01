#!/usr/bin/env python2

import sys
import json
import copy
import wx

class ObjectType(object):
    GameObject = 0
    Prefab = 1

class Editor(wx.Frame):
    def __init__(self, *args, **kwargs):
        self.model = kwargs['model']
        del kwargs['model']
        super(Editor, self).__init__(*args, **kwargs) 
        self.ComponentWidgets = list()
        self.selectedObject = None
        self.selectedObjectType = None
        menubar = wx.MenuBar()
        fileMenu = wx.Menu()
        saveItem = fileMenu.Append(wx.ID_SAVE, 'Save', 'Save')
        fitem = fileMenu.Append(wx.ID_EXIT, 'Quit', 'Quit application')
        menubar.Append(fileMenu, '&File')
        self.SetMenuBar(menubar)

        self.Bind(wx.EVT_MENU, self.OnSave, saveItem)
        self.Bind(wx.EVT_MENU, self.OnQuit, fitem)

        self.SetSize((800, 600))
        self.SetTitle('Editor')
        self.Centre()
        self.Show(True)

        self.panel = wx.Panel(self)
        self.mainbox = wx.BoxSizer(wx.HORIZONTAL)
        self.leftbox = wx.BoxSizer(wx.VERTICAL)
        self.rightbox = wx.BoxSizer(wx.VERTICAL)
        st1 = wx.StaticText(self.panel, label='New Game Object')
        self.leftbox.Add(st1, flag=wx.RIGHT, border=8)

        self.newObjCtrl = wx.TextCtrl(self.panel)
        self.leftbox.Add(self.newObjCtrl)

        self.mainbox.Add(self.leftbox, flag=wx.EXPAND|wx.LEFT|wx.TOP, border=10)
        self.mainbox.Add(self.rightbox, flag=wx.EXPAND|wx.RIGHT, border=10)
        self.panel.SetSizerAndFit(self.mainbox)

        self.newObjCtrl.Bind(wx.EVT_KEY_UP, self.OnNewObject)

        self.tree = wx.TreeCtrl(self.panel, size=(400,300))
        self.tree.Bind(wx.EVT_TREE_SEL_CHANGED, self.OnActiveTree)
        self.leftbox.Add(self.tree, flag=wx.EXPAND|wx.LEFT|wx.RIGHT|wx.TOP|wx.BOTTOM, border=10)

        self.removeObjButton = wx.Button(self.panel, label='Delete selected object')
        self.removeObjButton.Bind(wx.EVT_BUTTON, self.OnRemoveObject)
        self.leftbox.Add(self.removeObjButton)

        self.createPrefabButton = wx.Button(self.panel, label='Create prefab')
        self.createPrefabButton.Bind(wx.EVT_BUTTON, self.OnCreatePrefab)
        self.leftbox.Add(self.createPrefabButton)

        self.compTypeCtrl = wx.ComboBox(self.panel, choices=self.model.getAvailableComponentTypes(), style=wx.CB_READONLY)
        self.leftbox.Add(self.compTypeCtrl)
        self.compTypeCtrl.SetSelection(0)

        self.newCompCtrl = wx.Button(self.panel, label='Add component')
        self.leftbox.Add(self.newCompCtrl)
        self.newCompCtrl.Bind(wx.EVT_BUTTON, self.OnNewComponent)

        self.updateGUI()

    def OnActiveTree(self, e):
        tid = self.tree.GetSelection()
        objname = self.tree.GetItemText(tid)

        par = self.tree.GetItemParent(tid)
        if not par:
            self.postActiveTree()
            return

        ot = self.tree.GetItemText(par)
        if ot == "Objects":
            self.selectedObjectType = ObjectType.GameObject
        elif ot == "Prefabs":
            self.selectedObjectType = ObjectType.Prefab
        else: # root
            self.postActiveTree()
            return

        try:
            self.selectedObject = self.model.getObject(objname, self.selectedObjectType)
        except KeyError:
            self.selectedObject = None
            self.selectedObjectType = None
        finally:
            self.postActiveTree()

    def postActiveTree(self):
        self.updateComponentView()

    def OnNewComponent(self, e):
        self.addComponent(self.compTypeCtrl.GetValue())

    def OnNewObject(self, e):
        key = e.GetKeyCode()
        if key == wx.WXK_RETURN:
            self.addObject(self.newObjCtrl.GetValue())

    def addComponent(self, comptype):
        if self.model.addComponent(self.selectedObject['name'], self.selectedObjectType, comptype):
            self.updateComponentView()

    def addObject(self, objname):
        if self.model.addObject(objname):
            self.newObjCtrl.SetValue("")
            self.updateTreeView()

    def OnSave(self, e):
        self.model.save()

    def OnQuit(self, e):
        self.Close()

    def updateGUI(self):
        self.updateTreeView()
        self.updateComponentView()

    def updateTreeView(self):
        self.tree.DeleteAllItems()
        self.treeRoot = self.tree.AddRoot("Game")
        objects = self.tree.AppendItem(self.treeRoot, "Objects")
        for objname, obj in sorted(self.model.objects.items()):
            tid = self.tree.AppendItem(objects, objname)
        prefabs = self.tree.AppendItem(self.treeRoot, "Prefabs")
        for objname, obj in sorted(self.model.prefabs.items()):
            tid = self.tree.AppendItem(prefabs, objname)
        self.tree.ExpandAll()

    def updateComponentView(self):
        for widget in self.ComponentWidgets:
            widget.Destroy()
        self.panel.SetSizerAndFit(self.mainbox)
        self.ComponentWidgets = list()

        if self.selectedObject:
            objname = self.selectedObject['name']
            self.ComponentWidgets.append(wx.StaticText(self.panel, label=objname))

            for comp in sorted(self.selectedObject['components'], key = lambda x: x['type']):
                compname = comp['type']
                complayout = self.model.components[compname]
                self.ComponentWidgets.append(wx.StaticText(self.panel, label=compname))

                if compname != 'TransformComponent':
                    b = wx.Button(self.panel, label='Delete')
                    b.Bind(wx.EVT_BUTTON, lambda event, compdata=compname: self.RemoveComponent(event, compdata))
                    self.ComponentWidgets.append(b)

                for vname, vtype in sorted(complayout['values'].items()):
                    self.ComponentWidgets.append(wx.StaticText(self.panel, label=vname))
                    if vtype == 'Vector2':
                        widgets = list()

                        for i in xrange(2):
                            w = wx.TextCtrl(self.panel)
                            values = comp['values']
                            try:
                                vals = values[vname]
                            except KeyError: # inteface changed (new public variable)
                                values[vname] = self.model.getDefault(vtype)
                                vals = values[vname]
                            val = vals[i]
                            w.SetValue(str(val))
                            self.ComponentWidgets.append(w)
                            widgets.append(w)

                        for widget in widgets:
                            for ev in [wx.EVT_KEY_UP, wx.EVT_KILL_FOCUS]:
                                widget.Bind(ev, lambda event,
                                        compdata=(compname,vname,vtype,widgets): self.ComponentChanged(event, compdata))

                    else:
                        w = wx.TextCtrl(self.panel)
                        values = comp['values']
                        try:
                            val = values[vname]
                        except KeyError: # interface changed (new public variable)
                            values[vname] = self.model.getDefault(vtype)
                            val = values[vname]
                        w.SetValue(str(val))
                        for ev in [wx.EVT_KEY_UP, wx.EVT_KILL_FOCUS]:
                            w.Bind(ev, lambda event, compdata=(compname,vname,vtype,w): self.ComponentChanged(event, compdata))
                        self.ComponentWidgets.append(w)

        for widget in self.ComponentWidgets:
            self.rightbox.Add(widget, wx.EXPAND)
        self.panel.SetSizerAndFit(self.mainbox)

    def ComponentChanged(self, e, compdata):
        compname, vname, vtype, w = compdata
        if isinstance(w, list): # Vector2
            newVal = [widget.GetValue() for widget in w]
        else:
            newVal = w.GetValue()
        self.model.setComponentValue(self.selectedObject['name'], self.selectedObjectType, compname, vname, vtype, newVal)

    def RemoveComponent(self, e, compdata):
        compname = compdata
        self.model.removeComponent(self.selectedObject['name'], self.selectedObjectType, compname)
        self.updateComponentView()

    def OnRemoveObject(self, e):
        if self.selectedObject:
            self.model.removeObject(self.selectedObject['name'], self.selectedObjectType)
            self.selectedObject = None
            self.updateGUI()

    def OnCreatePrefab(self, e):
        if self.selectedObject:
            self.model.createPrefab(self.selectedObject['name'])
            self.selectedObject = None
            self.updateGUI()

class Model(object):
    def __init__(self, compdata, gamedata, gamefilename):
        self.gamefilename = gamefilename
        self.components = dict()
        for c in compdata['components']:
            self.components[c['name']] = c
        self.objects = dict()
        for o in gamedata['objects']:
            self.objects[o['name']] = o

        self.prefabs = dict()
        for o in gamedata.get('prefabs', dict()):
            self.prefabs[o['name']] = o

    def getAvailableComponentTypes(self):
        return [c['name'] for c in self.components.values()]

    def addComponent(self, objname, objtype, comptype):
        obj = self._getObjectOnType(objname, objtype)

        for c in obj['components']:
            if c['type'] == comptype:
                print 'Object already has component of type', comptype
                return False

        comp = dict()
        comp['type'] = comptype
        comp['values'] = dict()
        for valuename, valuetype in self.components[comptype]['values'].items():
            comp['values'][valuename] = self.getDefault(valuetype)
        obj['components'].append(comp)
        return True

    def getDefault(self, valuetype):
        if valuetype == 'string':
            return ''
        elif valuetype == 'int':
            return 0
        elif valuetype == 'bool':
            return False
        elif valuetype == 'float64':
            return 0.0
        elif valuetype == 'Vector2':
            return [0.0, 0.0]

    def addObject(self, objname):
        if objname and objname not in self.objects:
            obj = dict()
            obj['name'] = objname
            obj['components'] = [{'type':'TransformComponent', 'values':{'Position':[0.0, 0.0]}}]
            self.objects[objname] = obj
            return True
        else:
            return False

    def _getObjectOnType(self, objname, objtype):
        if objtype == ObjectType.GameObject:
            return self.objects[objname]
        else:
            return self.prefabs[objname]

    def getObject(self, objname, objtype):
        return self._getObjectOnType(objname, objtype)

    def setComponentValue(self, objname, objtype, comptype, vname, vtype, newVal):
        def convert(val, t):
            if t == 'string':
                return val
            elif t == 'int':
                return int(val)
            elif t == 'bool':
                return bool(val)
            elif t == 'float64':
                return float(val)
            elif t == 'Vector2':
                return [float(v) for v in newVal]

        obj = self._getObjectOnType(objname, objtype)

        for c in obj['components']:
            if c['type'] == comptype:
                assert vname in c['values']
                c['values'][vname] = convert(newVal, vtype)
                return
        assert False, 'Component %s not found in object %s' % (comptype, objname)

    def removeComponent(self, objname, objtype, comptype):
        obj = self._getObjectOnType(objname, objtype)
        obj['components'] = [c for c in obj['components'] if c['type'] != comptype]

    def removeObject(self, objname, objtype):
        if objtype == ObjectType.GameObject:
            del self.objects[objname]
        else:
            del self.prefabs[objname]

    def createPrefab(self, objname):
        self.prefabs[objname] = copy.deepcopy(self.objects[objname])

    def save(self):
        game = dict()
        game['objects'] = self.objects.values()
        game['prefabs'] = self.prefabs.values()
        with open(self.gamefilename, 'w') as f:
            f.write(json.dumps(game, indent=4))

def main():
    ed = wx.App()
    try:
        compfilename = sys.argv[1]
        gamefilename = sys.argv[2]
    except IndexError:
        print "Usage: %d <component information JSON file> <game JSON file>"
        sys.exit(1)
    compdata = json.loads(open(compfilename, 'r').read())
    try:
        gamedata = json.loads(open(gamefilename, 'r').read())
    except IOError:
        gamedata = {'objects':[]}
    model = Model(compdata, gamedata, gamefilename)
    e = Editor(None, model=model)
    ed.MainLoop()    

if __name__ == '__main__':
    main()

