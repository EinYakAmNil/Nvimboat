local M = {}

---State of Nvimboat
M._enabled = false

M.config = require("nvimboat.config")
M.actions = require("nvimboat.actions")
M.pages = require("nvimboat.pages")
M.keymaps = require("nvimboat.keymaps")

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

---@return nil
function M.enable()
	if M._enabled then
		return
	end

	M.keymaps.activate_keymaps(M.keymaps.keymaps)
	if #M.pages < 1 then
		M.actions.show_main_page()
	end
	M._enabled = true
end

---@return nil
function M.disable()
end

function M.setup(opts)
	opts = opts or {}
	M.config.engine = opts.engine or M.config.engine
	M.config.linkHandler = opts.linkHandler or M.config.linkHandler
	M.config.logPath = opts.logPath or M.config.logPath
	M.config.cachePath = opts.cachePath or M.config.cachePath
	M.config.cacheTime = opts.cacheTime or M.config.cacheTime
	M.config.dbPath = opts.dbPath or M.config.dbPath
	M.config.separator = opts.separator or M.config.separator
	M.feeds = opts.feeds or {}
	M.filters = opts.filters or {}
	M.keymaps.configure(opts.keymaps or {})

	vim.fn["remote#host#Register"]("nvimboat", 'x', start_engine)
	vim.fn["remote#host#RegisterPlugin"]("nvimboat", '0', { {
		type = 'command',
		name = 'Nvimboat',
		sync = 1,
		opts = { complete = "customlist,CompleteNvimboat", nargs = "+" }
	}, {
		type = 'function', name = 'CompleteNvimboat', sync = 1, opts = { _ = "" }
	} })
end

return M
