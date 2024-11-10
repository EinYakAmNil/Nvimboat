local M = {}
local defaults = require("nvimboat.defaults")

local function start_engine()
	return vim.fn.jobstart({ M.config.engine }, { rpc = true })
end

function M.setup(opts)
	opts = opts or {}
	M.config.engine = opts.go or defaults.go
	vim.fn["remote#host#Register"]("nvimboat", 'x', start_engine)
	vim.fn["remote#host#RegisterPlugin"]("nvimboat", '0', {
		{ type = 'command',  name = 'Nvimboat',         sync = 1, opts = { complete = "customlist,CompleteNvimboat", nargs = "+" } },
		{ type = 'function', name = 'CompleteNvimboat', sync = 1, opts = { _ = "" } },
	})
	M.feeds = opts.feeds or {}
	M.filters = opts.filters or {}
	M.config.linkhandler = opts.linkhandler or defaults.linkhandler
	M.config.log = opts.log or defaults.log
end

M.actions = require("nvimboat.actions")
M.config = {}

return M
