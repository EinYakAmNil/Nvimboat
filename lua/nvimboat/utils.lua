local api = vim.api
local M = {}

function M.seek_id(line, separator)
	local _, url_start = line:reverse():find(separator)

	if url_start then
		local url = line:reverse():sub(1, url_start - #separator):reverse()
		return url
	end
end

function M.seek_ids_visual(separator)
	local start_row = vim.fn.getpos("v")[2]
	local end_row = vim.fn.getcurpos()[2]
	local direction = 0
	if start_row < end_row then
		direction = 1
	else
		direction = -1
	end

	local urls = {}
	for row_num = start_row, end_row, direction do
		local line = api.nvim_buf_get_lines(0, row_num - 1, row_num, true)[1]
		table.insert(urls, M.seek_id(line, separator))
	end

	return urls
end

function M.line_id(separator)
	local row_num = api.nvim_win_get_cursor(0)[1]
	local line = api.nvim_buf_get_lines(0, row_num - 1, row_num, true)[1]
	local url = M.seek_id(line, separator)

	return url or "no url detected"
end

function M.play_videos(urls)
	local playlist_file = "/tmp/nvimboat.playlist"
	local playlist = io.open(playlist_file, "w")
	for _, url in ipairs(urls) do
		if playlist then
			playlist:write(url .. "\n")
		end
	end
	if playlist then
		playlist:close()
	end
	vim.fn.jobstart({ "mpv", "--no-terminal", "--profile=builtin-pseudo-gui", "--playlist=" .. playlist_file },
		{ detach = true })
end

return M
