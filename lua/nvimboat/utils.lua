local pages = require("nvimboat.pages")
local M = {}

---Split `input_str` by `separator`. If no separator is given, then whitespace is used
---@param input_str string
---@param separator string
---@return string[]
function M.string_split(input_str, separator)
	separator = separator or "%s"
	local splits = {}
	for str in string.gmatch(input_str, "([^" .. separator .. "]+)") do
		splits[#splits+1] = str:match("^%s*(.-)%s*$")
	end
	return splits
end

---Return the link of the current feed/article into the clipboard
---@return string
function M.get_link()
	assert(vim.api.nvim_get_mode()["mode"] == "n")
	local link = ""
	local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)

	if pages[#pages].type == "Article" then
		for _, buf_lines in lines do
			link = string.match(buf_lines, "^Link: (.*)")
			if link then
				return link
			end
		end
	end
	if pages[#pages].type == "Feed" then
		local cursor_xy = vim.api.nvim_win_get_cursor(0)
		local splits = M.string_split(lines[cursor_xy[1]], " │ ")
		link = splits[#splits]
	end

	return link
end

return M
