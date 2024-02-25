local page = require("nvimboat.page")
local utils = require("nvimboat.utils")
local api = vim.api
local Nvimboat = vim.cmd.Nvimboat
local M = {}

function M.setup(opts)
	M.separator = opts.separator
	M.reloader = opts.reloader
	M.feeds = opts.feeds
	M.linkhandler = opts.linkhandler
end

function M.select()
	if page.page_type == "TagsPage" then
		local tag = utils.seek_tag()
		Nvimboat("select", tag)
		return
	end
	local id = utils.line_id(M.separator)
	if page.page_type == "Article" then
		vim.fn.jobstart({ M.linkhandler, id }, { detach = true })
		return
	end
	Nvimboat("select", id)
end

function M.back()
	Nvimboat("back")
end

function M.show_main_menu()
	Nvimboat("show-main")
end

function M.show_tags()
	if page.page_type ~= "TagsPage" then
		Nvimboat("show-tags")
	end
end

function M.next_unread()
	if page.page_type == "TagsPage" then
		return
	end
	if page.page_type == "Article" then
		Nvimboat("next-unread")
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
		Nvimboat("prev-unread")
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
	local vim_mode = vim.fn.mode()
	local curpos = api.nvim_win_get_cursor(0)

	if vim_mode == 'n' then
		if page.page_type == "Article" then
			local url = utils.article_url()
			-- It doesn't actually matter what we pass as the second argument
			Nvimboat("toggle-read", url)
			return
		end
		local id = utils.line_id(M.separator)
		Nvimboat("toggle-read", id)
	elseif vim_mode == 'v' or vim_mode == 'V' then
		local ids = utils.seek_ids_visual(M.separator)
		local escape = api.nvim_replace_termcodes("<Esc>", true, false, true)
		api.nvim_feedkeys(escape, "v", false)
		Nvimboat("toggle-read", unpack(ids))
	end
	api.nvim_win_set_cursor(0, curpos)
end

function M.open_media()
	if page.page_type == "Article" then
		local url = utils.article_url()
		vim.fn.jobstart({ M.linkhandler, url }, { detach = true })
		return
	end

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
	if page.page_type ~= "Article" then
		return
	end
	Nvimboat("next-article")
end

function M.prev_article()
	if page.page_type ~= "Article" then
		return
	end
	Nvimboat("prev-article")
end

function M.reload_all()
	local default_reload, special_reload = utils.sort_by_reloader(M.feeds)
	utils.reload_feed(M.reloader, default_reload)
	for reloader, feeds in pairs(special_reload) do
		utils.reload_feed(reloader, feeds)
	end
end

function M.delete()
	local vim_mode = vim.fn.mode()
	local curpos = api.nvim_win_get_cursor(0)

	if vim_mode == 'n' then
		if page.page_type == "Article" then
			local url = utils.article_url()
			-- It doesn't actually matter what we pass as the second argument
			Nvimboat("delete", url)
			return
		end
		local id = utils.line_id(M.separator)
		local max_lines = api.nvim_buf_line_count(0)
		if max_lines  <= curpos then
			api.nvim_win_set_cursor(0, { max_lines, 0 })
		end
		Nvimboat("delete", id)
	elseif vim_mode == 'v' or vim_mode == 'V' then
		local ids = utils.seek_ids_visual(M.separator)
		local escape = api.nvim_replace_termcodes("<Esc>", true, false, true)
		api.nvim_feedkeys(escape, "v", false)
		local max_lines = api.nvim_buf_line_count(0)
		if max_lines  <= curpos then
			api.nvim_win_set_cursor(0, { max_lines, 0 })
		end
		Nvimboat("delete", unpack(ids))
	end
	api.nvim_win_set_cursor(0, curpos)
end

return M
