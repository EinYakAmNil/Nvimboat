local api = vim.api
local eq = assert.are.same
local Nvimboat = vim.cmd.Nvimboat

local function dump_buffer()
	local lines = api.nvim_buf_get_lines(0, 0, -1, false)
	for _, l in ipairs(lines) do
		print(l)
	end
	print()
	print()
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
		Nvimboat("show-main")
		os.execute("cp " .. nvimboat.config.dbpath .. ".orig " .. nvimboat.config.dbpath)
	end)
	it("setup feeds.", function()
		nvimboat.setup({
			feeds = feeds_config,
			filters = filter_config,
			keymaps = {
				n = {
					z = {
						rhs = function()
							print("hello")
						end,
						opts = { silent = true },
					}
				},
			}
		})
		eq(nvimboat.feeds, feeds_config)
		Nvimboat("enable")
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
		Nvimboat("enable")
		eq("MainMenu", nvimboat.page.page_type)
	end)
	it("can select a feed then a article and go back", function()
		Nvimboat("select", "https://lukesmith.xyz/rss.xml")
		eq("Feed", nvimboat.page.page_type)
		Nvimboat("select", "https://lukesmith.xyz/updates/lindypress-bug-fix/")
		eq("Article", nvimboat.page.page_type)
		Nvimboat("show-main")
		eq("MainMenu", nvimboat.page.page_type)
		Nvimboat("select", "https://lukesmith.xyz/rss.xml")
		eq(nvimboat.page.page_type, "Feed")
	end)
	it("can toggle the unread state of articles", function()
		Nvimboat("select", "https://lukesmith.xyz/rss.xml")
		-- dump_buffer()
		-- print()
		Nvimboat("select", "https://lukesmith.xyz/updates/lindypress-bug-fix/")
		Nvimboat("toggle-read", "https://lukesmith.xyz/updates/lindypress-bug-fix/")
		-- dump_buffer()
		Nvimboat("back")
		-- dump_buffer()
	end)
	it("syncs feed -> article -> back -> back correctly with main menu", function()
		Nvimboat("select", "https://lukesmith.xyz/rss.xml")
		-- dump_buffer()
		-- print()
		Nvimboat("select", "https://lukesmith.xyz/articles/blockchain-blasphemy/")
		Nvimboat("back")
		-- dump_buffer()
		eq("Feed", nvimboat.page.page_type)
		Nvimboat("back")
		eq("MainMenu", nvimboat.page.page_type)
	end)
	it("can toggle articles of a feed", function()
		Nvimboat("select", "https://lukesmith.xyz/rss.xml")
		-- dump_buffer()
		-- print()
		Nvimboat("toggle-read", "https://lukesmith.xyz/updates/lindypress-bug-fix/",
			"https://lukesmith.xyz/articles/blockchain-blasphemy/")
		-- dump_buffer()
	end)
	it("syncs filter -> article -> back -> back correctly with main menu", function()
		-- dump_buffer()
		Nvimboat("select", "query: unread = 1, tags: YouTube")
		-- dump_buffer()
		eq("Filter", nvimboat.page.page_type)
		Nvimboat("select", "https://www.youtube.com/watch?v=M3jpevc6rj8")
		-- dump_buffer()
		-- vim.cmd.sleep("1000ms")
		Nvimboat("back")
		-- dump_buffer()
		Nvimboat("back")
		-- dump_buffer()
	end)
	it("can show the tags page", function()
		Nvimboat("show-tags")
		eq("TagsPage", nvimboat.page.page_type)
		-- dump_buffer()
	end)
	it("can select a tag", function()
		Nvimboat("show-tags")
		local id = nvimboat.utils.seek_tag()
		Nvimboat("select", id)
		eq("TagFeeds", nvimboat.page.page_type)
		Nvimboat("select", "https://www.youtube.com/feeds/videos.xml?channel_id=UCBzYAtjNpudOSCP5_8W1iAw")
		eq("Feed", nvimboat.page.page_type)
		-- dump_buffer()
	end)
	it("can show the next article", function()
		Nvimboat("select", "https://lukesmith.xyz/rss.xml")
		-- dump_buffer()
		Nvimboat("select", "https://lukesmith.xyz/updates/lindypress-bug-fix/")
		-- dump_buffer()
		Nvimboat("next-article")
		-- dump_buffer()
		Nvimboat("back")
		Nvimboat("select", "https://lukesmith.xyz/articles/why-do-i-so-rarely-talk-about-politics-on-my-channel/")
		-- dump_buffer()
		Nvimboat("next-article")
		-- dump_buffer()
		Nvimboat("next-article")
		-- dump_buffer()
		Nvimboat("next-article")
	end)
	it("can find the next unread article", function()
		Nvimboat("select", "https://lukesmith.xyz/rss.xml")
		-- dump_buffer()
		Nvimboat("select", "https://lukesmith.xyz/updates/lindypress-bug-fix/")
		Nvimboat("next-unread")
		-- dump_buffer()
		Nvimboat("back")
		Nvimboat("select", "https://lukesmith.xyz/articles/why-do-i-so-rarely-talk-about-politics-on-my-channel/")
		Nvimboat("back")
		-- dump_buffer()
	end)
	it("can toggle articles in filter", function()
		Nvimboat("select", "query: unread = 1, tags: YouTube")
		Nvimboat("toggle-read",
			"https://www.youtube.com/watch?v=07fkEAdKOC8",
			"https://www.youtube.com/watch?v=31WWZKYG5j8"
		)
		Nvimboat("toggle-read",
			"https://www.youtube.com/watch?v=07fkEAdKOC8",
			"https://www.youtube.com/watch?v=31WWZKYG5j8"
		)
		Nvimboat("select", "https://www.youtube.com/watch?v=31WWZKYG5j8")
		-- dump_buffer()
	end)
	it("delete articles", function()
		dump_buffer()
		Nvimboat("select", "https://lukesmith.xyz/rss.xml")
		Nvimboat("select", "https://lukesmith.xyz/updates/lindypress-bug-fix/")
		Nvimboat("delete", "https://lukesmith.xyz/updates/lindypress-bug-fix/")
		vim.cmd.sleep("500ms")
		Nvimboat("show-main")
		dump_buffer()
	end)
end)
