package nvimboat

import (
	"github.com/neovim/go-client/nvim"
)

type Page interface {
	Select(nb *Nvimboat, id string) (err error)
	Render(nv *nvim.Nvim, buf nvim.Buffer) (err error)
}

func (nb *Nvimboat) Show(p Page) (err error) {
	p.Render(nb.Nvim, *nb.Buffer)
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
