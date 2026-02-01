local M = {}

---Neovim echo wrapper.
---@param msg string
---@return nil
function M.echo(msg)
	vim.api.nvim_echo(
		{ { msg } },
		true,
		{}
	)
end

---Split `input_str` by `separator`. If no separator is given, then whitespace is used.
---@param input_str string
---@param separator string
---@return string[]
local function string_split(input_str, separator)
	separator = separator or "%s"

	---@type string[]
	local splits = {}
	for str in input_str:gmatch("(.-)" .. separator) do
		splits[#splits + 1] = str:match("^%s*(.-)%s*$")
	end

	local last_col = input_str:match(".*" .. separator .. "(.*)") or input_str
	if input_str ~= "" then
		splits[#splits + 1] = last_col:match("^%s*(.-)%s*$")
	end

	return splits
end

---@param line string
---@param separator string
---@return string
local function get_last_column(line, separator)
	local splits = string_split(line, separator)
	return splits[#splits]
end

---@param line string
---@param separator string
---@return string
local function get_filter_name(line, separator)
	local splits = string_split(line, separator)
	return splits[2]
end

---Return the id by page type.
---@param line string
---@param page Page
---@param separator string
---@return string
function M.get_select_id(line, page, separator)
	local selection_method = {
		Article = function()
			if page.type == "Article" then
				return page.id
			end
			error(page.type .. " != Article page type", 0)
		end,
		Feed = function()
			return get_last_column(line, separator)
		end,
		Filter = function()
			return get_filter_name(line, separator)
		end,
		MainMenu = function()
			return get_last_column(line, separator)
		end,
		TagFeeds = function()
			return get_last_column(line, separator)
		end,
		TagsOverview = function()
			---@type string
			local id = line:gsub(" %(%d+%)$", "")
			return id
		end,
	}
	return selection_method[page.type]()
end

---Return the ids by page type during visual mode.
---@param page Page
---@param separator string
---@return string[]
function M.multi_select_id(page, separator)
	---@type string[]
	local ids = {}
	---@type integer
	local start_row = vim.fn.getpos("v")[2]
	---@type integer
	local end_row = vim.fn.getpos(".")[2]
	---@type integer
	local direction
	if start_row < end_row then
		direction = 1
	else
		direction = -1
	end
	for row_num = start_row, end_row, direction do
		---@type string
		local line = vim.api.nvim_buf_get_lines(0, row_num - 1, row_num, true)[1]
		ids[#ids + 1] = M.get_select_id(line, page, separator)
	end
	return ids
end

return M
