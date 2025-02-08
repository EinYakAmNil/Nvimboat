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
	dbPath = os.getenv("HOME") .. "/.cache/nvimboat-test/lua-test.db",
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
			" │ N (10/10) │ Arch Linux: Recent news updates   │ https://www.archlinux.org/feeds/news/",
			" │ N (15/15) │ CaravanPalace                     │ https://www.youtube.com/feeds/videos.xml?user=CaravanPalace",
			" │ N (17/17) │ Not Related! A Big-Braned Podcast │ https://notrelated.xyz/rss",
			" │ N (30/30) │ Path of Exile News                │ https://www.pathofexile.com/news/rss",
			" │ N (50/50) │ ShortFatOtaku on Odysee           │ https://odysee.com/$/rss/@ShortFatOtaku:1",
			" │ N (10/10) │ Starsector                        │ https://fractalsoftworks.com/feed/",
			" │ N (12/12) │ 依云's Blog                       │ https://blog.lilydjwg.me/feed",
		}
		eq_buf(main_menu_buf)
		eq("MainMenu", nvimboat.pages[1].type)
		eq("", nvimboat.pages[1].id)
	end)
	it("can select a feed", function()
		local url = "https://www.archlinux.org/feeds/news/"
		vim.cmd.Nvimboat("select", url)
		local feed_buf = {
			" │ N │ 03 Feb 25 │ Frederik Schwan        │ Glibc 2.41 corrupting Discord installation                              │ https://archlinux.org/news/glibc-241-corrupting-discord-installation/",
			" │ N │ 16 Jan 25 │ Robin Candau           │ Critical rsync security release 3.4.0                                   │ https://archlinux.org/news/critical-rsync-security-release-340/",
			" │ N │ 19 Nov 24 │ Rafael Epplée          │ Providing a license for package sources                                 │ https://archlinux.org/news/providing-a-license-for-package-sources/",
			" │ N │ 14 Sep 24 │ Morten Linderud        │ Manual intervention for pacman 7.0.0 and local repositories required    │ https://archlinux.org/news/manual-intervention-for-pacman-700-and-local-repositories-required/",
			" │ N │ 01 Jul 24 │ Robin Candau           │ The sshd service needs to be restarted after upgrading to openssh-9.8p1 │ https://archlinux.org/news/the-sshd-service-needs-to-be-restarted-after-upgrading-to-openssh-98p1/",
			" │ N │ 15 Apr 24 │ Christian Heusel       │ Arch Linux 2024 Leader Election Results                                 │ https://archlinux.org/news/arch-linux-2024-leader-election-results/",
			" │ N │ 07 Apr 24 │ Robin Candau           │ Increasing the default vm.max_map_count value                           │ https://archlinux.org/news/increasing-the-default-vmmax_map_count-value/",
			" │ N │ 29 Mar 24 │ David Runge            │ The xz package has been backdoored                                      │ https://archlinux.org/news/the-xz-package-has-been-backdoored/",
			" │ N │ 04 Mar 24 │ Morten Linderud        │ mkinitcpio hook migration and early microcode                           │ https://archlinux.org/news/mkinitcpio-hook-migration-and-early-microcode/",
			" │ N │ 09 Jan 24 │ Jan Alexander Steffens │ Making dbus-broker our default D-Bus daemon                             │ https://archlinux.org/news/making-dbus-broker-our-default-d-bus-daemon/",
		}
		eq_buf(feed_buf)
		eq("Feed", nvimboat.pages[2].type)
		eq(url, nvimboat.pages[2].id)
	end)
	it("can select an article", function()
		local url = "https://archlinux.org/news/critical-rsync-security-release-340/"
		vim.cmd.Nvimboat("select", url)
		local article_buf = {
			"Feed: https://www.archlinux.org/feeds/news/",
			"Title: Critical rsync security release 3.4.0",
			"Author: Robin Candau",
			"Date: 16 Jan 25",
			"Link: https://archlinux.org/news/critical-rsync-security-release-340/",
			"== Article Begin ==",
			"We'd like to raise awareness about the rsync security release version `3.4.0-1` as described in our advisory [ASA-202501-1](https://security.archlinux.org/ASA-202501-1).",
			"",
			"An attacker only requires anonymous read access to a vulnerable rsync server, such as a public mirror, to execute arbitrary code on the machine the server is running on.",
			"Additionally, attackers can take control of an affected server and read/write arbitrary files of any connected client.",
			"Sensitive data can be extracted, such as OpenPGP and SSH keys, and malicious code can be executed by overwriting files such as `~/.bashrc` or `~/.popt`.",
			"",
			"We highly advise anyone who runs an rsync daemon or client prior to version `3.4.0-1` to upgrade and reboot their systems immediately.",
			"As Arch Linux mirrors are mostly synchronized using rsync, we highly advise any mirror administrator to act immediately, even though the hosted package files themselves are cryptographically signed.",
			"",
			"All infrastructure servers and mirrors maintained by Arch Linux have already been updated.",
			"",
			"# Links",
			"https://security.archlinux.org/ASA-202501-1",
		}
		eq_buf(article_buf)
		eq("Article", nvimboat.pages[3].type)
		eq(url, nvimboat.pages[3].id)
	end)
end)
