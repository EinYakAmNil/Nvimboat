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
	vim.cmd.Nvimboat("show-tags")
end

function M.toggle_article_read()
end

function M.next_unread()
end

function M.prev_unread()
end

function M.open_media()
end

function M.next_article()
end

function M.prev_article()
end

function M.reload()
end

return M
