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
		print_buf()
	end)
end)
