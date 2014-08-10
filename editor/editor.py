#!/usr/bin/env python2

import sys
import os
import json
import copy
from contextlib import contextmanager
import tempfile
import subprocess
import fcntl

import wx

class ObjectType(object):
    GameObject = 0
    Prefab = 1

class TextDialog(wx.Dialog):
    def __init__(self, *args, **kwargs):
        super(TextDialog, self).__init__(*args, **kwargs)
        pnl = wx.Panel(self)
        self.mainbox = wx.BoxSizer(wx.VERTICAL)
        st1 = wx.StaticText(pnl, label=kwargs['title'])
        self.mainbox.Add(st1)
        self.textCtrl = wx.TextCtrl(pnl)
        self.mainbox.Add(self.textCtrl)
        okbutton = wx.Button(pnl, label='OK')
        okbutton.Bind(wx.EVT_BUTTON, self.OnOK)
        self.mainbox.Add(okbutton)
        cancelbutton = wx.Button(pnl, label='Cancel')
        cancelbutton.Bind(wx.EVT_BUTTON, self.OnCancel)
        self.mainbox.Add(cancelbutton)
        pnl.SetSizerAndFit(self.mainbox)

    def OnOK(self, e):
        self.EndModal(wx.ID_OK)
        self.Destroy()

    def OnCancel(self, e):
        self.EndModal(wx.ID_CANCEL)
        self.Destroy()

