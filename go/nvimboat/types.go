package nvimboat

import (
	"time"

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
		Pages       PageStack
	}
	PageStack struct {
		Pages []Page
	}
	NvimboatAction func(*Nvimboat, *nvim.Nvim, ...string) error
)
