local M = {}

local function get_plugin_path()
	local path = debug.getinfo(1, "S").source:sub(2)
	path = vim.fn.fnamemodify(path, ":p:h")
	path = string.gsub(path, "lua/nvimboat", "")
	return path
end

M.plugin_path = get_plugin_path()
M.linkhandler = os.getenv("BROWSER")
M.go = get_plugin_path() .. "go/"

return M
