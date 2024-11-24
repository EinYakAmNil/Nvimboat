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
	eq(#expected_buf, #rendered)
	for idx, line in ipairs(rendered) do
		eq(expected_buf[idx], line)
	end
end

nvimboat.setup({
	linkHandler = "firefox",
	-- dbPath = os.getenv("HOME") .. "/.cache/nvimboat-test/reload_test-20241123.db",
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
		-- vim.cmd.Nvimboat("reload")
		vim.cmd.Nvimboat("show-main")
		local main_menu_buf = {
			" | N (10/10) | Arch Linux: Recent news updates   | https://www.archlinux.org/feeds/news/",
			" | N (15/15) | CaravanPalace                     | https://www.youtube.com/feeds/videos.xml?user=CaravanPalace",
			" | N (17/17) | Not Related! A Big-Braned Podcast | https://notrelated.xyz/rss",
			" | N (30/30) | Path of Exile News                | https://www.pathofexile.com/news/rss",
			" | N (50/50) | ShortFatOtaku on Odysee           | https://odysee.com/$/rss/@ShortFatOtaku:1",
			" | N (10/10) | Starsector                        | https://fractalsoftworks.com/feed/",
		}
		eq_buf(main_menu_buf)
		eq("MainMenu", nvimboat.pages[1].type)
		eq("", nvimboat.pages[1].id)
	end)
	it("can select a feed", function()
		local url = "https://www.archlinux.org/feeds/news/"
		vim.cmd.Nvimboat("select", url)
		local feed_buf = {
			" | N | 19 Nov 24 | Rafael Epplée          | Providing a license for package sources                                   | https://archlinux.org/news/providing-a-license-for-package-sources/",
			" | N | 14 Sep 24 | Morten Linderud        | Manual intervention for pacman 7.0.0 and local repositories required      | https://archlinux.org/news/manual-intervention-for-pacman-700-and-local-repositories-required/",
			" | N | 01 Jul 24 | Robin Candau           | The sshd service needs to be restarted after upgrading to openssh-9.8p1   | https://archlinux.org/news/the-sshd-service-needs-to-be-restarted-after-upgrading-to-openssh-98p1/",
			" | N | 15 Apr 24 | Christian Heusel       | Arch Linux 2024 Leader Election Results                                   | https://archlinux.org/news/arch-linux-2024-leader-election-results/",
			" | N | 07 Apr 24 | Robin Candau           | Increasing the default vm.max_map_count value                             | https://archlinux.org/news/increasing-the-default-vmmax_map_count-value/",
			" | N | 29 Mar 24 | David Runge            | The xz package has been backdoored                                        | https://archlinux.org/news/the-xz-package-has-been-backdoored/",
			" | N | 04 Mar 24 | Morten Linderud        | mkinitcpio hook migration and early microcode                             | https://archlinux.org/news/mkinitcpio-hook-migration-and-early-microcode/",
			" | N | 09 Jan 24 | Jan Alexander Steffens | Making dbus-broker our default D-Bus daemon                               | https://archlinux.org/news/making-dbus-broker-our-default-d-bus-daemon/",
			" | N | 04 Dec 23 | Christian Heusel       | Bugtracker migration to GitLab completed                                  | https://archlinux.org/news/bugtracker-migration-to-gitlab-completed/",
			" | N | 02 Nov 23 | Frederik Schwan        | Incoming changes in JDK / JRE 21 packages may require manual intervention | https://archlinux.org/news/incoming-changes-in-jdk-jre-21-packages-may-require-manual-intervention/",
			-- " | N | 22 Sep 23 | David Runge            | Changes to default password hashing algorithm and umask settings          | https://archlinux.org/news/changes-to-default-password-hashing-algorithm-and-umask-settings/",
		}
		eq_buf(feed_buf)
		eq("Feed", nvimboat.pages[2].type)
		eq(url, nvimboat.pages[2].id)
	end)
end)
