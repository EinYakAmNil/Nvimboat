local M = {}
local runtime_path = ""

for _, rtp in ipairs(vim.api.nvim_list_runtime_paths()) do
	if rtp:match("/Nvimboat$") then
		runtime_path = rtp .. "/"
	end
end

if runtime_path == "" then
	runtime_path = "./"
end

M.godir = runtime_path .. "go/"
M.cachedir = runtime_path .. "cache/"
M.cachetime = 600
M.dbpath = M.cachedir .. "cache.db"
M.log = runtime_path .. "nvimboat.log"
M.separator = " | "
M.reloader = runtime_path .. "python/reloader.py -v -t " .. M.cachetime .. " -d " .. M.cachedir
M.linkhandler = os.getenv("BROWSER")

return M
