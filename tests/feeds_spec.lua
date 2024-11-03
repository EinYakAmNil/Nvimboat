local feeds = {
	{ rssurl = "https://notrelated.xyz/rss",                                  tags = { "Science", "Dead" } },
	{ rssurl = "https://www.archlinux.org/feeds/news/",                       tags = { "Tech", "Linux" }, },
	{ rssurl = "https://www.pathofexile.com/news/rss",                        tags = { "Gaming", "Path of Exile" } },
	{ rssurl = "https://fractalsoftworks.com/feed/",                          tags = { "Gaming", "Starsector" } },
	{ rssurl = "https://odysee.com/$/rss/@ShortFatOtaku:1",                   tags = { "Video", "Odysee", "Politics" } },
	{ rssurl = "https://www.youtube.com/feeds/videos.xml?user=CaravanPalace", tags = { "Video", "YouTube", "Music" } },
}
local nvimboat = require("nvimboat")
nvimboat.setup({
	feeds = feeds
})

describe("the database", function ()
	it("can initialize the database correctly", function ()
	end)
end)
