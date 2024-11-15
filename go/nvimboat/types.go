package nvimboat

import (
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type (
	Nvimboat struct {
		Nvim        *nvim.Nvim
		Buffer      *nvim.Buffer
		Window      *nvim.Window
		Feeds       []*Feed
		Filters     []map[string]any
		LogPath     string
		CachePath   string
		CacheTime   time.Duration
		DbPath      string
		LinkHandler string
	}

	Feed struct {
		rssdb.RssFeed
		Tags []string
	}

	NvimboatAction func(*Nvimboat, *nvim.Nvim, ...string) error
)
