package main

import (
	"errors"
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/nvimboat"
	nvimPlugin "github.com/neovim/go-client/nvim/plugin"
)

func main() {
	nb := new(nvimboat.Nvimboat)
	nb.ChanAsync = make(chan nvimboat.Async)
	go execAsync(nb)
	nvimPlugin.Main(func(p *nvimPlugin.Plugin) (err error) {
		p.HandleCommand(
			&nvimPlugin.CommandOptions{
				Name:     "Nvimboat",
				NArgs:    "+",
				Complete: "customlist, CompleteNvimboat",
			},
			nb.HandleAction,
		)
		p.HandleFunction(
			&nvimPlugin.FunctionOptions{Name: "CompleteNvimboat"},
			nvimboat.CompleteNvimboat,
		)
		return
	})
}

func execAsync(nb *nvimboat.Nvimboat) (err error) {
	if nb.ChanAsync == nil {
		err = fmt.Errorf(`No channel.`)
		err = errors.Join(err, errors.New("main/execAsync"))
		return
	}

	maxWorkers := make(chan struct{}, 4)

	for task := range nb.ChanAsync {
		maxWorkers <- struct{}{}
		go func(t nvimboat.Async) {
			defer func() { <-maxWorkers }()
			task.Function(task.Args...)
		}(task)
	}
	return
}
