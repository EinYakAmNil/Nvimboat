package nvimboat

import (
	"time"

	"github.com/neovim/go-client/nvim"
)

type (
	Nvimboat struct {
		Buffer       *nvim.Buffer
		CachePath    string
		CacheTime    time.Duration
		DbPath       string
		FeedConfig   map[string][]string
		FilterConfig []*Filter
		LinkHandler  string
		LogPath      string
		Nvim         *nvim.Nvim
		Pages        []Page
		Window       *nvim.Window
	}
	NvimboatAction func(*Nvimboat, *nvim.Nvim, ...string) error
)
