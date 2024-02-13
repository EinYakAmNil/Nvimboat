local api = vim.api
local eq = assert.are.same

local function dump_buffer()
	local lines = api.nvim_buf_get_lines(0, 0, -1, false)
	for _, l in ipairs(lines) do
		print(l)
	end
end

describe("nvimboat", function()
	local nvimboat = require("nvimboat")
	local filter_config = {
		{ name = "New YouTube", query = "unread = 1", tags = { "YouTube" } },
		{ name = "New Music",   query = "unread = 1", tags = { "Music" } },
	}
	local feeds_config = {
		{ rssurl = "https://lukesmith.xyz/rss.xml",                                                tags = { "Tech", "Linux", "Politics" }, },
		{ rssurl = "https://notrelated.xyz/rss",                                                   tags = { "Science" } },
		{ rssurl = "https://www.pathofexile.com/news/rss",                                         tags = { "Gaming", "Path, of, Exile" } },
		{ rssurl = "https://fractalsoftworks.com/feed/",                                           tags = { "Gaming", "Starsector" } },
		{ rssurl = "https://www.spreaker.com/show/3639061/episodes/feed",                          tags = { "Art", "Audiobook" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCXe00APH0fgg8XvN3sPh3nQ", tags = { "YouTube", "Media" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCBzYAtjNpudOSCP5_8W1iAw", tags = { "YouTube", "Animation" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?user=ajollywanker",                   tags = { "YouTube", "Gaming" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCy6Q3wg-PgsgO2XtQxZpZEg", tags = { "YouTube", "Politics" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCbBlqVT59IUpsxNT9Dg1zIw", tags = { "YouTube", "Music" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCs_JLrcQMNMHZgclnYwmAcQ", tags = { "YouTube", "Music" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCB2tP2QfRG7hTra0KTOtTBg", tags = { "YouTube", "Music" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCd8b39MwklavavmTafRtlVw", tags = { "YouTube", "Music" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCe0hvpy81ay5dJvagaeLulA", tags = { "YouTube", "Politics" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCZrFPUErPLEKLvQykfPxwLw", tags = { "YouTube", "Music" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?user=an0nymooose",                    tags = { "YouTube", "Animation" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?user=audiomachine1",                  tags = { "YouTube", "Music" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCUowFWIWGw6Pv2JqfEj8njQ", tags = { "YouTube", "Politics", "Science" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?user=MrBlackReborn",                  tags = { "YouTube", "Animation", "TF2" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCvigl2g67gl18hJgFex-3zg", tags = { "YouTube", "Blender", "Art" } },
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?channel_id=UCurYTpc1LEzC_DZUArseVwg", tags = { "YouTube", "Animation", "TF2" } },
		{
			rssurl = "https://www.pixiv.net/en/users/26040235",
			tags = { "Pixiv", "Art" },
			reloader = "pixivboat --cache-dir " .. nvimboat.config.cachedir
		},
		{
			rssurl = "https://www.pixiv.net/en/users/17509087",
			tags = { "Pixiv", "Art" },
			reloader = "pixivboat --cache-dir " .. nvimboat.config.cachedir
		},
	}
	it("setup feeds.", function()
		nvimboat.setup({
			feeds = feeds_config,
			filters = filter_config
		})
		eq(nvimboat.feeds, feeds_config)
	end)
	it("can sort feeds by their reloader", function()
		local pixiv_feeds = {
			"https://www.pixiv.net/en/users/17509087",
			"https://www.pixiv.net/en/users/26040235",
		}
		local reload_sort = nvimboat.utils.sort_by_reloader
		local default_feeds, special_feeds = reload_sort(nvimboat.feeds)
		table.sort(pixiv_feeds)
		table.sort(special_feeds["pixivboat --cache-dir " .. nvimboat.config.cachedir])
		eq(#default_feeds, 21)
		eq(pixiv_feeds, special_feeds["pixivboat --cache-dir " .. nvimboat.config.cachedir])
	end
	)
	it("can call the Nvimboat command", function()
		-- dump_buffer()
		vim.cmd.Nvimboat("enable")
		-- dump_buffer()
		eq(nvimboat.page.page_type, "MainMenu")
	end)
	it("can select a feed then a article and go back", function()
		vim.cmd.Nvimboat("select", "https://lukesmith.xyz/rss.xml")
		eq(nvimboat.page.page_type, "Feed")
		vim.cmd.Nvimboat("select", "https://lukesmith.xyz/updates/lindypress-bug-fix/")
		eq(nvimboat.page.page_type, "Article")
		dump_buffer()
		vim.cmd.Nvimboat("show-main")
		eq(nvimboat.page.page_type, "MainMenu")
	end)
end)
