package nvimboat

const (
	luaPackage  = "return package.loaded.nvimboat"
	luaEnable   = luaPackage + ".actions.enable()"
	luaDisable  = luaPackage + ".actions.disable()"
	luaConfig   = luaPackage + ".config"
	luaFeeds    = luaPackage + ".feeds"
	luaPages    = luaPackage + ".pages"
	luaPushPage = luaPackage + ".pages:push(...)"
	luaPopPage  = luaPackage + ".pages:pop()"
)
