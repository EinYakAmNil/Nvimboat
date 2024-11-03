package main

import (
	"github.com/EinYakAmNil/Nvimboat/go/engine/commands"
	nvimPlugin "github.com/neovim/go-client/nvim/plugin"
)

func main() {
	nvimPlugin.Main(func(p *nvimPlugin.Plugin) (err error) {
		p.HandleCommand(&nvimPlugin.CommandOptions{Name: "Nvimboat", NArgs: "+", Complete: "customlist, CompleteNvimboat"}, commands.HandleCommand)
		p.HandleFunction(&nvimPlugin.FunctionOptions{Name: "CompleteNvimboat"}, commands.CompleteNvimboat)
		return
	})
}