class RedirectText(object):
    def __init__(self, aWxTextCtrl):
        self.streams = list()
        self.out = aWxTextCtrl

    def addStream(self, p):
        flags = fcntl.fcntl(p.stdout, fcntl.F_GETFL)
        fcntl.fcntl(p.stdout, fcntl.F_SETFL, flags|os.O_NONBLOCK)
        self.streams.append(p)

    def addData(self, string):
        self.out.AppendText(string)

    def update(self):
        leftstreams = list()
        for p in self.streams:
            p.poll()
            if p.returncode is None:
                leftstreams.append(p)

        self.streams = leftstreams
        for p in self.streams:
            try:
                line = p.stdout.readline()
            except IOError:
                pass
            else:
                if line:
                    self.out.AppendText(line)

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
        newItem = fileMenu.Append(wx.ID_NEW, 'New', 'New')
        saveItem = fileMenu.Append(wx.ID_SAVE, 'Save', 'Save')
        fitem = fileMenu.Append(wx.ID_EXIT, 'Quit', 'Quit application')
        menubar.Append(fileMenu, '&File')
        self.SetMenuBar(menubar)

        self.Bind(wx.EVT_MENU, self.OnNew, newItem)
        if self.model:
            self.Bind(wx.EVT_MENU, self.OnSave, saveItem)
        self.Bind(wx.EVT_MENU, self.OnQuit, fitem)

        self.Bind(wx.EVT_IDLE, self.OnIdle)

        self.SetSize((800,600))
        self.SetTitle('Editor')
        self.Centre()
        self.Show(True)

        if not self.model:
            return

        self.panel = wx.Panel(self)
        self.mainbox = wx.BoxSizer(wx.VERTICAL)
        self.editorbox = wx.BoxSizer(wx.HORIZONTAL)
        self.consolebox = wx.BoxSizer(wx.HORIZONTAL)
        self.leftbox = wx.BoxSizer(wx.VERTICAL)
        self.rightbox = wx.BoxSizer(wx.VERTICAL)

        self.objectsbox = wx.BoxSizer(wx.VERTICAL)
        self.modifierbox = wx.BoxSizer(wx.HORIZONTAL)
        self.consolebox = wx.BoxSizer(wx.HORIZONTAL)
        self.simplemodifierbox = wx.BoxSizer(wx.VERTICAL)
        self.advancedmodifierbox = wx.BoxSizer(wx.VERTICAL)

        st1 = wx.StaticText(self.panel, label='New Game Object')

        self.newObjCtrl = wx.TextCtrl(self.panel)

        self.newObjCtrl.Bind(wx.EVT_KEY_UP, self.OnNewObject)

        self.tree = wx.TreeCtrl(self.panel, size=(400,300))
        self.tree.Bind(wx.EVT_TREE_SEL_CHANGED, self.OnActiveTree)

        self.removeObjButton = wx.Button(self.panel, label='Delete selected object')
        self.removeObjButton.Bind(wx.EVT_BUTTON, self.OnRemoveObject)

        self.createPrefabButton = wx.Button(self.panel, label='Create prefab')
        self.createPrefabButton.Bind(wx.EVT_BUTTON, self.OnCreatePrefab)

        self.compTypeCtrl = wx.ComboBox(self.panel, choices=self.model.getAvailableComponentTypes(), style=wx.CB_READONLY)
        self.compTypeCtrl.SetSelection(0)

        self.addCompCtrl = wx.Button(self.panel, label='Add component to object')
        self.addCompCtrl.Bind(wx.EVT_BUTTON, self.OnAddComponent)

        self.newCompCtrl = wx.Button(self.panel, label='New component')
        self.newCompCtrl.Bind(wx.EVT_BUTTON, self.OnNewComponent)

        self.playButton = wx.Button(self.panel, label='Play')
        self.playButton.Bind(wx.EVT_BUTTON, self.OnPlay)

        self.editButton = wx.Button(self.panel, label='Edit component')
        self.editButton.Bind(wx.EVT_BUTTON, self.OnEdit)

        self.consoleCtrl = wx.TextCtrl(self.panel, style=wx.TE_MULTILINE, size=(790, 80))
        self.consoleCtrl.SetEditable(False)
        self.redirectOutput = RedirectText(self.consoleCtrl)
        self.model.setOutputTarget(self.redirectOutput)

        # objects
        self.objectsbox.Add(st1, flag=wx.RIGHT)
        self.objectsbox.Add(self.newObjCtrl)
        self.objectsbox.Add(self.tree, flag=wx.EXPAND)
        self.leftbox.Add(self.objectsbox)

        # modifiers
        self.simplemodifierbox.Add(self.removeObjButton)
        self.simplemodifierbox.Add(self.compTypeCtrl)
        self.simplemodifierbox.Add(self.addCompCtrl)
        self.simplemodifierbox.Add(self.editButton)
        self.advancedmodifierbox.Add(self.newCompCtrl)
        self.advancedmodifierbox.Add(self.createPrefabButton)
        self.advancedmodifierbox.Add(self.playButton)
        self.modifierbox.Add(self.simplemodifierbox)
        self.modifierbox.Add(self.advancedmodifierbox)
        self.leftbox.Add(self.modifierbox)

        self.editorbox.Add(self.leftbox, flag=wx.EXPAND|wx.ALL)
        self.editorbox.Add(self.rightbox, flag=wx.EXPAND|wx.RIGHT)
        self.mainbox.Add(self.editorbox, flag=wx.EXPAND|wx.ALL)

        self.consolebox.Add(self.consoleCtrl, flag=wx.EXPAND|wx.ALL)
        self.mainbox.Add(self.consolebox, flag=wx.EXPAND|wx.ALL)

        self.panel.SetSizerAndFit(self.mainbox)

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

    def OnAddComponent(self, e):
        self.addComponent(self.compTypeCtrl.GetValue())

    def OnNewComponent(self, e):
        compname = self._getDialogEntry('New Component')
        if compname:
            self.model.newComponent(compname)
            self.updateInterface()
            self.model.editComponent(compname)
            self.updateInterface()

    def updateInterface(self):
        self.model.updateInterface()
        self.compTypeCtrl.SetItems(self.model.getAvailableComponentTypes())

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

    def _getDialogEntry(self, title):
        dlg = TextDialog(None, title=title)
        res = dlg.ShowModal()
        if res == wx.ID_OK:
            val = dlg.textCtrl.GetValue()
        else:
            val = None
        dlg.Destroy()
        return val

    def OnNew(self, e):
        gamename = self._getDialogEntry('New Project')
        if gamename:
            model = loadModel(dlg.textCtrl.GetValue())
            ed = Editor(None, model=model)
            self.Destroy()

    def OnSave(self, e):
        self.model.save()

    def OnIdle(self, e):
        self.redirectOutput.update()

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

    def OnPlay(self, e):
        self.model.play()

    def OnEdit(self, e):
        self.model.editComponent(self.compTypeCtrl.GetValue())
        self.updateInterface()

