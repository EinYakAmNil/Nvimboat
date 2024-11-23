local nvimboat = require("nvimboat")
local eq = assert.are.equal

describe("lua plugin", function()
	it("can push pages", function()
		nvimboat.pages:push("MainMenu", "")
		nvimboat.pages:push("Feed", "feed_url")
		eq(2, #nvimboat.pages)
		nvimboat.pages:pop()
		nvimboat.pages:pop()
		eq(0, #nvimboat.pages)
	end)
end)
