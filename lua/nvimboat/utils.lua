local api = vim.api
local M = {}

function M.seek_id(line, separator)
	local _, url_start = line:reverse():find(separator)

	if url_start then
		local url = line:reverse():sub(1, url_start - #separator):reverse()
		return url
	end
	return line
end

function M.seek_tag()
	local row_num = api.nvim_win_get_cursor(0)[1]
	local line = api.nvim_buf_get_lines(0, row_num - 1, row_num, true)[1]
	local tag_end, _, _ = line:find("%s%(%d+%)")
	local tag = line:sub(1, tag_end - 1)
	return tag
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
	vim.fn.jobstart(
		"mpv --ytdl-raw-options=force-ipv4= --no-terminal --profile=builtin-pseudo-gui --playlist=" .. playlist_file,
		{ detach = true }
	)
end

function M.reload_feed(reloader, feeds)
	for i, url in ipairs(feeds) do
		feeds[i] = '"' .. url .. '"'
	end
	local urls = table.concat(feeds, " ")
	vim.fn.jobstart(reloader .. " " .. urls, {
		stderr_buffered = true,
		on_stderr = function(_, data)
			if data ~= "" then
				for _, msg in ipairs(data) do
					if msg ~= "" then
						print(msg)
					end
				end
			end
		end,
		on_stdout = function(_, data)
			if data ~= "" then
				for _, msg in ipairs(data) do
					if msg ~= "" then
						print(msg)
					end
				end
			end
		end
	})
end

function M.sort_by_reloader(feeds)
	local default_reload = {}
	local reloaders = {}
	for _, feed in ipairs(feeds) do
		if feed.reloader then
			if reloaders[feed.reloader] then
				table.insert(reloaders[feed.reloader], feed.rssurl)
			else
				reloaders[feed.reloader] = { feed.rssurl }
			end
		else
			table.insert(default_reload, feed.rssurl)
		end
	end
	return default_reload, reloaders
end

function M.article_url()
	local lines = #api.nvim_buf_get_lines(0, 0, -1, false)
	for i = 0, lines, 1 do
		local node_type = vim.treesitter.get_node({ pos = { i, 0 } }):type()
		local line = api.nvim_buf_get_lines(0, i, i + 1, false)[1]
		if node_type == "header" and line:sub(1, 6) == "Link: " then
			return line:gsub("Link: ", "")
		end
	end
end

return M
