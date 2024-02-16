package nvimboat

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
)

func (nb *Nvimboat) ShowMain(nv *nvim.Nvim, args ...string) (err error) {
	mainMenu, err := QueryMain(nb.DBHandler, nb.Feeds, nb.Filters)
	if err != nil {
		return
	}
	nb.Pages.Pages = nil
	err = nb.Push(mainMenu)
	return
}

func (nb *Nvimboat) Enable(nv *nvim.Nvim, args ...string) (err error) {
	err = nb.init(nv)
	if err != nil {
		return
	}
	err = nb.ShowMain(nv, args...)
	if err != nil {
		return
	}
	err = nv.ExecLua(nvimboatEnable, new(any))
	return
}

func (nb *Nvimboat) Disable(nv *nvim.Nvim, args ...string) (err error) {
	err = nv.ExecLua(nvimboatDisable, new(any))
	return
}

func (nb *Nvimboat) ShowTags(nv *nvim.Nvim, args ...string) (err error) {
	tags, err := QueryTags(nb.Feeds)
	if err != nil {
		return
	}
	err = nb.Push(tags)
	return
}

func (nb *Nvimboat) Select(nv *nvim.Nvim, args ...string) (err error) {
	defer nb.Nvim.SetWindowCursor(*nb.Window, [2]int{0, 1})
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments to call 'select'")
	}
	page, err := nb.Pages.Top().QueryChild(nb.DBHandler, args[1])
	if err != nil {
		return
	}
	err = nb.Push(page)
	return
}

func (nb *Nvimboat) Back(nv *nvim.Nvim, args ...string) (err error) {
	switch nb.Pages.Top().(type) {
	case *MainMenu:
		return
	default:
		err = nb.Pop()
		return
	}
}

func (nb *Nvimboat) NextUnread(nv *nvim.Nvim, args ...string) error {
	return nil
}

func (nb *Nvimboat) PrevUnread(nv *nvim.Nvim, args ...string) error {
	return nil
}

func (nb *Nvimboat) NextArticle(nv *nvim.Nvim, args ...string) error {
	return nil
}

func (nb *Nvimboat) PrevArticle(nv *nvim.Nvim, args ...string) error {
	return nil
}

func (nb *Nvimboat) ToggleArticleRead(nv *nvim.Nvim, args ...string) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments to call 'toggle-unread'")
	}
	err = nb.Pages.Top().ToggleUnread(*nb, args[1:]...)
	return
}
