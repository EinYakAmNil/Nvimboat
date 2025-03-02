package nvimboat

import (
	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type Filter struct {
	Name        string
	ID          string
	Query       string
	IncludeTags map[string]bool
	ExcludeTags map[string]bool
	Articles    []rssdb.QueryFilterRow
}

func (f *Filter) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	return
}

func (f *Filter) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	return
}

func (f *Filter) ChildIdx(p Page) (idx int, err error) {
	return
}

func (f *Filter) Back(nb *Nvimboat) (cursor_x int, err error) {
	return
}
