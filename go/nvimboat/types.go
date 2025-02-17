package nvimboat

import (
	"time"

	"github.com/neovim/go-client/nvim"
)

type (
	Nvimboat struct {
		FeedConfig   map[string][]string
		Buffer       *nvim.Buffer
		CachePath    string
		CacheTime    time.Duration
		DbPath       string
		FilterConfig []Filter
		LinkHandler  string
		LogPath      string
		Nvim         *nvim.Nvim
		Pages        []Page
		Window       *nvim.Window
	}
	NvimboatAction func(*Nvimboat, *nvim.Nvim, ...string) error
)
