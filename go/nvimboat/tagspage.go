package nvimboat

import (
	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type TagsPage struct {
}

func (tp *TagsPage) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	return
}

func (tp *TagsPage) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	return
}

func (tp *TagsPage) ChildIdx(p Page) (idx int, err error) {
	return
}

func (tp *TagsPage) Back(nb *Nvimboat) (cursor_x int, err error) {
	return
}
