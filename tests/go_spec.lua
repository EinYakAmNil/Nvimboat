local nvimboat = require("nvimboat")

nvimboat.setup()

describe("nvimboat engine", function()
	it("can handle commands", function()
		vim.cmd.Nvimboat("enable")
		vim.cmd.Nvimboat("disable")
	end)
end)
