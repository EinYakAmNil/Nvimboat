local api = vim.api
local eq = assert.are.equal

local M = {}

---Print current buffer.
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

---Checks Nvimboat buffer lines against the expected buffer lines.
---@param expected_buf string[] expeced lines of the buffer
---@return nil
function M.eq_buf(expected_buf)
	local rendered = api.nvim_buf_get_lines(0, 0, -1, false)
	eq(#expected_buf, #rendered)
	for idx, line in ipairs(rendered) do
		eq(expected_buf[idx], line)
	end
end

---Checks if Nvimboat has the cursor in the expected row.
---@param expected_row integer expeced row of the cursor
---@return nil
function M.eq_cursor_row(expected_row)
	cursor = api.nvim_win_get_cursor(0)
	eq(expected_row, cursor[1])
end

return M
