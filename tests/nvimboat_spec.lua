local api = vim.api
local eq = assert.are.same

describe("nvimboat", function()
	local nvimboat = require("nvimboat")
	after_each(function()
		vim.cmd.sleep("50m")
	end
	)

	it("has loaded the correct default configuration", function()
		nvimboat.setup()
		eq(nvimboat.config.cachedir, "./cache/")
		eq(nvimboat.config.cachetime, 5)
		eq(nvimboat.config.dbpath, "./cache/cache.db")
	end)

	it("can execute the 'Nvimboat' command", function()
		vim.cmd.Nvimboat()
		vim.cmd.Nvimboat("Hello")
		vim.cmd.Nvimboat("show-main")
		vim.cmd.Nvimboat("show-feed", "feed")
		vim.cmd.Nvimboat("show-feed", "feed")
	end)
end)
