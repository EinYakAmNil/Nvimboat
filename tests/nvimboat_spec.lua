local nvimboat = require("nvimboat")
local utils = require("tests.utils")
local expected = require("tests.expected")

local eq = assert.are.equal
local dbPath = os.getenv("HOME") .. "/.cache/nvimboat-test/lua-test.db"

nvimboat.setup({
	filters = { {
		name = "new Linux articles",
		unread = 1,
		tags = { "Linux" }
	}, {
		name = "new non political videos",
		unread = 1,
		tags = { "Video", "!Politics" }
	}, {
		name = "gaming articles",
		tags = { "Gaming" }
	} },
	feeds = {
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?user=CaravanPalace", tags = { "Video", "YouTube", "Music" } },
		{ rssurl = "https://www.archlinux.org/feeds/news/",                       tags = { "Tech", "Linux" }, },
		{ rssurl = "https://www.pathofexile.com/news/rss",                        tags = { "Gaming", "Path of Exile" } },
		{ rssurl = "https://fractalsoftworks.com/feed/",                          tags = { "Gaming", "Starsector" } },
		{ rssurl = "https://blog.lilydjwg.me/feed",                               tags = { "Linux", "Tech", "Chinese" } },
		{ rssurl = "https://odysee.com/$/rss/@ShortFatOtaku:1",                   tags = { "Video", "Politics", "Odysee" } },
		{ rssurl = "https://notrelated.xyz/rss",                                  tags = { "Science" } },
	},
	linkHandler = "firefox",
	dbPath = dbPath
})

local go_build = vim.system(
	{ "go", "build", "-C", "go/" },
	{ text = true }
):wait()

if go_build.stderr ~= "" then
	print(go_build.stderr)
end

describe("nvimboat", function()
	it("can be configured", function()
		eq("firefox", nvimboat.config.linkHandler)
	end)
	it("can show the main page", function()
		vim.cmd.Nvimboat("enable")
		vim.cmd.Nvimboat("show-main")
		utils.eq_buf(expected.main_menu_buf[1])
		eq("MainMenu", nvimboat.pages[1].type)
		eq("", nvimboat.pages[1].id)
	end)
	it("can select a feed", function()
		local url = "https://www.archlinux.org/feeds/news/"
		vim.cmd.Nvimboat("select", url)
		utils.eq_buf(expected.feed_buf[1])
		eq("Feed", nvimboat.pages[2].type)
		eq(url, nvimboat.pages[2].id)
	end)
	it("can select an article", function()
		local url = "https://archlinux.org/news/critical-rsync-security-release-340/"
		vim.cmd.Nvimboat("select", url)
		utils.eq_buf(expected.article_buf[1])
		eq("Article", nvimboat.pages[3].type)
		eq(url, nvimboat.pages[3].id)
	end)
	it("can go back to the feed with correct cursor position", function()
		vim.cmd.Nvimboat("back")
		eq(2, #nvimboat.pages)
		eq("Feed", nvimboat.pages[#nvimboat.pages].type)
		utils.eq_buf(expected.feed_buf[2])
		utils.eq_cursor_row(2)
	end)
	it("can go back to the main menu with correct cursor position", function()
		vim.cmd.Nvimboat("back")
		eq(1, #nvimboat.pages)
		eq("MainMenu", nvimboat.pages[#nvimboat.pages].type)
		utils.eq_buf(expected.main_menu_buf[2])
		utils.eq_cursor_row(4)
	end)
	it("can select a filter", function()
		vim.cmd.Nvimboat("select", "gaming articles")
		utils.eq_buf(expected.filter_buf[1])
	end)
	it("can select an article from filter", function()
		vim.cmd.Nvimboat("select", "https://www.pathofexile.com/forum/view-thread/3594080")
	end)
	it("can go back to the filter with correct cursor position", function()
		vim.cmd.Nvimboat("toggle-read", "a")
		utils.eq_cursor_row(19)
	end)
	it("can go back to the main menu with correct cursor position", function()
		vim.cmd.Nvimboat("back")
		utils.eq_cursor_row(3)
	end)
	it("can toggle read of an entire feed", function()
		vim.cmd.Nvimboat("toggle-read", "https://fractalsoftworks.com/feed/")
		utils.eq_buf(expected.main_menu_buf[3])
		vim.cmd.Nvimboat("select", "https://fractalsoftworks.com/feed/")
		utils.eq_buf(expected.feed_buf[3])
	end)
	it("can select the next/prev unread feed/filter", function()
		vim.cmd.Nvimboat("show-main")
		vim.api.nvim_win_set_cursor(0, { 7, 0 })
		vim.cmd.Nvimboat("next-unread")
		utils.eq_cursor_row(8)
		vim.cmd.Nvimboat("next-unread")
		utils.eq_cursor_row(10)
		vim.cmd.Nvimboat("next-unread")
		utils.eq_cursor_row(1)
		vim.cmd.Nvimboat("prev-unread")
		utils.eq_cursor_row(10)
		vim.cmd.Nvimboat("prev-unread")
		vim.cmd.Nvimboat("prev-unread")
		utils.eq_cursor_row(7)
	end)
	it("can delete articles of feeds", function()
		vim.cmd.Nvimboat("show-main")
		vim.cmd.Nvimboat("delete"
		, "https://fractalsoftworks.com/feed/"
		, "https://www.archlinux.org/feeds/news/"
		)
		utils.eq_buf(expected.main_menu_buf[4])
	end)
end)

vim.system({ "sqlite3", dbPath, "UPDATE rss_item SET unread = 1, deleted = 0;" })
