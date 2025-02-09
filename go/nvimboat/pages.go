package nvimboat

import (
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type Page interface {
	Select(dbh rssdb.DbHandle, id string) (p Page, err error)
	Render(nv *nvim.Nvim, buf nvim.Buffer) (err error)
	ChildIdx(p Page) (idx int, err error)
	Back(nb *Nvimboat) (cursor_x int, err error)
}

func (nb *Nvimboat) Show(p Page) (err error) {
	err = setLines(nb.Nvim, *nb.Buffer, []string{""})
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Show: %w\n", err)
		return
	}
	defer trimTrail(nb.Nvim, *nb.Buffer)
	defer nb.Nvim.SetWindowCursor(*nb.Window, [2]int{0, 1})
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Show: %w\n", err)
		return
	}
	p.Render(nb.Nvim, *nb.Buffer)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Show: %w\n", err)
		return
	}
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

func (ps *PageStack) Pop() (p Page) {
	if len(ps.Pages) > 1 {
		ps.Pages = ps.Pages[:len(ps.Pages)-1]
	}
	return ps.Top()
}
