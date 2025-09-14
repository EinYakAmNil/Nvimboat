package main

import (
	"github.com/EinYakAmNil/Nvimboat/go/engine/nvimboat"
	nvimPlugin "github.com/neovim/go-client/nvim/plugin"
)

func main() {
	nvimPlugin.Main(func(p *nvimPlugin.Plugin) (err error) {
		p.HandleCommand(
			&nvimPlugin.CommandOptions{
				Name:     "Nvimboat",
				NArgs:    "+",
				Complete: "customlist, CompleteNvimboat",
			},
			nvimboat.HandleAction,
		)
		p.HandleFunction(
			&nvimPlugin.FunctionOptions{Name: "CompleteNvimboat"},
			nvimboat.CompleteNvimboat,
		)
		return
	})
}
