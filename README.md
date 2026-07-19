<img src="nvimboat-logo.svg" width="30%" alt="nvimboat-logo">

# Nvimboat

A RSS/Atom feed reader in Neovim.
It aims to be fully compatible with the database schema of [newsboat](https://newsboat.org/), so migration can be done seamlessly.

# But why would I want to read my RSS-Feeds in my text-editor?

1. Extending newsboat is a pain. Everything has to be set to 'browser' before it can be executed.
2. Keymaps cannot be bound to custom functions. Only Macros can call user defined scripts and I don't want to press an extra key just to invoke that.
3. There is no visual mode in newsboat.
4. Vim movements are more comfortable.
5. Because you can.

# Installation

## Lazy.nvim

```lua
return {
    "EinYakAmNil/Nvimboat",
    build = function()
        local nvimboat_path = vim.fn.stdpath("data") .. "/lazy/Nvimboat/"
		vim.system({ "go", "build" },
			{ cwd = nvimboat_path .. "go" })
		vim.system(
			{ "gcc", "-shared", "-o", "../parser/nvimboat.so", "-I./src", "src/parser.c", "-Os" },
			{ cwd = nvimboat_path .. 'treesitter' })
		vim.system({ "cp", "nvimboat-logo.svg", vim.fn.expand("~/.local/share/icons/hicolor/48x48/apps") },
			{ cwd = nvimboat_path })
		vim.system({ "cp", "nvimboat.desktop", vim.fn.expand("~/.local/share/applications") },
			{ cwd = nvimboat_path })
    end,
    cmd = "Nvimboat",
    config = function()
        local nvimboat = require("nvimboat")
        nvimboat.setup({
            feeds = {
                -- see Configuration
            }
        })
    end
}
```

# Requirements

## Neovim

- I'm using the latest stable version (whatever that may be at the time of you reading this). So development will always be done on that.
- Offers a more versatile UI than newsboat.
- Plugin implements a special mode for viewing and managing feeds and articles.

## [Go](https://pkg.go.dev/github.com/neovim/go-client/nvim)

- Engine of the plugin
- Handles formatting and display logic
- Interacts with the database to fetch/update the required information
- Requests RSS feeds

## [SQLite](https://www.sqlite.org/index.html)

- Store the feed data in the same format as [newsboat](https://newsboat.org/)
- Can be used interchangeably between these two programs

## Lua

- Configures Neovim to function as a feed/article selector
- Creates a special mode with slightly different keymaps
- Colorscheme is based on [treesitter](https://tree-sitter.github.io/tree-sitter/) nodes

# Configuration

## LSP

I recommend adding a ```.luarc.json``` file to whereever you configure this plugin.
It should include the Nvimboat installation directory.

```json
{
	"workspace.library": [
		"${3rd}/luassert/library",
		"/usr/share/nvim/runtime/lua/",
		"~/.local/share/nvim/lazy/Nvimboat/lua/nvimboat/"
	],
	"runtime.version": "Lua 5.1"
}
```

## Nvimboat

- Feeds can be tagged to put them into categories and mark them for filters.
- A feed with any matching tags of a filter will be included.
- Putting an exclamation mark in front of a tag can be used to exclude any feed that has been tagged by that.

```lua
local nvimboat = require("nvimboat")

nvimboat.setup({
    feeds = {
        {
            rssurl = "https://www.youtube.com/feeds/videos.xml?user=Harry101UK",
            tags = { "YouTube", "Animation" },
        },
        { rssurl = "https://www.archlinux.org/feeds/news/", tags = { "Tech" } },
        { rssurl = "https://suckless.org/atom.xml", tags = { "Tech" } },
        {
            rssurl = "https://twitter.com/DoctorLalve",
            tags = { "YouTube", "Animation" },
        },
    },
    filters = {
        {
           name = "New YouTube tech videos, but not music",
           unread = 1,
           tags = { "YouTube", "Tech", "!Music" },
        },
        {
           name = "New Music",
           unread = 1,
           tags = { "Music" },
        },
    keymaps = {
        n = { -- keymaps for normal mode
            w = { -- key to be mapped
                rhs = function()
                -- do something
                end,
                opts = {}
            },
        }
        v = { -- keymaps for visual mode
            -- values are merged into default maps
            -- this doesn't remove the preconfigured keymaps for this mode
        }
    }
    -- Default values for the other options
    pluginPath = "~/.local/share/nvim/lazy/Nvimboat" -- Default will be determined dynamically
    logPath = pluginPath .. "nvimboat.log"
    cacheTime = "10m" -- Format: https://pkg.go.dev/time#Duration. Caches HTTP requests for this duration...
    cachePath = pluginPath .. "cache/" -- ... in this directory
    dbPath = cachePath .. "cache.db" -- You should set this to somewhere else, if you don't want it to be lost by deinstalling Nvimboat.
    userAgent = "nvimboat/v1.0"
    separator = " │ " -- separator for UI, changing it will break treesitter
})
```

# Usage

To start use the command: **Nvimboat enable** or use the included *nvimboat.desktop* file.
When in Nvimboat mode remaps are done for the local buffer. Disabling Nvimboat mode should restore any custom configuration.

Keymaps:
- **l** selects an item, while **h** goes back to the last page. The pages are stored in the Go plugin as a sort of stack.
- **n** shows or puts the cursor on the next unread feed/article. **N**/**p** does it for the previous one
- **t** shows all the tags similar to newsboat and let's you select them to view all feeds of a specific tag.
- While inside an article **J** and **K** can be used to show the next/previous article in the feed/filter.
- **o** in normal and visual mode: will attempt to play selected articles when mpv is installed.
- **q** goes back to the main menu.
- **R** updates all the feeds.

# Custom reload scripts (removed, write your own RSS Feed converter instead)

~When reloading, the plugin first sorts the feeds by their reloader an then passes all the URLs of each reloader as a long line of arguments.~
~This should be taken into account when making custom scripts.~

During the rewrite I decided to remove the Python script that would reload the feeds and reimplement it in Go.
This improved database handling as there were bugs with Go and Python (and any other reload scripts) locking each other out of the database.
Maintaining and testing the code also was much easier this way.

My first thought then was to implement the custom reloading logic in Go.
That turned out to be possible, but really inconvenient for the user, because you had to modify the plugin source code and recompile the engine.
This would probably lead to conflicts in plugin updates.

So my solution for my own custom reload scripts was to convert them into webservers which would serve parsable RSS feeds instead.
Most of the logic was there anyway and I only needed to read the [RSS specs](https://www.rssboard.org/rss-specification) to put the values into HTML instead of my database.

# Migration from newsboat

You can just copy your old newsboat database into the cache directory of the installation path.
If you use lazy.nvim it should look something like _$HOME/.local/share/nvim/lazy/Nvimboat/cache_.
