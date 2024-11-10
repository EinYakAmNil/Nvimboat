package nvimboat

import "github.com/neovim/go-client/nvim"

type Nvimboat struct {
	Nvim    *nvim.Nvim
	Buffer  *nvim.Buffer
	Window  *nvim.Window
	Config  map[string]any
	Feeds   []map[string]any
	Filters []map[string]any
}

type NvimboatAction func(*Nvimboat, *nvim.Nvim, ...string) error
