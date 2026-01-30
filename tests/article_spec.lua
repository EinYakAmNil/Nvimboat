local assert = require("luassert")
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
	for _, l in ipairs(go_build.stderr) do
		print(l)
	end
end

describe("nvimboat", function()
	it("can be configured", function()
		eq("firefox", nvimboat.config.linkHandler)
	end)
	it("can show the main page", function()
		vim.cmd.Nvimboat("enable")
		vim.cmd.Nvimboat("show-main")
		eq("MainMenu", nvimboat.pages[1].type)
		eq("", nvimboat.pages[1].id)
	end)
	it("can select a feed", function()
		local url = "https://www.archlinux.org/feeds/news/"
		vim.cmd.Nvimboat("select", url)
		eq("Feed", nvimboat.pages[2].type)
		eq(url, nvimboat.pages[2].id)
	end)
	it("can select an article", function()
		local url = "https://archlinux.org/news/providing-a-license-for-package-sources/"
		vim.cmd.Nvimboat("select", url)
		utils.eq_buf(expected.article_buf[2])
		vim.cmd.Nvimboat("back")
		utils.eq_buf(expected.feed_buf[5])
	end)
	it("can select the next unread article of a feed", function()
		local url = "https://archlinux.org/news/glibc-241-corrupting-discord-installation/"
		vim.cmd.Nvimboat("select", url)
		vim.cmd.Nvimboat("prev-unread")
		vim.cmd.Nvimboat("back")
		utils.eq_buf(expected.feed_buf[6])
		vim.cmd.Nvimboat("show-main")
	end)
	it("can select the next/previous (unread) article of a filter", function()
		vim.cmd.Nvimboat("select", "new Linux articles")
		local url = "https://blog.lilydjwg.me/posts/216867.html"
		utils.eq_buf(expected.filter_buf[4])
		vim.cmd.Nvimboat("select", url)
		vim.cmd.Nvimboat("prev-unread")
		vim.cmd.Nvimboat("back")
		utils.eq_buf(expected.filter_buf[5])
		vim.cmd.Nvimboat("select", url)
		vim.cmd.Nvimboat("next-article")
		vim.cmd.Nvimboat("next-article")
		vim.cmd.Nvimboat("next-article")
		vim.cmd.Nvimboat("next-article")
		vim.cmd.Nvimboat("next-article")
		vim.cmd.Nvimboat("next-article")
		vim.cmd.Nvimboat("next-article")
		vim.cmd.Nvimboat("next-article")
		vim.cmd.Nvimboat("next-article")
		vim.cmd.Nvimboat("next-article")
		vim.cmd.Nvimboat("next-article")
		vim.cmd.Nvimboat("next-article")
		vim.cmd.Nvimboat("next-article")
		vim.cmd.Nvimboat("prev-article")
		vim.cmd.Nvimboat("back")
		utils.eq_cursor_row(19)
	end)
	it("can delete an article", function()
		local url = "https://blog.lilydjwg.me/posts/216867.html"
		vim.cmd.Nvimboat("select", url)
		vim.cmd.Nvimboat("delete", url)
		vim.cmd.Nvimboat("show-main")
		vim.cmd.Nvimboat("select", "new Linux articles")
	end)
	it("can copy the article url", function()
		local copy_url = "https://blog.lilydjwg.me/posts/216896.html"
		vim.cmd.Nvimboat("select", copy_url)
		nvimboat.actions.copy_link()
		assert.are.equal(copy_url, vim.fn.getreg("+"))
	end)
end)

vim.system({ "sqlite3", dbPath, "UPDATE rss_item SET unread = 1, deleted = 0;" })
