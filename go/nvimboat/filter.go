package nvimboat

import (
	"fmt"

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

// TODO: Create a SQL-Query, that does not rely on injection anymore.
// Lua filter config won't have "query".
// Instead the keys will be the rss_item column names.
func (f *Filter) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	return
}

func (f *Filter) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	if len(f.Articles) == 0 {
		err = setLines(nv, buf, []string{"No Articles found."})
		if err != nil {
			err = fmt.Errorf("nvimboat/Filter.Render: %w\n", err)
			return
		}
		return
	}
	return
}

func (f *Filter) ChildIdx(p Page) (idx int, err error) {
	return
}

func (f *Filter) Back(nb *Nvimboat) (cursor_x int, err error) {
	return
}

func (f *Filter) ToggleRead(dbh rssdb.DbHandle, id string) (err error) {
	return
}