class Model(object):
    try:
        basepath = os.path.abspath(os.environ['PLURALITY_PROJECTPATH'])
    except KeyError:
        basepath = os.path.join(os.environ['HOME'], '.plurality', 'projects')

    def getProjectBasePath(self):
        return os.path.join(Model.basepath, self.gamename)

    def getGameFilePath(self):
        return os.path.join(self.getProjectBasePath(), 'game', 'game.json')

    def getEditorGameFilePath(self):
        return os.path.join(self.getProjectBasePath(), 'game', 'editor_last.json')

    def getInterfaceFilePath(self):
        return os.path.join(self.getProjectBasePath(), 'out.json')

    def getSourceDir(self):
        return os.path.join(self.getProjectBasePath(), 'src', self.gamename)

    def getComponentSourcePath(self, compname):
        return os.path.join(self.getSourceDir(), '%s.go' % compname)

    def createMain(self):
        os.makedirs(self.getSourceDir())
        mainstr = '''package main

import (
    "runtime"
    "plurality"
)

func main() {
    runtime.LockOSThread()
    plurality.Main()
}
'''
        with open(self.getComponentSourcePath('main'), 'w') as f:
            f.write(mainstr)

    def copyResources(self):
        sharepath = os.path.join(Model.basepath, 'share')
        if not os.path.exists(sharepath):
            os.symlink(os.path.abspath('../share'), sharepath)

    def _createNewGame(self):
        os.makedirs(os.path.join(self.getProjectBasePath(), 'game'))
        os.makedirs(os.path.join(self.getProjectBasePath(), 'share'))
        with open(self.gamefilename, 'w') as f:
            f.write(json.dumps({'objects':[], 'prefabs':[]}, indent=4))
        self.createMain()
        self.copyResources()
        succ = self._compileGame()
        if not succ:
            raise RuntimeError('Unable to compile created game')

    def setOutputTarget(self, out):
        self.outputTarget = out

    def __init__(self, gamename):
        self.gamename = gamename
        self.gamefilename = self.getGameFilePath()
        self.editorgamefilename = self.getEditorGameFilePath()
        self.compfilename = self.getInterfaceFilePath()

        if not os.path.exists(self.compfilename):
            self._createNewGame()

        self.updateInterface()

        gamedata = json.loads(open(self.gamefilename, 'r').read())
        self.objects = dict()
        for o in gamedata['objects']:
            self.objects[o['name']] = o

        self.prefabs = dict()
        for o in gamedata.get('prefabs', dict()):
            self.prefabs[o['name']] = o

    def getAvailableComponentTypes(self):
        return [c['name'] for c in self.components.values() if c['name'] != 'TransformComponent']

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

    def getSave(self):
        game = dict()
        game['objects'] = self.objects.values()
        game['prefabs'] = self.prefabs.values()
        return json.dumps(game, indent=4)

    def save(self):
        with open(self.gamefilename, 'w') as f:
            f.write(self.getSave())

    @contextmanager
    def goenv(self):
        plpath = os.path.join(os.getcwd(), '..')
        oldgopath = os.environ['GOPATH']
        os.environ['GOPATH'] = self.getProjectBasePath() + ':' + os.path.abspath(os.path.join(plpath, 'go')) + ':' + oldgopath
        try:
            yield
        finally:
            os.environ['GOPATH'] = oldgopath

    def _compileGame(self):
        with self.goenv():
            try:
                output = subprocess.check_output('cd %s && go install plurality && go install %s && bin/%s -o %s' % \
                        (self.getProjectBasePath(), self.gamename, self.gamename, self.compfilename),
                        stderr=subprocess.STDOUT, shell=True)
            except subprocess.CalledProcessError as e:
                self.outputTarget.addData(e.output)
                return False
            else:
                self.outputTarget.addData(output)
                return True

    def play(self):
        binpath = '%s/bin/%s' % (self.getProjectBasePath(), self.gamename)
        if os.path.exists(binpath):
            os.unlink(binpath)
        with open(self.editorgamefilename, 'w') as f:
            f.write(self.getSave())
        succ = self._compileGame()
        if succ:
            ln = 'cd %s && %s %s' % (self.getProjectBasePath(), binpath, self.editorgamefilename)
            p = subprocess.Popen(ln, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, shell=True)
            self.outputTarget.addStream(p)

    def updateInterface(self):
        compdata = json.loads(open(self.compfilename, 'r').read())
        self.components = dict()
        for c in compdata['components']:
            self.components[c['name']] = c

    def editComponent(self, compname):
        sourcepath = self.getComponentSourcePath(compname)
        if os.path.exists(sourcepath):
            os.system('gvim %s' % sourcepath)

    def _newComponentTemplate(self, compname):
        return '''package main

import (
        "plurality"
)

func (c *%s) Name() string {
        return "%s"
}

func init() {
        plurality.ComponentNameMap["%s"] = func() plurality.Componenter { return &%s{} }
}

/* Game code starts here */

type %s struct {
        plurality.Component
}

func (c *%s) Start() {
}

func (c *%s) Update() {
}

''' % ((compname,) * 7)

    def newComponent(self, compname):
        comppath = self.getComponentSourcePath(compname)
        assert not os.path.exists(comppath)
        with open(comppath, 'w') as f:
            f.write(self._newComponentTemplate(compname))

def loadModel(gamename):
    gamename = str(gamename)
    assert all([str.isalnum(l) for l in gamename]) and '/' not in gamename and ' ' not in gamename
    model = Model(gamename)
    return model

def main():
    ed = wx.App()
    try:
        gamename = sys.argv[1]
    except IndexError:
        e = Editor(None, model=None)
    else:
        m = loadModel(gamename)
        e = Editor(None, model=m)
    ed.MainLoop()    

if __name__ == '__main__':
    main()

