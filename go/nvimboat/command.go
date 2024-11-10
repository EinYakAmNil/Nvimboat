package nvimboat

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
)

func (nb *Nvimboat) HandleAction(nv *nvim.Nvim, args []string) (err error) {
	if len(args) == 0 {
		return fmt.Errorf("no arguments supplied")
	}
	action, ok := Actions[args[0]]
	if ok {
		err = action(nb, nv, args...)
		if err != nil {
			nb.Log(err)
		}
		return
	} else {
		err = fmt.Errorf("'%s' is not implemented", args[0])
		return
	}
}

func CompleteNvimboat(args *nvim.CommandCompletionArgs) (suggestions []string, err error) {
	if args.ArgLead != "" {
		for command := range Actions {
			lcd := min(len(command), len(args.ArgLead))
			if args.ArgLead[:lcd] == command[:lcd] {
				suggestions = append(suggestions, command)
			}
		}
		return
	}
	suggestions = sortMapKeys(Actions)
	return
}
