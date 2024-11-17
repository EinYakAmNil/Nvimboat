local nvimboat = require("nvimboat")
local eq = assert.are.equal

local function print_buf()
	print("\ncurrent buffer lines:\n")
	buf_lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
	for _, l in ipairs(buf_lines) do
		print(l)
	end
	print()
	return buf_lines
end

local function eq_buf(expected_buf)
	local rendered = vim.api.nvim_buf_get_lines(0, 0, -1, false)
	for idx, line in ipairs(expected_buf) do
		eq(line, rendered[idx])
	end
end

nvimboat.setup({
	linkHandler = "firefox",
	dbPath = os.getenv("HOME") .. "/.cache/nvimboat-test/reload_test.db",
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
		local main_menu_buf = {
			" | N (10/10) | Arch Linux: Recent news updates   | https://www.archlinux.org/feeds/news/",
			" | N (15/15) | CaravanPalace                     | https://www.youtube.com/feeds/videos.xml?user=CaravanPalace",
			" | N (16/17) | Not Related! A Big-Braned Podcast | https://notrelated.xyz/rss",
			" | N (30/30) | Path of Exile News                | https://www.pathofexile.com/news/rss",
			" | N (50/50) | ShortFatOtaku on Odysee           | https://odysee.com/$/rss/@ShortFatOtaku:1",
			" | N (10/10) | Starsector                        | https://fractalsoftworks.com/feed/",
		}
		eq_buf(main_menu_buf)
	end)
end)
