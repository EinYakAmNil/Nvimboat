# Nvimboat

A RSS/Atom/Twitter/Manga feed reader in Neovim.
It aims to be fully compatible with the database schema of [newsboat](https://newsboat.org/), so migration can be done seamlessly.

# But why would I want to read my RSS-Feeds in my text-editor?

1. Extending newsboat is a pain. Everything has to be set to 'browser' before it can be executed.
2. Keymaps cannot be bound to custom functions. Only Macros can call user defined scripts and I don't want to press an extra key just to invoke that.
3. There is no visual mode in newsboat.
4. Vim movements are more comfortable.
5. Because you can.

# Components 

## Neovim

- I'm using the latest stable version (whatever that may be at the time of you reading this). So development will always be done on that.
- Offers a more versatile UI than newsboat.
- Special mode for viewing and managing feeds and articles.

## [Go](https://go.dev/)

- Backend of the plugin
- Handles formatting and display logic
- Interacts with the database to fetch/update the required information

## [SQLite](https://www.sqlite.org/index.html)

- Store the feed data in the same format as [newsboat](https://newsboat.org/)
- Can be used interchangeably between these two programs

## Lua

- Configures Neovim to function as a feed/article selector
- Creates a special mode with slightly different keymaps
- Colorscheme is based on [treesitter](https://tree-sitter.github.io/tree-sitter/) nodes

## Python

- Reload scripts for different kinds of feeds.
- Default installation comes:
    - with a normal RSS/Atom feed parser
    - Twitter-scraper (relies on nitter instances)
    - Mangapill tracker

# Installation

## Lazy.nvim
```lua
require("lazy").setup({
    { "EinYakAmNil/Nvimboat" }
})
```
## Default values
```lua
nvimboat.godir = runtime_path .. "go/"
nvimboat.cachedir = runtime_path .. "cache/"
nvimboat.cachetime = 600
nvimboat.dbpath = M.cachedir .. "cache.db"
nvimboat.log = runtime_path .. "nvimboat.log"
nvimboat.separator = " | "
nvimboat.reloader = runtime_path .. "python/reloader.py"
```
# Configuration

- Feeds can be tagged to put them into categories and mark them for filters
- A feed needs to have all the tags defined in a filter to be shown
- Putting an exclamation mark in front of a tag can be used to exclude any feed that has been tagged by that

```lua
local nvimboat = require("nvimboat")

nvimboat.setup({
    urls = {
        {
            rssurl = "https://www.youtube.com/feeds/videos.xml?user=Harry101UK",
            tags = { "YouTube", "Animation" },
        },
        { rssurl = "https://www.archlinux.org/feeds/news/", tags = { "Tech" } },
        { rssurl = "https://suckless.org/atom.xml", tags = { "Tech" } },
    },
    filters = {
        {
           name = "New YouTube tech videos, but not music",
           query = "unread = 1",
           tags = { "YouTube", "Tech", "!Music" },
        },
        {
           name = "New Music",
           query = "unread = 1",
           tags = { "Music" },
        },
    db = 'path/to/database'
    separator = " | " -- separator for UI
    cache_dir = "path/to/xml/cache"
    cache_time = 600 -- time for which cache is valid
})
```
- If you change the default separator, then treesitter will be broken and some functionalities won't work.

# Usage

- To start use the command: **Nvimboat enable** or use the included *nvimboat.desktop* file.
- When in Nvimboat mode remaps are done for the local buffer. Disabling Nvimboat mode should restore any custom configuration.
- **l** selects an item, while **h** goes back to the last page. The pages are stored in the Go plugin as a sort of stack.
- **n** shows or puts the cursor on the next unread feed/article. **N**/**p** does it for the previous one. **TODO**: implement periodic behaviour, maybe with a ring buffer?
- **t** shows all the tags similar to newsboat and let's you select them to view all feeds of a specific tag.
- While inside an article **J** and **K** can be used to show the next/previous article in the feed/filter. **TODO**: Show first article of next feed when reaching the end of one feed. 
- **o** in normal and visual mode: will attempt to play selected articles when mpv is installed. **TODO**: Make it more general, but maybe a link handler shouldn't be part of this project. 
- **q** goes back to the main menu.
- **R** updates all the feeds. **TODO**: rework state tracking in neovim so individual feed reload can be done.
