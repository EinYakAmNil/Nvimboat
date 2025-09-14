local nvimboat = require("nvimboat")
local utils = require("tests.utils")

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

local toggle_urls = {
	"https://archlinux.org/news/providing-a-license-for-package-sources/",
	"https://archlinux.org/news/manual-intervention-for-pacman-700-and-local-repositories-required/",
	"https://archlinux.org/news/arch-linux-2024-leader-election-results/",
}

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
	it("can toggle the read status", function ()
		vim.cmd.Nvimboat("toggle-read", unpack(toggle_urls))
		utils.print_buf()
	end)
end)

vim.system({ "sqlite3", dbPath, "UPDATE rss_item SET unread = 1;" })
