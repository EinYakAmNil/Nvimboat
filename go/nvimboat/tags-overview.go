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
type TagsOverview struct {
	PrevCursorPosition [2]int
}

func (tp *TagsOverview) ID() string {
	return "Tags"
}

func (tp *TagsOverview) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	tag := new(TagFeeds)
	tag.Name = id
	tag.Urls = TagConfig[id]
	return tag, err
}

func (tp *TagsOverview) Open(urls ...string) (err error) {
	return
}

func (tp *TagsOverview) Render(nv *nvim.Nvim, buf nvim.Buffer) (err error) {
	if len(TagConfig) == 0 {
		err = setLines(nv, buf, []string{"No tags defined."})
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/TagsOverviewPage.Render"))
			return
		}
		return
	}
	var lines []string
	for tag, urls := range TagConfig {
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

func (tp *TagsOverview) ChildIdx(p Page) (idx int, err error) {
	switch tagFeeds := p.(type) {
	case *TagFeeds:
		tagNames := make([]string, 0, len(TagConfig))
		for t := range TagConfig {
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

func (tp *TagsOverview) Back() (cursor_x int, err error) {
	return tp.PrevCursorPosition[0], nil
}

func (tp *TagsOverview) ToggleRead(dbh rssdb.DbHandle, ids []string) (pos [2][2]int, err error) {
	pos, err = getCursorPositions()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/TagsOverview.ToggleRead"))
		return
	}
	err = Nvim.FeedKeys("\x1b", "n", false) // <Esc> = x1b
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/TagsOverview.ToggleRead"))
		return
	}
	Log("Read status toggling is not implemented for this page.")
	return
}

func (tp *TagsOverview) NextUnread(dbh rssdb.DbHandle) (err error)           { return }
func (tp *TagsOverview) PrevUnread(dbh rssdb.DbHandle) (err error)           { return }
func (tp *TagsOverview) Delete(dbh rssdb.DbHandle, ids []string) (err error) { return }
