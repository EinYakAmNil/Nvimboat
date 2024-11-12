package nvimboat

import (
	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type Nvimboat struct {
	Nvim    *nvim.Nvim
	Buffer  *nvim.Buffer
	Window  *nvim.Window
	Config  map[string]any
	Feeds   []*Feed
	Filters []map[string]any
}

type Feed struct {
	rssdb.RssFeed
	Tags []string
}

type NvimboatAction func(*Nvimboat, *nvim.Nvim, ...string) error
