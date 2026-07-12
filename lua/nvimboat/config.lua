---@class nvimboat.Config
local M = {}

---@return string
local function get_plugin_path()
	local path = debug.getinfo(1, "S").source:sub(2)
	path = vim.fn.fnamemodify(path, ":p:h")
	path = string.gsub(path, "lua/nvimboat", "")
	return path
end

---@type string
M.pluginPath = get_plugin_path()
---@type string
M.engine = M.pluginPath .. "go/engine"
---@type string
M.linkHandler = os.getenv("BROWSER") or "firefox"
---@type string
M.cachePath = M.pluginPath .. "cache/"
---Format: https://pkg.go.dev/time#Duration
---@type string
M.cacheTime = "10m"
---@type string
M.logPath = M.pluginPath .. "nvimboat.log"
---@type string
M.dbPath = M.cachePath .. "cache.db"
---@type string
M.userAgent = "nvimboat/v1.0"
---@type string
M.separator = " │ "

M.feeds = {}
M.filters = {}
M.keymaps = {}

return M
