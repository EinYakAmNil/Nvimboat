local utils = require("nvimboat.utils")

---@alias Action
---| "back"
---| "delete"
---| "disable"
---| "enable"
---| "next-article"
---| "next-unread"
---| "prev-article"
---| "prev-unread"
---| "reload"
---| "select"
---| "show-main"
---| "show-tags"
---| "toggle-read"

---@type fun(action: Action, ...)
local Nvimboat = vim.cmd.Nvimboat

local M = {}

---@return boolean
function M.enable()
	return true
end

---@return boolean
function M.disable()
	return true
end

---@return nil
function M.show_main_page()
	Nvimboat("show-main")
end

---@return nil
function M.toggle_read()
	Nvimboat("toggle-read")
end

---Copy the appropiate link of the current page type to the system clipboard.
---Return true on success
---@return boolean
function M.copy_link()
	local link = utils.get_link()
	if link then
		vim.fn.setreg('+', link)
		vim.api.nvim_echo(
			{ { "Copied '" .. link .. "' to the clipboard." } },
			true,
			{}
		)
		return true
	else
		vim.api.nvim_echo(
			{ { "No link found." } },
			true,
			{}
		)
	end
	return false
end

return M
