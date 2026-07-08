local M = {}

---State of Nvimboat
M._enabled = false

M.config = require("nvimboat.config")
M.actions = require("nvimboat.actions")
M.pages = require("nvimboat.pages")
M.keymaps = require("nvimboat.keymaps")
M.utils = require("nvimboat.utils")

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
	vim.bo.filetype = "nvimboat"
	vim.bo.buftype = "nofile"
	M._enabled = true
end

---@return nil
function M.disable()
end

---@param opts nvimboat.Config
function M.setup(opts)
	opts = opts or {}
	for key, value in pairs(opts) do
		if key ~= "keymaps" then
			M.config[key] = value
		end
	end
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
