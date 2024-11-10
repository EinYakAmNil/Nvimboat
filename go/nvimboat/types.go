package nvimboat

import "github.com/neovim/go-client/nvim"

type NvimboatAction func(*Nvimboat, *nvim.Nvim, ...string) error

type Nvimboat struct {
	Nvim   *nvim.Nvim
	Buffer *nvim.Buffer
	Window *nvim.Window
	Config map[string]any
}

