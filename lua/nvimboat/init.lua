local M = {}
local default = require("nvimboat.default")
local action = require("nvimboat.action")

M.config = default
M.utils = require("nvimboat.utils")

local function init_nvimboat()
	return vim.fn.jobstart({ M.config.godir .. "gonboat" }, {
		rpc = true,
		on_stderr = function(_, data)
			local msg = ""
			for _, d in ipairs(data) do
				msg = msg .. "\n" .. d
			end
			print(msg)
		end
	})
end

local function load_config(opts)
	local C = {}
	C.runtime_path = opts.runtime_path or default.runtime_path
	C.godir = opts.godir or default.godir
	C.cachedir = opts.cachedir or default.cachedir
	C.cachetime = opts.cachetime or default.cachetime
	C.dbpath = opts.dbpath or default.dbpath
	C.log = opts.log or default.log
	C.separator = opts.separator or default.separator
	C.reloader = opts.reloader or default.reloader
	C.linkhandler = opts.linkhandler or default.linkhandler
	C.feeds = opts.feeds or {}
	action.setup(C)

	return C
end

function M.setup(opts)
	opts = opts or {}
	M.feeds = opts.feeds or {}
	M.filters = opts.filters or {}
	M.config = load_config(opts)
	M.action = require("nvimboat.action")
	M.enable = require("nvimboat.mode").enable
	M.disable = require("nvimboat.mode").disable
	vim.fn["remote#host#Register"]("nvimboat", 'x', init_nvimboat)
	vim.fn["remote#host#RegisterPlugin"]("nvimboat", '0', {
		{ type = 'command',  name = 'Nvimboat',         sync = 1, opts = { complete = "customlist,CompleteNvimboat", nargs = "+" } },
		{ type = 'function', name = 'CompleteNvimboat', sync = 1, opts = { _ = "" } },
	})
	M.page = require("nvimboat.page")
	M.keymaps = require("nvimboat.keymaps").configure(opts)
end

return M
