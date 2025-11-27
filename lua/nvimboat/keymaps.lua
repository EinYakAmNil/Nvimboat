local actions = require("nvimboat.actions")

local M = {}

---@class KeymapOpts
local keymap_opts = { silent = true, buffer = 0, nowait = true }

---@class RHS
---@field rhs function
---@field opts KeymapOpts

---@class Keymap
---@field [string] RHS

---@type RHS
local copy = { rhs = actions.copy_link, opts = keymap_opts }

---@type RHS
local quit = { rhs = actions.show_main_page, opts = keymap_opts }

---@type RHS
local toggle_read = { rhs = actions.toggle_read, opts = keymap_opts }

---@class Keys
---@field [string] Keymap

---@type Keys
M.keymaps = {
	n = {
		a = toggle_read,
		q = quit,
		y = copy,
	}
}

---@param keymaps Keys
---@return Keys
function M.set(keymaps)
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

return M
