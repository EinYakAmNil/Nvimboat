local utils = require "nvimboat.utils"
local M = {}

function M.enable()
	return true
end

function M.disable()
	return true
end

function M.copy_link()
	local link = utils.get_link()
	if link then
		vim.fn.setreg('+', link)
		vim.api.nvim_echo(
			{ { "Copied '" .. link .. "' to the clipboard." } },
			true,
			{}
		)
	else
		vim.api.nvim_echo(
			{ { "No link found." } },
			true,
			{}
		)
	end
end

return M
