package nvimboat

import (
	"fmt"
	"log"

	"github.com/neovim/go-client/nvim"
)

func (nb *Nvimboat) Command(nv *nvim.Nvim, args []string) (err error) {
	if len(args) == 0 {
		return fmt.Errorf("no command arguments supplied")
	}
	for command, action := range actions {
		if args[0] == command {
			err = action(nb, nv, args...)
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
	err = fmt.Errorf("'%s' command is not implemented", args[0])
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
		return suggestions, nil
	}
	return sortedMapKeys(actions), nil
}

var actions = map[string]Action{
	"enable":       (*Nvimboat).Enable,
	"disable":      (*Nvimboat).Disable,
	"show-main":    (*Nvimboat).ShowMain,
	"show-tags":    (*Nvimboat).ShowTags,
	"select":       (*Nvimboat).Select,
	"back":         (*Nvimboat).Back,
	"next-unread":  (*Nvimboat).NextUnread,
	"prev-unread":  (*Nvimboat).PrevUnread,
	"next-article": (*Nvimboat).NextArticle,
	"prev-article": (*Nvimboat).PrevArticle,
	"toggle-read":  (*Nvimboat).ToggleArticleRead,
	"delete":       (*Nvimboat).Delete,
}
