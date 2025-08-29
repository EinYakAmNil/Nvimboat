local nvimboat = require("nvimboat")
local utils = require("tests.utils")

local eq = assert.are.equal
local dbPath = os.getenv("HOME") .. "/.cache/nvimboat-test/lua-test.db"

local go_build = vim.system(
	{ "go", "build", "-C", "go/" },
	{ text = true }
):wait()

if go_build.stderr ~= "" then
	for _, l in ipairs(go_build.stderr) do
		print(l)
	end
end

nvimboat.setup({
	filters = { {
		name = "new Linux articles",
		query = "unread = 1",
		tags = { "Linux" }
	}, {
		name = "new non political videos",
		query = "unread = 1",
		tags = { "Video", "!Politics" }
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

local tags_page_buf = {
	"Chinese (1)",
	"Gaming (2)",
	"Linux (2)",
	"Music (1)",
	"Odysee (1)",
	"Path of Exile (1)",
	"Politics (1)",
	"Science (1)",
	"Starsector (1)",
	"Tech (2)",
	"Video (2)",
	"YouTube (1)",
}

local main_menu = {
	"N (22/22) │ new Linux articles                │ query: unread = 1, tags: Linux",
	"N (15/15) │ new non political videos          │ query: unread = 1, tags: Video, !Politics",
	"N (10/10) │ Arch Linux: Recent news updates   │ https://www.archlinux.org/feeds/news/",
	"N (15/15) │ CaravanPalace                     │ https://www.youtube.com/feeds/videos.xml?user=CaravanPalace",
	"N (17/17) │ Not Related! A Big-Braned Podcast │ https://notrelated.xyz/rss",
	"N (30/30) │ Path of Exile News                │ https://www.pathofexile.com/news/rss",
	"N (50/50) │ ShortFatOtaku on Odysee           │ https://odysee.com/$/rss/@ShortFatOtaku:1",
	"N (10/10) │ Starsector                        │ https://fractalsoftworks.com/feed/",
	"N (12/12) │ 依云's Blog                       │ https://blog.lilydjwg.me/feed",
}

local gaming = {
	"N (30/30) │ Path of Exile News │ https://www.pathofexile.com/news/rss",
	"N (10/10) │ Starsector         │ https://fractalsoftworks.com/feed/",
}

describe("nvimboat", function()
	it("can select the tags page from the main menu", function()
		vim.cmd.Nvimboat("enable")
		vim.cmd.Nvimboat("show-main")
		eq("MainMenu", nvimboat.pages[1].type)
		eq("", nvimboat.pages[1].id)
	end)
	it("can show tags", function()
		vim.cmd.Nvimboat("show-tags")
		utils.eq_buf(tags_page_buf)
	end)
	it("can go back", function()
		vim.cmd.Nvimboat("back")
		utils.eq_buf(main_menu)
	end)
	it("can show select the tags page from another page", function()
		local url = "https://www.archlinux.org/feeds/news/"
		vim.cmd.Nvimboat("select", url)
		vim.cmd.Nvimboat("show-tags")
	end)
	it("can select a tag", function()
		local tag = "Gaming"
		vim.cmd.Nvimboat("select", tag)
		utils.eq_buf(gaming)
	end)
	it("can select a feed", function ()
		-- local url = "https://www.pathofexile.com/news/rss"
		local url = "https://fractalsoftworks.com/feed/"
		vim.cmd.Nvimboat("select", url)
		print(vim.inspect(nvimboat.pages))
		vim.cmd.Nvimboat("back")
		utils.print_buf()
		print(vim.inspect(nvimboat.pages))
	end)
end)
