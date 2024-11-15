local M = {}
local defaults = require("nvimboat.defaults")

local function start_engine()
	return vim.fn.jobstart({ M.config.engine }, {
		rpc = true,
		on_stderr = function(_, data)
			local log = io.open(M.config.logPath, "a")
			if log == nil then
				print("cannot open log file")
				return
			end
			for _, line in ipairs(data) do
				log:write(line)
			end
			log:write("\n")
			log:close()
		end,
		stderr_buffered = true
	})
end

function M.setup(opts)
	opts = opts or {}
	M.config.engine = opts.engine or defaults.engine
	M.config.linkHandler = opts.linkHandler or defaults.linkHandler
	M.config.logPath = opts.logPath or defaults.logPath
	M.config.cachePath = opts.cachePath or defaults.cachePath
	M.config.cacheTime = opts.cacheTime or defaults.cacheTime
	M.config.dbPath = opts.dbPath or defaults.dbPath
	M.feeds = opts.feeds or {}
	M.filters = opts.filters or {}

	vim.fn["remote#host#Register"]("nvimboat", 'x', start_engine)
	vim.fn["remote#host#RegisterPlugin"]("nvimboat", '0', {
		{ type = 'command',  name = 'Nvimboat',         sync = 1, opts = { complete = "customlist,CompleteNvimboat", nargs = "+" } },
		{ type = 'function', name = 'CompleteNvimboat', sync = 1, opts = { _ = "" } },
	})
end

M.actions = require("nvimboat.actions")
M.config = {}

return M
