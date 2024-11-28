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
	ChildIdx(id string) (p Page)
}

func (nb *Nvimboat) Show(p Page, id string) (err error) {
	err = setLines(nb.Nvim, *nb.Buffer, []string{""})
	if err != nil {
		err = fmt.Errorf("Show: %w", err)
		return
	}
	defer trimTrail(nb.Nvim, *nb.Buffer)
	defer nb.Nvim.SetWindowCursor(*nb.Window, [2]int{0, 1})
	if err != nil {
		err = fmt.Errorf("Show: %w", err)
		return
	}
	p.Render(nb.Nvim, *nb.Buffer)
	if err != nil {
		err = fmt.Errorf("Show: %w", err)
		return
	}
	pageType := fmt.Sprintf("%T", p)
	_, pageType, _ = strings.Cut(pageType, "nvimboat.")
	err = nb.Nvim.ExecLua(luaPushPage, new(any), pageType, id)
	if err != nil {
		err = fmt.Errorf("Show: %w", err)
		return
	}
	nb.Pages.Push(p)
	return
}

func (ps *PageStack) Top() (p Page) {
	if pageCount := len(ps.Pages); pageCount > 0 {
		return ps.Pages[pageCount-1]
	}
	return nil
}

func (ps *PageStack) Push(p Page) {
	ps.Pages = append(ps.Pages, p)
}
