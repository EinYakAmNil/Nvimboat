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

describe("nvimboat", function()
	it("can show the main page", function()
		vim.cmd.Nvimboat("enable")
		vim.cmd.Nvimboat("show-main")
		eq("MainMenu", nvimboat.pages[1].type)
		eq("", nvimboat.pages[1].id)
	end)
	it("can show tags", function()
		vim.cmd.Nvimboat("show-tags")
		utils.eq_buf(tags_page_buf)
	end)
end)
