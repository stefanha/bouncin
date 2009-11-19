package plugins

type Plugin interface {
	Enable() bool;
	Disable();
}

var plugins = make(map[string] Plugin)

func Register(name string, plugin Plugin) {
	plugins[name] = plugin;
}

func Enable(name string) bool {
	plugin, ok := plugins[name];
	if ok {
		ok = plugin.Enable()
	}
	return ok;
}

func Disable(name string) {
	plugin, ok := plugins[name];
	if !ok {
		return
	}
	plugin.Disable()
}
