package nvimboat

import (
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type Filter struct {
	Name string
	rssdb.QueryFilterParams
	FilterDescription string
	IncludeTags       map[string]bool
	ExcludeTags       map[string]bool
	Articles          []rssdb.QueryFilterRow
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

func parseFilters(rawFilters []map[string]any) (filterConfig []*Filter, err error) {
	for _, filter := range rawFilters {
		f := new(Filter)
		configMapping := map[string]*string{
			"name":              &f.Name,
			"guid":              &f.Guid,
			"title":             &f.Title,
			"author":            &f.Author,
			"url":               &f.Url,
			"content":           &f.Content,
			"content_mime_type": &f.ContentMimeType,
		}
		for luaValue, filterAttr := range configMapping {
			assignFilterVarcharAttr(filterAttr, filter[luaValue])
		}
		if f.Name == "" {
			err = fmt.Errorf(
				"nvimboat/parseFilters: no name for filter in: %+v\n",
				prettyStruct(filter),
			)
			return
		}
		if value, ok := filter["unread"].(int); ok {
			f.UnreadStates = []int{value}
		} else {
			f.UnreadStates = []int{0, 1}
		}
		f.ExcludeTags = make(map[string]bool)
		f.IncludeTags = make(map[string]bool)
		var descriptionTags []string
		if tags, okTags := filter["tags"].([]any); okTags {
			f.FilterDescription += ", tags: "
			for _, tag := range tags {
				if t, ok := tag.(string); ok {
					if len(t) == 0 {
						err = fmt.Errorf("nvimboat/parseFilters: cannot parse %+v\n", filter)
						return
					} else if t[0] == '!' {
						f.ExcludeTags[t[1:]] = true
						descriptionTags = append(descriptionTags, t)
					} else {
						f.IncludeTags[t] = true
						descriptionTags = append(descriptionTags, t)
					}
				}
			}
		} else {
			fmt.Println(tags, filter["tags"])
			err = fmt.Errorf("nvimboat/parseFilters: cannot parse %+v\n", filter)
			return
		}
		filterConfig = append(filterConfig, f)
	filterFeed:
		for _, feed := range Feeds {
			for excTag := range f.ExcludeTags {
				if feed.Tags[excTag] == true {
					continue filterFeed
				}
			}
			for incTag := range f.IncludeTags {
				if feed.Tags[incTag] == true {
					f.Feedurls = append(f.Feedurls, feed.Rssurl)
					continue filterFeed
				}
			}
		}
	}
	return
}

func assignFilterVarcharAttr(attribute *string, luaValue any) {
	if value, ok := luaValue.(string); ok {
		*attribute = value
	} else {
		*attribute = "%"
	}
}
