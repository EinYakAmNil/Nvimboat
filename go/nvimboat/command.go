package nvimboat

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
)

type NvimboatAction func(*nvim.Nvim, ...string) error

var actions = map[string]NvimboatAction{
	"back":         Back,
	"delete":       Delete,
	"disable":      Disable,
	"enable":       Enable,
	"next-article": NextArticle,
	"next-unread":  NextUnread,
	"prev-article": PrevArticle,
	"prev-unread":  PrevUnread,
	"reload":       Reload,
	"select":       Select,
	"show-main":    ShowMain,
	"show-tags":    ShowTags,
	"toggle-read":  ToggleRead,
}

func HandleAction(nv *nvim.Nvim, args []string) (err error) {
	if len(args) == 0 {
		return fmt.Errorf("no arguments supplied")
	}
	action, ok := actions[args[0]]
	if ok {
		err = action(nv, args...)
		if err != nil {
			Log(err)
		}
		return
	} else {
		err = fmt.Errorf("'%s' is not implemented", args[0])
		return
	}
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
