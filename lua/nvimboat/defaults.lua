local M = {}

local function get_plugin_path()
	local path = debug.getinfo(1, "S").source:sub(2)
	path = vim.fn.fnamemodify(path, ":p:h")
	path = string.gsub(path, "lua/nvimboat", "")
	return path
end

M.pluginPath = get_plugin_path()
M.engine = get_plugin_path() .. "go/engine"
M.linkHandler = os.getenv("BROWSER")
M.cachePath = get_plugin_path() .. "cache/"
M.cacheTime = "10m"
M.logPath = get_plugin_path() .. "nvimboat.log"
M.dbPath = M.cachePath .. "cache.db"

return M
