package main

import (
	"nvimboat"

	_ "github.com/neovim/go-client/nvim"
	nvimPlugin "github.com/neovim/go-client/nvim/plugin"
)

func main() {
	chanNvimboat := make(chan *nvimboat.Nvimboat)
	go nvimboatLoop(chanNvimboat)
	nb := <-chanNvimboat
	dbUpdate(nb)
}

func nvimboatLoop(cnb chan *nvimboat.Nvimboat) {
	nb := new(nvimboat.Nvimboat)
	cnb <- nb
	nb.SyncDBchan = make(chan nvimboat.SyncDB)
	nvimPlugin.Main(func(p *nvimPlugin.Plugin) (err error) {
		p.HandleCommand(&nvimPlugin.CommandOptions{Name: "Nvimboat", NArgs: "+", Complete: "customlist,CompleteNvimboat"}, nb.Command)
		p.HandleFunction(&nvimPlugin.FunctionOptions{Name: "CompleteNvimboat"}, nvimboat.CompleteNvimboat)
		return
	})
}
