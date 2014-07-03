#!/usr/bin/env python2

import sys
import json
import wx

class Editor(wx.Frame):
    def __init__(self, *args, **kwargs):
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

        self.tc = wx.TextCtrl(self.panel)
        self.vbox1.Add(self.tc)
        self.vbox1.Add(self.hbox2)
        self.hbox.Add(self.vbox1, flag=wx.EXPAND|wx.LEFT|wx.TOP, border=10)
        self.hbox.Add(self.vbox2, flag=wx.EXPAND|wx.RIGHT, border=10)
        self.panel.SetSizerAndFit(self.hbox)

        self.tc.Bind(wx.EVT_KEY_UP, self.OnNewObject)

        self.tree = wx.TreeCtrl(self.panel, size=(400,300))
        self.tree.Bind(wx.EVT_TREE_SEL_CHANGED, self.OnActiveTree)
        self.hbox2.Add(self.tree, flag=wx.EXPAND|wx.LEFT|wx.RIGHT|wx.TOP|wx.BOTTOM, border=10)

    def OnActiveTree(self, e):
        tid = self.tree.GetSelection()
        objname = self.tree.GetItemText(tid)
        try:
            self.selectedObject = self.model.objects[objname]
        except KeyError:
            self.selectedObject = None
        finally:
            self.updateComponentView()

    def OnNewObject(self, e):
        key = e.GetKeyCode()
        if key == wx.WXK_RETURN:
            self.addObject(self.tc.GetValue())

    def addObject(self, objname):
        if self.model.addObject(objname):
            self.tc.SetValue("")
            self.updateGUI()

    def OnSave(self, e):
        self.model.save()

    def OnQuit(self, e):
        self.Close()

    def setModel(self, m):
        self.model = m
        self.updateGUI()

    def updateGUI(self):
        self.tree.DeleteAllItems()
        self.treeRoot = self.tree.AddRoot("Objects")
        for objname, obj in sorted(self.model.objects.items()):
            tid = self.tree.AppendItem(self.treeRoot, objname)
        self.tree.ExpandAll()

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

class Model(object):
    def __init__(self, compdata, gamedata, gamefilename):
        self.gamefilename = gamefilename
        self.components = dict()
        for c in compdata['components']:
            self.components[c['name']] = c
        self.objects = dict()
        for o in gamedata['objects']:
            self.objects[o['name']] = o

    def addObject(self, objname):
        if objname and objname not in self.objects:
            obj = dict()
            obj['name'] = objname
            obj['components'] = dict()
            self.objects[objname] = obj
            return True
        else:
            return False

    def setComponentValue(self, objname, compname, vname, vtype, newVal):
        def convert(val, t):
            if t == 'string':
                return val
            elif t == 'int':
                return int(val)
            elif t == 'bool':
                return bool(val)

        obj = self.objects[objname]
        for c in obj['components']:
            if c['type'] == compname:
                assert vname in c['values']
                c['values'][vname] = convert(newVal, vtype)
                return
        assert False, 'Component %s not found in object %s' % (compname, objname)

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
    e = Editor(None)
    e.setModel(model)
    ed.MainLoop()    

if __name__ == '__main__':
    main()

