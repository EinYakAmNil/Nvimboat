package nvimboat

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
)

type NvimboatAction func(*Nvimboat, *nvim.Nvim, ...string) error

var actions = map[string]NvimboatAction{
	"back":         (*Nvimboat).Back,
	"delete":       (*Nvimboat).Delete,
	"disable":      (*Nvimboat).Disable,
	"enable":       (*Nvimboat).Enable,
	"next-article": (*Nvimboat).NextArticle,
	"next-unread":  (*Nvimboat).NextUnread,
	"prev-article": (*Nvimboat).PrevArticle,
	"prev-unread":  (*Nvimboat).PrevUnread,
	"reload":       (*Nvimboat).Reload,
	"select":       (*Nvimboat).Select,
	"show-main":    (*Nvimboat).ShowMain,
	"show-tags":    (*Nvimboat).ShowTags,
	"toggle-read":  (*Nvimboat).ToggleRead,
}

type (
	Nvimboat struct {
		ChanAsync chan Async
	}
	Async struct {
		Function func(...any) error
		Args     []any
	}
)

func (nb *Nvimboat) HandleAction(nv *nvim.Nvim, args []string) (err error) {
	if len(args) == 0 {
		return fmt.Errorf("no arguments supplied")
	}
	var (
		action NvimboatAction
		ok     bool
	)
	if action, ok = actions[args[0]]; !ok {
		err = fmt.Errorf("'%s' is not implemented", args[0])
		return
	}
	err = action(nb, nv, args...)
	if err != nil {
		Log(err)
	}
	return
}

func CompleteNvimboat(args *nvim.CommandCompletionArgs) (suggestions []string, err error) {
	if args.ArgLead != "" {
		for command := range actions {
			lcd := min(len(command), len(args.ArgLead))
			if args.ArgLead[:lcd] == command[:lcd] {
				suggestions = append(suggestions, command)
			}
		}
		return
	}
	suggestions = sortMapKeys(actions)
	return
}
