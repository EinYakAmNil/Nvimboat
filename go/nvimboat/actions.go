package nvimboat

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
)

var Actions = map[string]NvimboatAction{
	"enable":       (*Nvimboat).Enable,
	"disable":      (*Nvimboat).Disable,
	"reload":       (*Nvimboat).Reload,
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

func (nb *Nvimboat) init(nv *nvim.Nvim) (err error) {
	rawConfig := make(map[string]any)
	nb.Nvim = nv
	nb.Buffer = new(nvim.Buffer)
	nb.Window = new(nvim.Window)
	execBatch := nv.NewBatch()
	execBatch.CurrentWindow(nb.Window)
	execBatch.CurrentBuffer(nb.Buffer)
	execBatch.ExecLua(luaConfig, &rawConfig)
	execBatch.ExecLua(luaFeeds, &nb.Feeds)
	err = execBatch.Execute()
	if err != nil {
		err = fmt.Errorf("Nvimboat init: %w", err)
		return
	}
	err = parseConfig(nb, rawConfig)
	if err != nil {
		err = fmt.Errorf("Nvimboat init parse lua config: %w", err)
		return
	}
	err = SetupLogging(nb.LogPath)
	if err != nil {
		err = fmt.Errorf("Nvimboat init logging: %w", err)
		return
	}
	return
}

func (nb *Nvimboat) Enable(nv *nvim.Nvim, args ...string) (err error) {
	err = nb.init(nv)
	if err != nil {
		err = fmt.Errorf("Nvimboat enable: %w", err)
		return
	}
	err = nb.Nvim.ExecLua(luaEnable, new(any))
	if err != nil {
		err = fmt.Errorf("Nvimboat enable: %w", err)
		return
	}
	nb.Log("enabled Nvimboat")
	return
}

func (nb *Nvimboat) Disable(nv *nvim.Nvim, args ...string) (err error) {
	err = nb.Nvim.ExecLua(luaDisable, new(any))
	if err != nil {
		err = fmt.Errorf("Nvimboat disable: %w", err)
		return
	}
	return
}

func (nb *Nvimboat) Reload(nv *nvim.Nvim, args ...string) (err error) {
	if len(args) < 1 {
		err = fmt.Errorf("reload: expected at least one argument")
		return
	}
	// reload all feeds if no arguments are given to the subcommand
	var feedUrls []string
	if len(args) == 1 {
		for _, fu := range nb.Feeds {
			feedUrls = append(feedUrls, fu.Rssurl)
		}
	} else {
		feedUrls = args[1:]
	}
	err = nb.ReloadFeeds(feedUrls)
	if err != nil {
		err = fmt.Errorf("reload: %w", err)
		return
	}
	return
}

func (nb *Nvimboat) ShowMain(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func (nb *Nvimboat) ShowTags(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func (nb *Nvimboat) Select(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func (nb *Nvimboat) Back(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func (nb *Nvimboat) NextUnread(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func (nb *Nvimboat) PrevUnread(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func (nb *Nvimboat) NextArticle(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func (nb *Nvimboat) PrevArticle(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func (nb *Nvimboat) ToggleArticleRead(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func (nb *Nvimboat) Delete(nv *nvim.Nvim, args ...string) (err error) {
	return
}
