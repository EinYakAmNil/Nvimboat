local nvimboat = require("nvimboat")
local utils = require("tests.utils")

local eq = assert.are.equal
local dbPath = os.getenv("HOME") .. "/.cache/nvimboat-test/lua-test.db"

nvimboat.setup({
	filters = { {
		name = "new Linux articles",
		unread = 1,
		tags = { "Linux" }
	}, {
		name = "new non political videos",
		unread = 1,
		tags = { "Video", "!Politics" }
	}, {
		name = "gaming articles",
		tags = { "Gaming" }
	} },
	feeds = {
		{ rssurl = "https://www.youtube.com/feeds/videos.xml?user=CaravanPalace", tags = { "Video", "YouTube", "Music" } },
		{ rssurl = "https://www.archlinux.org/feeds/news/",                       tags = { "Tech", "Linux" }, },
		{ rssurl = "https://www.pathofexile.com/news/rss",                        tags = { "Gaming", "Path of Exile" } },
		{ rssurl = "https://fractalsoftworks.com/feed/",                          tags = { "Gaming", "Starsector" } },
		{ rssurl = "https://blog.lilydjwg.me/feed",                               tags = { "Linux", "Tech", "Chinese" } },
		{ rssurl = "https://odysee.com/$/rss/@ShortFatOtaku:1",                   tags = { "Video", "Politics", "Odysee" } },
		{ rssurl = "https://notrelated.xyz/rss",                                  tags = { "Science" } },
	},
	linkHandler = "firefox",
	dbPath = dbPath
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

local main_menu_buf_0 = {
	"N (40/40) │ gaming articles                   │ filter: tags: Gaming",
	"N (22/22) │ new Linux articles                │ filter: unread: 1, tags: Linux",
	"N (15/15) │ new non political videos          │ filter: unread: 1, tags: Video, !Politics",
	"N (10/10) │ Arch Linux: Recent news updates   │ https://www.archlinux.org/feeds/news/",
	"N (15/15) │ CaravanPalace                     │ https://www.youtube.com/feeds/videos.xml?user=CaravanPalace",
	"N (17/17) │ Not Related! A Big-Braned Podcast │ https://notrelated.xyz/rss",
	"N (30/30) │ Path of Exile News                │ https://www.pathofexile.com/news/rss",
	"N (50/50) │ ShortFatOtaku on Odysee           │ https://odysee.com/$/rss/@ShortFatOtaku:1",
	"N (10/10) │ Starsector                        │ https://fractalsoftworks.com/feed/",
	"N (12/12) │ 依云's Blog                       │ https://blog.lilydjwg.me/feed",
}
local main_menu_buf_1 = {
	"N (40/40) │ gaming articles                   │ filter: tags: Gaming",
	"N (21/21) │ new Linux articles                │ filter: unread: 1, tags: Linux",
	"N (15/15) │ new non political videos          │ filter: unread: 1, tags: Video, !Politics",
	"N (9/10)  │ Arch Linux: Recent news updates   │ https://www.archlinux.org/feeds/news/",
	"N (15/15) │ CaravanPalace                     │ https://www.youtube.com/feeds/videos.xml?user=CaravanPalace",
	"N (17/17) │ Not Related! A Big-Braned Podcast │ https://notrelated.xyz/rss",
	"N (30/30) │ Path of Exile News                │ https://www.pathofexile.com/news/rss",
	"N (50/50) │ ShortFatOtaku on Odysee           │ https://odysee.com/$/rss/@ShortFatOtaku:1",
	"N (10/10) │ Starsector                        │ https://fractalsoftworks.com/feed/",
	"N (12/12) │ 依云's Blog                       │ https://blog.lilydjwg.me/feed",
}
local feed_buf_0 = {
	"N │ 03 Feb 25 │ Frederik Schwan        │ Glibc 2.41 corrupting Discord installation                              │ https://archlinux.org/news/glibc-241-corrupting-discord-installation/",
	"N │ 16 Jan 25 │ Robin Candau           │ Critical rsync security release 3.4.0                                   │ https://archlinux.org/news/critical-rsync-security-release-340/",
	"N │ 19 Nov 24 │ Rafael Epplée          │ Providing a license for package sources                                 │ https://archlinux.org/news/providing-a-license-for-package-sources/",
	"N │ 14 Sep 24 │ Morten Linderud        │ Manual intervention for pacman 7.0.0 and local repositories required    │ https://archlinux.org/news/manual-intervention-for-pacman-700-and-local-repositories-required/",
	"N │ 01 Jul 24 │ Robin Candau           │ The sshd service needs to be restarted after upgrading to openssh-9.8p1 │ https://archlinux.org/news/the-sshd-service-needs-to-be-restarted-after-upgrading-to-openssh-98p1/",
	"N │ 15 Apr 24 │ Christian Heusel       │ Arch Linux 2024 Leader Election Results                                 │ https://archlinux.org/news/arch-linux-2024-leader-election-results/",
	"N │ 07 Apr 24 │ Robin Candau           │ Increasing the default vm.max_map_count value                           │ https://archlinux.org/news/increasing-the-default-vmmax_map_count-value/",
	"N │ 29 Mar 24 │ David Runge            │ The xz package has been backdoored                                      │ https://archlinux.org/news/the-xz-package-has-been-backdoored/",
	"N │ 04 Mar 24 │ Morten Linderud        │ mkinitcpio hook migration and early microcode                           │ https://archlinux.org/news/mkinitcpio-hook-migration-and-early-microcode/",
	"N │ 09 Jan 24 │ Jan Alexander Steffens │ Making dbus-broker our default D-Bus daemon                             │ https://archlinux.org/news/making-dbus-broker-our-default-d-bus-daemon/",
}
local feed_buf_1 = {
	"N │ 03 Feb 25 │ Frederik Schwan        │ Glibc 2.41 corrupting Discord installation                              │ https://archlinux.org/news/glibc-241-corrupting-discord-installation/",
	"  │ 16 Jan 25 │ Robin Candau           │ Critical rsync security release 3.4.0                                   │ https://archlinux.org/news/critical-rsync-security-release-340/",
	"N │ 19 Nov 24 │ Rafael Epplée          │ Providing a license for package sources                                 │ https://archlinux.org/news/providing-a-license-for-package-sources/",
	"N │ 14 Sep 24 │ Morten Linderud        │ Manual intervention for pacman 7.0.0 and local repositories required    │ https://archlinux.org/news/manual-intervention-for-pacman-700-and-local-repositories-required/",
	"N │ 01 Jul 24 │ Robin Candau           │ The sshd service needs to be restarted after upgrading to openssh-9.8p1 │ https://archlinux.org/news/the-sshd-service-needs-to-be-restarted-after-upgrading-to-openssh-98p1/",
	"N │ 15 Apr 24 │ Christian Heusel       │ Arch Linux 2024 Leader Election Results                                 │ https://archlinux.org/news/arch-linux-2024-leader-election-results/",
	"N │ 07 Apr 24 │ Robin Candau           │ Increasing the default vm.max_map_count value                           │ https://archlinux.org/news/increasing-the-default-vmmax_map_count-value/",
	"N │ 29 Mar 24 │ David Runge            │ The xz package has been backdoored                                      │ https://archlinux.org/news/the-xz-package-has-been-backdoored/",
	"N │ 04 Mar 24 │ Morten Linderud        │ mkinitcpio hook migration and early microcode                           │ https://archlinux.org/news/mkinitcpio-hook-migration-and-early-microcode/",
	"N │ 09 Jan 24 │ Jan Alexander Steffens │ Making dbus-broker our default D-Bus daemon                             │ https://archlinux.org/news/making-dbus-broker-our-default-d-bus-daemon/",
}
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
local filter_buf = {
	"N │ 11 Feb 25 │            │ The Legacy of Phrecia FAQ                                                        │ https://www.pathofexile.com/forum/view-thread/3721171",
	"N │ 10 Feb 25 │            │ Legacy of Phrecia Endgame Specialisation System                                  │ https://www.pathofexile.com/forum/view-thread/3720737",
	"N │ 07 Feb 25 │            │ The Legacy of Phrecia Teasers                                                    │ https://www.pathofexile.com/forum/view-thread/3718913",
	"N │ 05 Feb 25 │            │ More Information About the Legacy of Phrecia Event                               │ https://www.pathofexile.com/forum/view-thread/3717515",
	"N │ 04 Feb 25 │            │ Path of Exile 2 - Patch 0.1.1c Preview                                           │ https://www.pathofexile.com/forum/view-thread/3716846",
	"N │ 03 Feb 25 │            │ Update on Path of Exile 1                                                        │ https://www.pathofexile.com/forum/view-thread/3716196",
	"N │ 30 Jan 25 │            │ A Message to Path of Exile 1 Players                                             │ https://www.pathofexile.com/forum/view-thread/3713258",
	"N │ 24 Jan 25 │            │ Weekend Stash Tab Sale                                                           │ https://www.pathofexile.com/forum/view-thread/3707790",
	"N │ 20 Jan 25 │            │ Customer Support Update                                                          │ https://www.pathofexile.com/forum/view-thread/3703227",
	"N │ 16 Jan 25 │            │ Path of Exile 2 - Patch 0.1.1 Patch Note Preview                                 │ https://www.pathofexile.com/forum/view-thread/3695606",
	"N │ 15 Jan 25 │            │ Patch 0.1.1 Q&A VOD                                                              │ https://www.pathofexile.com/forum/view-thread/3694287",
	"N │ 12 Jan 25 │            │ Upcoming Changes in Path of Exile 2 0.1.1                                        │ https://www.pathofexile.com/forum/view-thread/3691520",
	"N │ 10 Jan 25 │            │ Find Out about Patch 0.1.1 on January 12th PST                                   │ https://www.pathofexile.com/forum/view-thread/3687933",
	"N │ 20 Dec 24 │ Alex       │ Anubis-class Cruiser                                                             │ https://fractalsoftworks.com/2024/12/20/anubis-class-cruiser/",
	"N │ 16 Dec 24 │            │ Path of Exile 2: Upcoming Changes and Improvements                               │ https://www.pathofexile.com/forum/view-thread/3642235",
	"N │ 13 Dec 24 │            │ Weekend Stash Tab Sale                                                           │ https://www.pathofexile.com/forum/view-thread/3626462",
	"N │ 11 Dec 24 │            │ Support Email Response Delays                                                    │ https://www.pathofexile.com/forum/view-thread/3616595",
	"N │ 10 Dec 24 │            │ Addressing your Early Access Post-launch Feedback                                │ https://www.pathofexile.com/forum/view-thread/3611705",
	"N │ 06 Dec 24 │            │ Path of Exile 2 Early Access Launch - Live Updates 🔴                            │ https://www.pathofexile.com/forum/view-thread/3594080",
	"N │ 06 Dec 24 │            │ Path of Exile 2 Launches in Early Access Soon - What You Need to Know            │ https://www.pathofexile.com/forum/view-thread/3592995",
	"N │ 04 Dec 24 │            │ Path of Exile 2 - Ascendancy Classes in Early Access                             │ https://www.pathofexile.com/forum/view-thread/3592012",
	"N │ 03 Dec 24 │            │ Path of Exile 2 Early Access Twitch Drops                                        │ https://www.pathofexile.com/forum/view-thread/3591907",
	"N │ 03 Dec 24 │            │ Path of Exile 2 Early Access Pre-Download Information                            │ https://www.pathofexile.com/forum/view-thread/3591631",
	"N │ 22 Nov 24 │            │ Path of Exile 2 Early Access FAQ                                                 │ https://www.pathofexile.com/forum/view-thread/3587981",
	"N │ 21 Nov 24 │            │ Announcing Path of Exile 2 in Early Access                                       │ https://www.pathofexile.com/forum/view-thread/3587754",
	"N │ 20 Nov 24 │            │ GGG Live Twitch Drops                                                            │ https://www.pathofexile.com/forum/view-thread/3587517",
	"N │ 19 Nov 24 │            │ Core Supporter Packs End Soon - An Update About Shipping Physical Items          │ https://www.pathofexile.com/forum/view-thread/3587332",
	"N │ 14 Nov 24 │            │ Watch GGG Live on November 21st - Everything You Need to Know about Early Access │ https://www.pathofexile.com/forum/view-thread/3586621",
	"N │ 14 Nov 24 │            │ Incident Report for Today's Deploy                                               │ https://www.pathofexile.com/forum/view-thread/3586510",
	"N │ 13 Nov 24 │            │ Changes to Path of Exile's Account System                                        │ https://www.pathofexile.com/forum/view-thread/3586288",
	"N │ 04 Nov 24 │            │ Tokyo Game Show ft. Koji Igarashi and Shuhei Yoshida                             │ https://www.pathofexile.com/forum/view-thread/3584989",
	"N │ 13 Jul 24 │ Alex       │ Planet Search Overhaul                                                           │ https://fractalsoftworks.com/2024/07/13/planet-search-overhaul/",
	"N │ 12 Jun 24 │ Stian      │ New music for Galatia Academy                                                    │ https://fractalsoftworks.com/2024/06/12/new-music-for-galatia-academy/",
	"N │ 11 May 24 │ Alex       │ Codex Overhaul                                                                   │ https://fractalsoftworks.com/2024/05/11/codex-overhaul/",
	"N │ 10 Apr 24 │ Alex       │ Save/Load UI, Autosave, Intel Map Markers, and More                              │ https://fractalsoftworks.com/2024/04/10/save-load-ui-autosave-intel-map-markers-and-more/",
	"N │ 13 Mar 24 │ Alex       │ Simulator Enhancements                                                           │ https://fractalsoftworks.com/2024/03/13/simulator-enhancements/",
	"N │ 02 Feb 24 │ Alex       │ Starsector 0.97a Release                                                         │ https://fractalsoftworks.com/2024/02/02/starsector-0-97a-release/",
	"N │ 13 Dec 23 │ Alex       │ Skill Tweaks                                                                     │ https://fractalsoftworks.com/2023/12/12/skill-tweaks/",
	"N │ 24 Nov 23 │ Alex       │ Colony Crises                                                                    │ https://fractalsoftworks.com/2023/11/24/colony-crises/",
	"N │ 13 Nov 23 │ dgbaumgart │ You Merely Adopted Rules.csv, I Was Born Into It                                 │ https://fractalsoftworks.com/2023/11/13/you-merely-adopted-rules-csv-i-was-born-into-it/",
}

describe("nvimboat", function()
	it("can be configured", function()
		eq("firefox", nvimboat.config.linkHandler)
	end)
	it("can show the main page", function()
		vim.cmd.Nvimboat("enable")
		vim.cmd.Nvimboat("show-main")
		utils.eq_buf(main_menu_buf_0)
		eq("MainMenu", nvimboat.pages[1].type)
		eq("", nvimboat.pages[1].id)
	end)
	it("can select a feed", function()
		local url = "https://www.archlinux.org/feeds/news/"
		vim.cmd.Nvimboat("select", url)
		utils.eq_buf(feed_buf_0)
		eq("Feed", nvimboat.pages[2].type)
		eq(url, nvimboat.pages[2].id)
	end)
	it("can select an article", function()
		local url = "https://archlinux.org/news/critical-rsync-security-release-340/"
		vim.cmd.Nvimboat("select", url)
		utils.eq_buf(article_buf)
		eq("Article", nvimboat.pages[3].type)
		eq(url, nvimboat.pages[3].id)
	end)
	it("can go back to the feed with correct cursor position", function()
		vim.cmd.Nvimboat("back")
		eq(2, #nvimboat.pages)
		eq("Feed", nvimboat.pages[#nvimboat.pages].type)
		utils.eq_buf(feed_buf_1)
		utils.eq_cursor_row(2)
	end)
	it("can go back to the main menu with correct cursor position", function()
		vim.cmd.Nvimboat("back")
		eq(1, #nvimboat.pages)
		eq("MainMenu", nvimboat.pages[#nvimboat.pages].type)
		utils.eq_buf(main_menu_buf_1)
		utils.eq_cursor_row(4)
	end)
	it("can select a filter", function()
		vim.cmd.Nvimboat("select", "gaming articles")
		utils.eq_buf(filter_buf)
	end)
end)

vim.system({ "sqlite3", dbPath, "UPDATE rss_item SET unread = 1;" })
