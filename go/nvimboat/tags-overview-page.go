package nvimboat

import (
	"fmt"
	"slices"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

// The attribute "Tags" maps the name of the tag to the corresponding feed URLs.
// It will be initialized by Nvimboat.ShowTags().
// Tag names are keys, URLs are values.
type TagsOverviewPage struct {
	Tags map[string][]string
}

func (tp *TagsOverviewPage) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	_, err = dbh.Queries.QueryTagFeeds(dbh.Ctx, tp.Tags[id])
	return
}

func (tp *TagsOverviewPage) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	if len(tp.Tags) == 0 {
		err = setLines(nv, buf, []string{"No tags defined."})
		if err != nil {
			err = fmt.Errorf("nvimboat/TagsPage.Render: %w\n", err)
			return
		}
		return
	}
	var lines []string
	for tag, urls := range tp.Tags {
		lines = append(lines, fmt.Sprintf(`%s (%d)`, tag, len(urls)))
	}
	slices.Sort(lines)
	err = setLines(nv, buf, lines)
	if err != nil {
		err = fmt.Errorf("nvimboat/TagsPage.Render: %w\n", err)
		return
	}
	return
}

func (tp *TagsOverviewPage) ChildIdx(p Page) (idx int, err error) {
	return
}

func (tp *TagsOverviewPage) Back(nb *Nvimboat) (cursor_x int, err error) {
	return
}
