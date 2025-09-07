package nvimboat

import (
	"github.com/neovim/go-client/nvim"
)

var (
	NbNvim   *nvim.Nvim
	NbBuffer *nvim.Buffer
	Feeds    []*Feed
)
