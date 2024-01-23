local default = require("nvimboat.default")
local M = {}

M.feeds = {}
M.reloader = default.reloader

function M.setup(opts)
	M.feeds = opts.feeds or {}
	M.reloader = opts.reloader or default.reloader
end

return M
