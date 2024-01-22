local page = require("nvimboat.page")
local utils = require("nvimboat.utils")
local api = vim.api
local M = {}

function M.setup(opts)
	M.separator = opts.separator
end

function M.select()
	if page.page_type == "Article" then
		return
	end
	if page.page_type == "TagsPage" then
		local tag = utils.line_tag()
		vim.cmd.Nvimboat("select", tag)
		return
	end
	local url_pages = { "MainMenu", "TagFeeds", "Feed", "Filter" }
	for _, p in ipairs(url_pages) do
		if page.page_type == p then
			local id = utils.line_id(M.separator)
			vim.cmd.Nvimboat("select", id)
			return
		end
	end
end

function M.back()
	vim.cmd.Nvimboat("back")
end

function M.show_main_menu()
	vim.cmd.Nvimboat("show-main")
end

function M.show_tags()
	if page.page_type ~= "TagsPage" then
		vim.cmd.Nvimboat("show-tags")
	end
end

function M.next_unread()
	if page.page_type == "TagsPage" then
		return
	end
	if page.page_type == "Article" then
		vim.cmd.Nvimboat("next-unread")
		return
	end

	local row = api.nvim_win_get_cursor(0)[1]
	local max_lines = #api.nvim_buf_get_lines(0, 0, -1, false)
	local set = {
		unread_feed = true,
		unread_filter = true,
		unread_article = true,
	}

	for i = row + 1, max_lines, 1 do
		local node_type = vim.treesitter.get_node({ pos = { i - 1, 0 } }):type()
		if set[node_type] ~= nil then
			api.nvim_win_set_cursor(0, { i, 0 })
			return
		end
	end
end

function M.prev_unread()
	if page.page_type == "TagsPage" then
		return
	end
	if page.page_type == "Article" then
		vim.cmd.Nvimboat("prev-unread")
		return
	end

	local set = {
		unread_feed = true,
		unread_filter = true,
		unread_article = true,
	}
	local row = api.nvim_win_get_cursor(0)[1]
	for i = row - 1, 1, -1 do
		local node_type = vim.treesitter.get_node({ pos = { i - 1, 0 } }):type()
		if set[node_type] ~= nil then
			api.nvim_win_set_cursor(0, { i, 0 })
			return
		end
	end
end

function M.toggle_article_read()
end

function M.open_media()
	local vim_mode = vim.fn.mode()

	if vim_mode == 'n' then
		local url = utils.line_id(M.separator)
		utils.play_videos({ url })
	elseif vim_mode == 'v' or vim_mode == 'V' then
		local urls = utils.seek_ids_visual(M.separator)
		utils.play_videos(urls)
	end
end

function M.next_article()
end

function M.prev_article()
end

function M.reload()
end

return M
