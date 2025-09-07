package nvimboat

import (
	"fmt"
	"strings"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type Page interface {
	Select(dbh rssdb.DbHandle, id string) (p Page, err error)
	Render(nv *nvim.Nvim, buf nvim.Buffer) (err error)
	ChildIdx(p Page) (idx int, err error)
	Back(nb *Nvimboat) (cursor_x int, err error)
	ToggleRead(dbh rssdb.DbHandle, id string) (err error)
}

func (nb *Nvimboat) Show(p Page) (err error) {
	err = setLines(nb.Nvim, *NbBuffer, []string{""})
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Show: %w\n", err)
		return
	}
	defer trimTrail(nb.Nvim, *NbBuffer)
	defer nb.Nvim.SetWindowCursor(*nb.Window, [2]int{0, 1})
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Show: %w\n", err)
		return
	}
	p.Render(nb.Nvim, *NbBuffer)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Show: %w\n", err)
		return
	}
	return
}

func (nb *Nvimboat) Top() (p Page) {
	if pageCount := len(nb.Pages); pageCount > 0 {
		return nb.Pages[pageCount-1]
	}
	return nil
}

func (nb *Nvimboat) PushPage(p Page, id string) (err error) {
	nb.Pages = append(nb.Pages, p)
	pageType := fmt.Sprintf("%T", p)
	_, pageType, _ = strings.Cut(pageType, "nvimboat.")
	err = nb.Nvim.ExecLua(luaPushPage, new(any), pageType, id)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ShowMain: %w\n", err)
		return
	}
	return
}

func (nb *Nvimboat) PopPage() (p Page, err error) {
	if len(nb.Pages) > 1 {
		nb.Pages = nb.Pages[:len(nb.Pages)-1]
	}
	err = nb.Nvim.ExecLua(luaPopPage, new(any))
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.PopPage: %w\n", err)
		return
	}
	return nb.Top(), nil
}

func (nb *Nvimboat) ResetPages() (err error) {
	nb.Pages = []Page{}
	err = nb.Nvim.ExecLua(luaResetPages, new(any))
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ResetPages: %w\n", err)
		return
	}
	return
}
