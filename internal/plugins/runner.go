package plugins

import (
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

func RunAll(ip string, port int, banner string) {
	pluginList, _ := List()

	for _, name := range pluginList {
		path := filepath.Join(PluginDir(), name)
		runPlugin(path, ip, port, banner)
	}
}

func runPlugin(path string, ip string, port int, banner string) {
	L := lua.NewState()
	defer L.Close()

	if err := L.DoFile(path); err != nil {
		return
	}

	err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("scan"),
		NRet:    0,
		Protect: true,
	},
		lua.LString(ip),
		lua.LNumber(port),
		lua.LString(banner),
	)
	if err != nil {
		return
	}
}
