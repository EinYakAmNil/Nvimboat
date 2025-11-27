local pages = require("nvimboat.pages")

local M = {}

---Check if item is contained in list
---@param item any
---@param list any[]
---@return boolean
local function contained_in(item, list)
	for _, value in ipairs(list) do
		if item == value then
			return true
		end
	end

	return false
end

---Split `input_str` by `separator`. If no separator is given, then whitespace is used
---@param input_str string
---@param separator string
---@return string[]
function string_split(input_str, separator)
	separator = separator or "%s"
	local splits = {}
	for str in string.gmatch(input_str, "([^" .. separator .. "]+)") do
		splits[#splits + 1] = str:match("^%s*(.-)%s*$")
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

	local allowed_page_types = {
		"Feed",
		"Filter",
		"MainPage",
		"TagFeeds",
	}
	if contained_in(pages[#pages].type, allowed_page_types) then
		local cursor_xy = vim.api.nvim_win_get_cursor(0)
		local splits = string_split(lines[cursor_xy[1]], " │ ")
		link = splits[#splits]
	end

	return link
end

return M
