package nvimboat

const (
	luaPackage = "return package.loaded.nvimboat"
	luaEnable  = luaPackage + ".actions.enable()"
	luaDisable = luaPackage + ".actions.disable()"
	luaConfig  = luaPackage + ".config"
)
