# Vimboat

A RSS/Atom/Twitter/Manga feed reader in Neovim

# Components 

## Go

- Handles formatting and display logic
- Makes GET-Requests to fetch the feeds
- Parses feed information for the database

## SQLite

- Store the feed data in the same format as [newsboat](https://newsboat.org/)
- Can be used interchangeably between these two programs

## Lua

- Configures Neovim to function as a feed/article selector
- Creates a special mode with slightly different keymaps
- Colorscheme is based on [treesitter](https://tree-sitter.github.io/tree-sitter/) nodes

## Neovim

- Offers a more versatile UI than newsboat
- Special mode for viewing and managing feeds and articles

# Requirements not declareable in Neovim

- [Go](https://go.dev/)
- [SQLite](https://www.sqlite.org/index.html)

# Installation

## Lazy.nvim
```lua
require("lazy").setup({
    { "EinYakAmNil/Nvimboat" }
})
```

# Configuration

+ Tags are currently unused, but the plan is to use them like in Newsboat

```lua
local nvimboat = require("nvimboat")

nvimboat.setup({
    urls = {
        {
            rssurl = ""https://www.youtube.com/feeds/videos.xml?user=Harry101UK,
            tags = { "YouTube", "Animation" },
        },
        { rssurl = "https://www.archlinux.org/feeds/news/", tags = { "Tech" } },
        { rssurl = "https://suckless.org/atom.xml", tags = { "Tech" } },
    },
    db = 'path/to/database'
    separator = " | " -- separator for UI
    cache_dir = "path/to/xml/cache"
    cache_time = 600 -- time for which cache is valid
})
```

# TODO

- Tags for feeds
- Speed up reload
- Make filters based on tags
- single feed reload

# Usage

- To start use the command: **Nvimboat enable**
- When in Nvimboat mode use **h, j, k, l** to select items
