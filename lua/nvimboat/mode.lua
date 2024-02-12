local keymaps = require("nvimboat.keymaps")
local page = require("nvimboat.page")
local api = vim.api
local M = {}

M._enabled = false
M._overlap = {
	n = {},
	v = {}
}

local function find_overlap(mode, lhs)
	local buf_local_maps = api.nvim_buf_get_keymap(0, mode)
	for _, prev_map in ipairs(buf_local_maps) do
		if prev_map["lhs"] == lhs then
			if prev_map["rhs"] then
				M._overlap[mode][prev_map["lhs"]] = { rhs = prev_map["rhs"], opts = { buffer = 0 } }
			elseif prev_map["callback"] then
				M._overlap[mode][prev_map["lhs"]] = { callback = prev_map["callback"], opts = { buffer = 0 } }
			end
		end
	end
end

local function activate_keymaps(maps)
	for mode, map in pairs(maps) do
		for lhs, args in pairs(map) do
			find_overlap(mode, lhs)
			local keymap_opts = args["opts"] or {}
			keymap_opts["buffer"] = 0
			vim.keymap.set(mode, lhs, args["rhs"], keymap_opts)
		end
	end
end

local function restore_mappings(nvimboat_maps, orig_maps)
	for mode, map in pairs(nvimboat_maps) do
		for lhs, _ in pairs(map) do
			vim.keymap.del(mode, lhs, { buffer = 0 })
		end
	end

	for mode, map in pairs(orig_maps) do
		for lhs, args in pairs(map) do
			if args.rhs then
				vim.keymap.set(mode, lhs, args.rhs, { buffer = 0 })
			end
			if args.callback then
				vim.keymap.set(mode, lhs, args.callback, { buffer = 0 })
			end
		end
	end

	M._overlap = {
		n = {},
		v = {}
	}
end

function M.enable()
	if M._enabled then
		return
	end

	activate_keymaps(keymaps.keymaps)
	if #page.page_type == "" then
		require("nvimboat.action").show_main_menu()
	end
	M._enabled = true
end

function M.disable()
	if not M._enabled then
		return
	end

	restore_mappings(keymaps.keymaps, M._overlap)
	M._enabled = false
end

return M
