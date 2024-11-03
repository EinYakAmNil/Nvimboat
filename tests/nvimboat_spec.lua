local nvimboat = require("nvimboat")
local eq = assert.are.equal

nvimboat.setup({
	linkhandler = "firefox"
})

describe("nvimboat", function()
	it("can be configured", function()
		eq("/home/linkai/Projekte/Nvimboat/go/", nvimboat.go)
		eq("firefox", nvimboat.linkhandler)
	end)
end)
