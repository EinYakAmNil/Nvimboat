local nvimboat = require("nvimboat")
local eq = assert.are.equal

nvimboat.setup({
	linkHandler = "firefox",
	dbPath = os.getenv("HOME") .. "/.cache/nvimboat-test/reload_test.db",
})

describe("nvimboat", function()
	it("can be configured", function()
		eq("firefox", nvimboat.config.linkHandler)
	end)
	it("can show the main page", function()
		vim.cmd.Nvimboat("enable")
		vim.cmd.Nvimboat("show-main")
	end)
end)
