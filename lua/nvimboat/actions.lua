local utils = require("nvimboat.utils")
local pages = require("nvimboat.pages")
local config = require("nvimboat.config")

---@alias Action
---| "back"
---| "delete"
---| "disable"
---| "enable"
---| "next-article"
---| "next-unread"
---| "prev-article"
---| "prev-unread"
---| "reload"
---| "select"
---| "show-main"
---| "show-tags"
---| "toggle-read"
---| "open"

---@type fun(action: Action, ...: string)
local Nvimboat = vim.cmd.Nvimboat

local M = {}

---Show the main menu.
---Resets the page stack.
---@return nil
function M.show_main_page()
	Nvimboat("show-main")
end

---Toggle the read status of an article or feed.
---@return nil
function M.toggle_read()
	local ids = utils.multi_select_id(pages[#pages], config.separator)
	Nvimboat("toggle-read", unpack(ids))
end

---Return to the previous page.
---@return nil
function M.back()
	Nvimboat("back")
end

---@return string
function M.select()
	---@type integer[]
	local cursor_xy = vim.api.nvim_win_get_cursor(0)
	---@type string[]
	local buf_lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
	local line = buf_lines[cursor_xy[1]]
	local id = utils.get_select_id(line, pages[#pages], config.separator)
	if pages[#pages].type == "Article" then
		vim.system({ config.linkHandler, id }, { detach = true })
		return id
	end
	Nvimboat("select", id)
	return id
end

---Copy the appropiate link of the current page type to the system clipboard.
---Also returns it as a string.
---@return string|nil
function M.copy_link()
	---@type integer[]
	local cursor_xy = vim.api.nvim_win_get_cursor(0)
	---@type string[]
	local buf_lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
	local line = buf_lines[cursor_xy[1]]
	local link = utils.get_select_id(line, pages[#pages], config.separator)
	if link then
		vim.fn.setreg('+', link)
		utils.echo("Copied '" .. link .. "' to the clipboard.")
		return link
	else
		utils.echo("No link found.")
	end
	return nil
end

---Return to the previous page.
---@return string|string[]
function M.delete()
	local mode_map = {
		n = function()
			---@type integer[]
			local cursor_xy = vim.api.nvim_win_get_cursor(0)
			---@type string[]
			local buf_lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
			local line = buf_lines[cursor_xy[1]]
			local id = utils.get_select_id(line, pages[#pages], config.separator)
			Nvimboat("delete", id)
			return id
		end,
		v = function()
			local ids = utils.multi_select_id(pages[#pages], config.separator)
			Nvimboat("delete", unpack(ids))
			return unpack(ids)
		end,
		V = function()
			local ids = utils.multi_select_id(pages[#pages], config.separator)
			Nvimboat("delete", unpack(ids))
			return unpack(ids)
		end,
	}
	return mode_map[vim.fn.mode()]()
end

function M.next_article()
	Nvimboat("next-article")
end

function M.prev_article()
	Nvimboat("prev-article")
end

function M.next_unread()
	if pages[#pages].type == "TagsOverview" then
		return
	end
	if pages[#pages].type == "Article" then
		Nvimboat("next-unread")
	end
	local ts_root = assert(vim.treesitter.get_parser():parse()[1]:root(),
		"tree-sitter parsing failed. `ts_root` is nil")
	local count = ts_root:named_child_count()
	local row, col = unpack(vim.api.nvim_win_get_cursor(0))
	local treesitter_unread_set = {
		unread_feed = true,
		unread_filter = true,
		unread_article = true,
	}
	for i = 0, count - 1 do
		local child = assert(ts_root:named_child((row + i) % count),
			"tree-sitter parsing failed. `ts_root` is nil")
		if treesitter_unread_set[child:type()] then
			vim.api.nvim_win_set_cursor(0, { child:start() + 1, col })
			return
		end
	end
end

function M.prev_unread()
	if pages[#pages].type == "TagsOverview" then
		return
	end
	if pages[#pages].type == "Article" then
		Nvimboat("prev-unread")
	end
	local ts_root = assert(vim.treesitter.get_parser():parse()[1]:root(),
		"tree-sitter parsing failed. `ts_root` is nil")
	local count = ts_root:named_child_count()
	local row, col = unpack(vim.api.nvim_win_get_cursor(0))
	local treesitter_unread_set = {
		unread_feed = true,
		unread_filter = true,
		unread_article = true,
	}
	for i = 0, count - 1 do
		local child = assert(ts_root:named_child((row - i - 2) % count),
			"tree-sitter parsing failed. `ts_root` is nil")
		if treesitter_unread_set[child:type()] then
			vim.api.nvim_win_set_cursor(0, { child:start() + 1, col })
			return
		end
	end
end

function M.reload()
	Nvimboat("reload")
end

function M.show_tags()
	Nvimboat("show-tags")
end

function M.open()
	local mode_map = {
		n = function()
			---@type integer[]
			local cursor_xy = vim.api.nvim_win_get_cursor(0)
			---@type string[]
			local buf_lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
			local line = buf_lines[cursor_xy[1]]
			local id = utils.get_select_id(line, pages[#pages], config.separator)
			Nvimboat("open", id)
			return id
		end,
		v = function()
			local ids = utils.multi_select_id(pages[#pages], config.separator)
			Nvimboat("open", unpack(ids))
			return unpack(ids)
		end,
		V = function()
			local ids = utils.multi_select_id(pages[#pages], config.separator)
			Nvimboat("open", unpack(ids))
			return unpack(ids)
		end,
	}
	return mode_map[vim.fn.mode()]()
end

return M
