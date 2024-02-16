local api = vim.api
local eq = assert.are.same

local function dump_buffer()
	local lines = api.nvim_buf_get_lines(0, 0, -1, false)
	for _, l in ipairs(lines) do
		print(l)
	end
end

local nvimboat = require("nvimboat")
local filter_config = {
	{ name = "New YouTube", query = "unread = 1", tags = { "YouTube" } },
	{ name = "New Music",   query = "unread = 1", tags = { "Music" } },
}
local feeds_config = {
	{ rssurl = "https://lukesmith.xyz/rss.xml",                                                tags = { "Tech", "Linux", "Politics" }, },
	{ rssurl = "https://notrelated.xyz/rss",                                                   tags = { "Science" } },
	{ rssurl = "https://www.pathofexile.com/news/rss",                                         tags = { "Gaming", "Path of Exile" } },
	{ rssurl = "https://fractalsoftworks.com/feed/",                                           tags = { "Gaming", "Starsector" } },
	{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCXe00APH0fgg8XvN3sPh3nQ", tags = { "YouTube", "Media" } },
	{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCBzYAtjNpudOSCP5_8W1iAw", tags = { "YouTube", "Animation" } },
	{ rssurl = "https://www.youtube.com/feeds/videos.xml?user=ajollywanker",                   tags = { "YouTube", "Gaming" } },
	{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCbBlqVT59IUpsxNT9Dg1zIw", tags = { "YouTube", "Music" } },
	{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCs_JLrcQMNMHZgclnYwmAcQ", tags = { "YouTube", "Music" } },
	{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCB2tP2QfRG7hTra0KTOtTBg", tags = { "YouTube", "Music" } },
	{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCd8b39MwklavavmTafRtlVw", tags = { "YouTube", "Music" } },
	{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCZrFPUErPLEKLvQykfPxwLw", tags = { "YouTube", "Music" } },
	{ rssurl = "https://www.youtube.com/feeds/videos.xml?user=an0nymooose",                    tags = { "YouTube", "Animation" } },
	{ rssurl = "https://www.youtube.com/feeds/videos.xml?user=audiomachine1",                  tags = { "YouTube", "Music" } },
	{ rssurl = "https://www.youtube.com/feeds/videos.xml?user=MrBlackReborn",                  tags = { "YouTube", "Animation", "TF2" } },
	{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCvigl2g67gl18hJgFex-3zg", tags = { "YouTube", "Blender", "Art" } },
	{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCurYTpc1LEzC_DZUArseVwg", tags = { "YouTube", "Animation", "TF2" } },
	-- {
	-- 	rssurl = "https://www.pixiv.net/en/users/26040235",
	-- 	tags = { "Pixiv", "Art" },
	-- 	reloader = "pixivboat --cache-dir " .. nvimboat.config.cachedir
	-- },
	-- {
	-- 	rssurl = "https://www.pixiv.net/en/users/17509087",
	-- 	tags = { "Pixiv", "Art" },
	-- 	reloader = "pixivboat --cache-dir " .. nvimboat.config.cachedir
	-- },
}

describe("nvimboat", function()
	after_each(function()
		vim.cmd.Nvimboat("show-main")
		os.execute("cp " .. nvimboat.config.dbpath .. ".orig " .. nvimboat.config.dbpath)
	end)
	it("setup feeds.", function()
		nvimboat.setup({
			feeds = feeds_config,
			filters = filter_config
		})
		eq(nvimboat.feeds, feeds_config)
		vim.cmd.Nvimboat("enable")
	end)
	-- it("can sort feeds by their reloader", function()
	-- 	local pixiv_feeds = {
	-- 		"https://www.pixiv.net/en/users/17509087",
	-- 		"https://www.pixiv.net/en/users/26040235",
	-- 	}
	-- 	local reload_sort = nvimboat.utils.sort_by_reloader
	-- 	local default_feeds, special_feeds = reload_sort(nvimboat.feeds)
	-- 	table.sort(pixiv_feeds)
	-- 	table.sort(special_feeds["pixivboat --cache-dir " .. nvimboat.config.cachedir])
	-- 	eq(#default_feeds, 21)
	-- 	eq(pixiv_feeds, special_feeds["pixivboat --cache-dir " .. nvimboat.config.cachedir])
	-- end)
	it("can call the Nvimboat command", function()
		vim.cmd.Nvimboat("enable")
		eq("MainMenu", nvimboat.page.page_type)
	end)
	it("can select a feed then a article and go back", function()
		vim.cmd.Nvimboat("select", "https://lukesmith.xyz/rss.xml")
		eq("Feed", nvimboat.page.page_type)
		vim.cmd.Nvimboat("select", "https://lukesmith.xyz/updates/lindypress-bug-fix/")
		eq("Article", nvimboat.page.page_type)
		vim.cmd.Nvimboat("show-main")
		eq("MainMenu", nvimboat.page.page_type)
		vim.cmd.Nvimboat("select", "https://lukesmith.xyz/rss.xml")
		eq(nvimboat.page.page_type, "Feed")
	end)
	it("can toggle the unread state of articles", function()
		vim.cmd.Nvimboat("select", "https://lukesmith.xyz/rss.xml")
		-- dump_buffer()
		-- print()
		vim.cmd.Nvimboat("select", "https://lukesmith.xyz/updates/lindypress-bug-fix/")
		vim.cmd.Nvimboat("toggle-read", "https://lukesmith.xyz/updates/lindypress-bug-fix/")
		-- dump_buffer()
		vim.cmd.Nvimboat("back")
		-- dump_buffer()
	end)
	it("syncs feed -> article -> back -> back correctly with main menu", function()
		vim.cmd.Nvimboat("select", "https://lukesmith.xyz/rss.xml")
		-- dump_buffer()
		-- print()
		vim.cmd.Nvimboat("select", "https://lukesmith.xyz/articles/blockchain-blasphemy/")
		vim.cmd.Nvimboat("back")
		-- dump_buffer()
		eq("Feed", nvimboat.page.page_type)
		vim.cmd.Nvimboat("back")
		eq("MainMenu", nvimboat.page.page_type)
	end)
	it("can toggle articles of a feed", function()
		vim.cmd.Nvimboat("select", "https://lukesmith.xyz/rss.xml")
		-- dump_buffer()
		-- print()
		vim.cmd.Nvimboat("toggle-read", "https://lukesmith.xyz/updates/lindypress-bug-fix/",
			"https://lukesmith.xyz/articles/blockchain-blasphemy/")
		-- dump_buffer()
	end)
	it("syncs filter -> article -> back -> back correctly with main menu", function()
		dump_buffer()
		vim.cmd.Nvimboat("select", "query: unread = 1, tags: YouTube")
		-- dump_buffer()
		eq("Filter", nvimboat.page.page_type)
		vim.cmd.Nvimboat("select", "https://www.youtube.com/watch?v=M3jpevc6rj8")
		dump_buffer()
		vim.cmd.Nvimboat("back")
		vim.cmd.Nvimboat("back")
		dump_buffer()
	end)
	it("can show the tags page", function()
		vim.cmd.Nvimboat("show-tags")
		eq("TagsPage", nvimboat.page.page_type)
		-- dump_buffer()
	end)
	it("can select a tag", function()
		vim.cmd.Nvimboat("show-tags")
		local id = nvimboat.utils.seek_tag()
		vim.cmd.Nvimboat("select", id)
		eq("TagFeeds", nvimboat.page.page_type)
		-- dump_buffer()
	end)
end)
