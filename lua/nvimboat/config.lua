---@class nvimboat.Config
local M = {}

---@return string
local function get_plugin_path()
	local path = debug.getinfo(1, "S").source:sub(2)
	path = vim.fn.fnamemodify(path, ":p:h")
	path = string.gsub(path, "lua/nvimboat", "")
	return path
end

M.pluginPath = get_plugin_path()
M.engine = M.pluginPath .. "go/engine"
M.linkHandler = os.getenv("BROWSER")
M.cachePath = M.pluginPath .. "cache/"
---Format: https://pkg.go.dev/time#Duration
M.cacheTime = "10m"
M.logPath = M.pluginPath .. "nvimboat.log"
M.dbPath = M.cachePath .. "cache.db"
M.userAgent = "nvimboat/v0.3"
M.separator = " │ "

M.feeds = {}
M.filters = {}
M.keymaps = {}

return M
