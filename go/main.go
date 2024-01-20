package main

import (
	"github.com/neovim/go-client/nvim"
	nvimPlugin "github.com/neovim/go-client/nvim/plugin"
	"nvimboat"
)

func main() {
	nb := new(nvimboat.Nvimboat)
	defer nb.LogFile.Close()
	// defer nb.DB.Close()

	nvimPlugin.Main(func(p *nvimPlugin.Plugin) error {
		if p.Nvim != nil {
			err := nb.Init(p)
			if err != nil {
				nb.Log(err)
			}
		}
		p.HandleCommand(&nvimPlugin.CommandOptions{Name: "Nvimboat", NArgs: "+", Complete: "customlist,CompleteNvimboat"}, nb.Command)
		p.HandleFunction(&nvimPlugin.FunctionOptions{Name: "CompleteNvimboat"}, func(c *nvim.CommandCompletionArgs) ([]string, error) {
			defer func() {
				err := recover()
				if err != nil {
					nb.Log(err)
				}
			}()
			var suggestions []string

			if c.ArgLead != "" {
				for _, a := range nvimboat.Actions {
					lcd := min(len(a), len(c.ArgLead))
					if c.ArgLead[:lcd] == a[:lcd] {
						suggestions = append(suggestions, a)
					}
				}
				return suggestions, nil
			}
			return nvimboat.Actions, nil
		})
		return nil
	})
}
