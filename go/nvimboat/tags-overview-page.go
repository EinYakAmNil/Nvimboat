package nvimboat

import (
	"errors"
	"fmt"
	"slices"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

// The attribute "Tags" maps the name of the tag to the corresponding feed URLs.
// It will be initialized by Nvimboat.ShowTags().
// Tag names are keys, URLs are values.
type TagsOverviewPage struct {
	Tags               map[string][]string
	PrevCursorPosition [2]int
}

func (tp *TagsOverviewPage) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	tag := new(TagFeeds)
	tag.Name = id
	feeds, err := dbh.Queries.QueryTagFeeds(dbh.Ctx, tp.Tags[id])
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/TagsOverviewPage.Select"))
		return
	}
	tag.Feeds = feeds
	tag.Urls = tp.Tags[id]
	return tag, err
}

func (tp *TagsOverviewPage) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	if len(tp.Tags) == 0 {
		err = setLines(nv, buf, []string{"No tags defined."})
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/TagsOverviewPage.Render"))
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
		err = errors.Join(err, errors.New("nvimboat/TagsOverviewPage.Render"))
		return
	}
	return
}

func (tp *TagsOverviewPage) ChildIdx(p Page) (idx int, err error) {
	switch tagFeeds := p.(type) {
	case *TagFeeds:
		tagNames := make([]string, 0, len(tp.Tags))
		for t := range tp.Tags {
			tagNames = append(tagNames, t)
		}
		slices.Sort(tagNames)
		for i, tagName := range tagNames {
			if tagFeeds.Name == tagName {
				return i + 1, nil
			}
		}
		err = fmt.Errorf(
			"nvimboat/TagsPage.Render: Could not find the tag: %s",
			tagFeeds.Name,
		)
		err = errors.Join(err, errors.New("nvimboat/TagsOverviewPage.ChildIdx"))
		return -1, err
	default:
		err = fmt.Errorf(
			"nvimboat/TagsPage.Render: Bad Page type: TagFeeds. Got: %T\n",
			p,
		)
		err = errors.Join(err, errors.New("nvimboat/TagsOverviewPage.ChildIdx"))
		return -1, err
	}
}

func (tp *TagsOverviewPage) Back() (cursor_x int, err error) {
	return tp.PrevCursorPosition[0], nil
}

func (tp *TagsOverviewPage) ToggleRead(dbh rssdb.DbHandle, ids []string) (err error) {
	Log("Read status toggling is not implemented for this page.")
	return
}

func (tp *TagsOverviewPage) NextUnread(dbh rssdb.DbHandle) (err error)           { return }
func (tp *TagsOverviewPage) PrevUnread(dbh rssdb.DbHandle) (err error)           { return }
func (tp *TagsOverviewPage) NextArticle(dbh rssdb.DbHandle) (err error)          { return }
func (tp *TagsOverviewPage) PrevArticle(dbh rssdb.DbHandle) (err error)          { return }
func (tp *TagsOverviewPage) Delete(dbh rssdb.DbHandle, ids []string) (err error) { return }
