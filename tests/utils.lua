local api = vim.api
local eq = assert.are.equal

local M = {}

---print current buffer
---@return table buf_lines list of lines in buffer.
function M.print_buf()
	print("\ncurrent buffer lines:\n")
	buf_lines = api.nvim_buf_get_lines(0, 0, -1, false)
	for _, l in ipairs(buf_lines) do
		print(l)
	end
	print()
	return buf_lines
end

function M.eq_buf(expected_buf)
	local rendered = api.nvim_buf_get_lines(0, 0, -1, false)
	eq(#expected_buf, #rendered)
	for idx, line in ipairs(rendered) do
		eq(expected_buf[idx], line)
	end
end

function M.eq_cursor_row(expected_row)
	cursor = api.nvim_win_get_cursor(0)
	eq(expected_row, cursor[1])
end

return M
