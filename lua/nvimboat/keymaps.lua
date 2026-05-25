local actions = require("nvimboat.actions")

local M = {}

---@class RHS
---@field rhs function|string
---@field opts vim.keymap.set.Opts

---@class Keymap
---@field [string] RHS

---@class Keys
---@field [string] Keymap

---@type vim.keymap.set.Opts
local keymap_opts = { buffer = 0, nowait = true, silent = true }

---@type RHS
local copy = { rhs = actions.copy_link, opts = keymap_opts }

---@type RHS
local show_main_menu = { rhs = actions.show_main_page, opts = keymap_opts }

---@type RHS
local toggle_read = { rhs = actions.toggle_read, opts = keymap_opts }

---@type RHS
local select = { rhs = actions.select, opts = keymap_opts }

---@type RHS
local open = { rhs = actions.open, opts = keymap_opts }

---@type RHS
local back = { rhs = actions.back, opts = keymap_opts }

---@type RHS
local delete = { rhs = actions.delete, opts = keymap_opts }

---@type RHS
local next_article = { rhs = actions.next_article, opts = keymap_opts }

---@type RHS
local prev_article = { rhs = actions.prev_article, opts = keymap_opts }

---@type RHS
local next_unread = { rhs = actions.next_unread, opts = keymap_opts }

---@type RHS
local prev_unread = { rhs = actions.prev_unread, opts = keymap_opts }

---@type RHS
local reload = { rhs = actions.reload, opts = keymap_opts }

---@type RHS
local show_tags = { rhs = actions.show_tags, opts = keymap_opts }

---@type Keys
M.keymaps = {
	n = {
		a = toggle_read,
		h = back,
		l = select,
		o = open,
		q = show_main_menu,
		y = copy,
		D = delete,
		J = next_article,
		K = prev_article,
		n = next_unread,
		N = prev_unread,
		p = prev_unread,
		R = reload,
		t = show_tags,
	},
	v = {
		a = toggle_read,
		D = delete,
		o = open,
	}
}

M._buffer_overlap = {}

M._global_overlap = {}


---@param mode string
---@param lhs string
local function save_overlap(mode, lhs)
	if not M._buffer_overlap[mode] then
		M._buffer_overlap[mode] = {}
	end
	if not M._global_overlap[mode] then
		M._global_overlap[mode] = {}
	end
	for _, keymap in ipairs(vim.api.nvim_buf_get_keymap(0, mode)) do
		if keymap.lhs == lhs then
			M._buffer_overlap[mode][lhs] = {
				opts = {
					buffer = 0,
					expr = keymap.expr,
					noremap = keymap.noremap,
					nowait = keymap.nowait,
					script = keymap.script,
					silent = keymap.silent,
				}
			}
			if keymap.rhs then
				M._buffer_overlap[mode][lhs]["rhs"] = keymap.rhs
			elseif keymap.callback then
				M._buffer_overlap[mode][lhs]["callback"] = keymap.callback
			end
		end
	end
	for _, keymap in ipairs(vim.api.nvim_get_keymap(mode)) do
		if keymap.lhs == lhs then
			M._buffer_overlap[mode][lhs] = {
				opts = {
					expr = keymap.expr,
					noremap = keymap.noremap,
					nowait = keymap.nowait,
					script = keymap.script,
					silent = keymap.silent,
				}
			}
			if keymap.rhs then
				M._buffer_overlap[mode][lhs]["rhs"] = keymap.rhs
			elseif keymap.callback then
				M._buffer_overlap[mode][lhs]["callback"] = keymap.callback
			end
		end
	end
end

---@param keymaps Keys
---@return Keys
function M.configure(keymaps)
	if not keymaps then
		return M.keymaps
	end

	for mode, keymap in pairs(keymaps) do
		---@cast mode string
		---@cast keymap Keymap
		for lhs, map in pairs(keymap) do
			---@cast lhs string
			---@cast map RHS
			M.keymaps[mode][lhs] = { rhs = map.rhs, opts = map.opts }
		end
	end
	return M.keymaps
end

---@param keymaps Keys
function M.activate_keymaps(keymaps)
	for mode, keymap in pairs(keymaps) do
		---@cast mode string
		---@cast keymap Keymap
		for lhs, map in pairs(keymap) do
			---@cast lhs string
			---@cast map RHS
			save_overlap(mode, lhs)
			vim.keymap.set(mode, lhs, map.rhs, map.opts or {})
		end
	end
end

return M
