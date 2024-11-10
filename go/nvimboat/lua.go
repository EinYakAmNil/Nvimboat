package nvimboat

const (
	luaPackage      = "return package.loaded.nvimboat"
	nvimboatEnable  = luaPackage + ".actions.enable()"
	nvimboatDisable = luaPackage + ".actions.disable()"
	nvimboatConfig  = luaPackage + ".config"
)
