#!/usr/bin/env python2

import sys
import json
import wx

class Editor(wx.Frame):
    def __init__(self, *args, **kwargs):
        self.model = kwargs['model']
        del kwargs['model']
        super(Editor, self).__init__(*args, **kwargs) 
        self.ComponentWidgets = list()
        self.selectedObject = None
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
        self.hbox = wx.BoxSizer(wx.HORIZONTAL)
        self.vbox1 = wx.BoxSizer(wx.VERTICAL)
        self.hbox2 = wx.BoxSizer(wx.HORIZONTAL)
        self.vbox2 = wx.BoxSizer(wx.VERTICAL)
        st1 = wx.StaticText(self.panel, label='New Game Object')
        self.vbox1.Add(st1, flag=wx.RIGHT, border=8)

        self.newObjCtrl = wx.TextCtrl(self.panel)
        self.vbox1.Add(self.newObjCtrl)
        self.vbox1.Add(self.hbox2)
        self.hbox.Add(self.vbox1, flag=wx.EXPAND|wx.LEFT|wx.TOP, border=10)
        self.hbox.Add(self.vbox2, flag=wx.EXPAND|wx.RIGHT, border=10)
        self.panel.SetSizerAndFit(self.hbox)

        self.newObjCtrl.Bind(wx.EVT_KEY_UP, self.OnNewObject)

        self.newCompCtrl = wx.Button(self.panel, label='Add component')
        self.vbox1.Add(self.newCompCtrl)
        self.newCompCtrl.Bind(wx.EVT_BUTTON, self.OnNewComponent)
        self.compTypeCtrl = wx.ComboBox(self.panel, choices=self.model.getAvailableComponentTypes(), style=wx.CB_READONLY)
        self.vbox1.Add(self.compTypeCtrl)
        self.compTypeCtrl.SetSelection(0)

        self.tree = wx.TreeCtrl(self.panel, size=(400,300))
        self.tree.Bind(wx.EVT_TREE_SEL_CHANGED, self.OnActiveTree)
        self.hbox2.Add(self.tree, flag=wx.EXPAND|wx.LEFT|wx.RIGHT|wx.TOP|wx.BOTTOM, border=10)
        self.updateGUI()

    def OnActiveTree(self, e):
        tid = self.tree.GetSelection()
        objname = self.tree.GetItemText(tid)
        try:
            self.selectedObject = self.model.objects[objname]
        except KeyError:
            self.selectedObject = None
        finally:
            self.updateComponentView()

    def OnNewComponent(self, e):
        self.addComponent(self.compTypeCtrl.GetValue())

    def OnNewObject(self, e):
        key = e.GetKeyCode()
        if key == wx.WXK_RETURN:
            self.addObject(self.newObjCtrl.GetValue())

    def addComponent(self, comptype):
        if self.model.addComponent(self.selectedObject['name'], comptype):
            self.updateGUI()

    def addObject(self, objname):
        if self.model.addObject(objname):
            self.newObjCtrl.SetValue("")
            self.updateGUI()

    def OnSave(self, e):
        self.model.save()

    def OnQuit(self, e):
        self.Close()

    def updateGUI(self):
        self.tree.DeleteAllItems()
        self.treeRoot = self.tree.AddRoot("Objects")
        for objname, obj in sorted(self.model.objects.items()):
            tid = self.tree.AppendItem(self.treeRoot, objname)
        self.tree.ExpandAll()
        self.updateComponentView()

    def updateComponentView(self):
        for widget in self.ComponentWidgets:
            widget.Destroy()
        self.ComponentWidgets = list()

        if self.selectedObject:
            objname = self.selectedObject['name']
            self.ComponentWidgets.append(wx.StaticText(self.panel, label=objname))

            for comp in sorted(self.selectedObject['components'], key = lambda x: x['type']):
                compname = comp['type']
                complayout = self.model.components[compname]
                self.ComponentWidgets.append(wx.StaticText(self.panel, label=compname))

                b = wx.Button(self.panel, label='delete')
                b.Bind(wx.EVT_BUTTON, lambda event, compdata=(objname, compname): self.RemoveComponent(event, compdata))
                self.ComponentWidgets.append(b)

                for vname, vtype in sorted(complayout['values'].items()):
                    self.ComponentWidgets.append(wx.StaticText(self.panel, label=vname))
                    w = wx.TextCtrl(self.panel)
                    w.SetValue(str(comp['values'][vname]))
                    for ev in [wx.EVT_KEY_UP, wx.EVT_KILL_FOCUS]:
                        w.Bind(ev, lambda event, compdata=(objname,compname,vname,vtype,w): self.ComponentChanged(event, compdata))
                    self.ComponentWidgets.append(w)

        for widget in self.ComponentWidgets:
            self.vbox2.Add(widget, wx.EXPAND)
        self.panel.SetSizerAndFit(self.hbox)

    def ComponentChanged(self, e, compdata):
        objname, compname, vname, vtype, w = compdata
        newVal = w.GetValue()
        self.model.setComponentValue(objname, compname, vname, vtype, newVal)

    def RemoveComponent(self, e, compdata):
        objname, compname = compdata
        self.model.removeComponent(objname, compname)
        self.updateComponentView()

class Model(object):
    def __init__(self, compdata, gamedata, gamefilename):
        self.gamefilename = gamefilename
        self.components = dict()
        for c in compdata['components']:
            self.components[c['name']] = c
        self.objects = dict()
        for o in gamedata['objects']:
            self.objects[o['name']] = o

    def getAvailableComponentTypes(self):
        return [c['name'] for c in self.components.values()]

    def addComponent(self, objname, comptype):
        obj = self.objects[objname]

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

    def addObject(self, objname):
        if objname and objname not in self.objects:
            obj = dict()
            obj['name'] = objname
            obj['components'] = list()
            self.objects[objname] = obj
            return True
        else:
            return False

    def setComponentValue(self, objname, comptype, vname, vtype, newVal):
        def convert(val, t):
            if t == 'string':
                return val
            elif t == 'int':
                return int(val)
            elif t == 'bool':
                return bool(val)

        obj = self.objects[objname]
        for c in obj['components']:
            if c['type'] == comptype:
                assert vname in c['values']
                c['values'][vname] = convert(newVal, vtype)
                return
        assert False, 'Component %s not found in object %s' % (comptype, objname)

    def removeComponent(self, objname, comptype):
        obj = self.objects[objname]
        obj['components'] = [c for c in obj['components'] if c['type'] != comptype]

    def save(self):
        game = dict()
        game['objects'] = self.objects.values()
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

