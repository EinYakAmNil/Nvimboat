local M = {}
local defaults = require("nvimboat.defaults")

local function start_engine()
	return vim.fn.jobstart({ M.go .. "/engine" }, {
		rpc = true,
		on_stderr = function(_, data)
			local msg = ""
			for _, line in ipairs(data) do
				msg = msg .. "\n" .. line
			end
		end
	})
end

function M.setup(opts)
	M.go = opts.go or defaults.go
	vim.fn["remote#host#Register"]("nvimboat", 'x', start_engine)
	vim.fn["remote#host#RegisterPlugin"]("nvimboat", '0', {
		{ type = 'command',  name = 'Nvimboat',         sync = 1, opts = { complete = "customlist,CompleteNvimboat", nargs = "+" } },
		{ type = 'function', name = 'CompleteNvimboat', sync = 1, opts = { _ = "" } },
	})
	opts = opts or {}
	M.feeds = opts.feeds or {}
	M.filters = opts.filters or {}
	M.linkhandler = opts.linkhandler or defaults.linkhandler
end

return M
