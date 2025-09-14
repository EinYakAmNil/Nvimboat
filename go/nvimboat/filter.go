package nvimboat

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

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

func (f *Filter) Select(dbh rssdb.DbHandle, id string) (p Page, err error) {
	articleInfo, err := dbh.Queries.GetArticle(dbh.Ctx, id)
	if err != nil {
		err = fmt.Errorf("nvimboat/Filter.Select: %w\n", err)
		return
	}
	err = dbh.Queries.SetArticlesRead(dbh.Ctx, []string{id})
	if err != nil {
		err = fmt.Errorf("nvimboat/Filter.Select: %w\n", err)
		return
	}
	p = &Article{articleInfo}
	idx, err := f.ChildIdx(p)
	if err != nil {
		err = fmt.Errorf("nvimboat/Filter.Select: %w\n", err)
		return
	}
	f.Articles[idx].Unread = 0
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
	var (
		readStatusCol []string
		parsedTime    string
		pubDateCol    []string
		authorCol     []string
		titleCol      []string
		urlCol        []string
	)
	for _, a := range f.Articles {
		switch a.Unread {
		case 0:
			readStatusCol = append(readStatusCol, " ")
		case 1:
			readStatusCol = append(readStatusCol, "N")
		default:
			err = fmt.Errorf(`nvimboat/Filter.Render: Bad unread number for "%s" in feed %s: %d\n`,
				a.Url,
				f.Name,
				a.Unread,
			)
			return
		}
		parsedTime, err = unixToDate(a.Pubdate)
		if err != nil {
			err = fmt.Errorf("nvimboat/Filter.Render: %w\n", err)
			return
		}
		pubDateCol = append(pubDateCol, parsedTime)
		authorCol = append(authorCol, a.Author)
		titleCol = append(titleCol, a.Title)
		urlCol = append(urlCol, a.Url)
	}
	for _, c := range [][]string{readStatusCol, pubDateCol, authorCol, titleCol, urlCol} {
		err = addColumn(nv, buf, c)
		if err != nil {
			err = fmt.Errorf("nvimboat/Filter.Render: %w\n", err)
			return
		}
	}
	return
}

func (f *Filter) ChildIdx(p Page) (idx int, err error) {
	childDate := p.(*Article).Pubdate
	var (
		section     = len(f.Articles)
		searchRange = f.Articles
	)
	for range f.Articles {
		if childDate > searchRange[section/2].Pubdate {
			searchRange = searchRange[:section/2]
		} else if childDate < searchRange[section/2].Pubdate {
			idx += section / 2
			searchRange = searchRange[section/2:]
		} else if childDate == searchRange[section/2].Pubdate {
			idx += section / 2
			return
		}
		section = len(searchRange)
	}
	return -1, fmt.Errorf(
		`"%v" doesn't contain: "%+v"`,
		prettyStruct(f),
		prettyStruct(p),
	)
}

func (f *Filter) Back() (cursor_x int, err error) {
	filterNames := make([]string, 0, len(Filters))
	for name := range Filters {
		filterNames = append(filterNames, name)
	}
	slices.Sort(filterNames)
	for idx, filterName := range filterNames {
		if f.Name == filterName {
			cursor_x = idx + 1
			return
		}
	}
	err = fmt.Errorf(
		"nvimboat/Filter.Back: cannot find index for %s\n",
		prettyStruct(f),
	)
	return -1, err
}

func (f *Filter) ToggleRead(dbh rssdb.DbHandle, ids []string) (err error) {
	setArticlesRead := false
checkAnyUnread:
	for _, a := range f.Articles {
		for _, id := range ids {
			if a.Url == id && a.Unread == 1 {
				setArticlesRead = true
				break checkAnyUnread
			}
		}
	}
	if setArticlesRead {
		err = dbh.Queries.SetArticlesRead(dbh.Ctx, ids)
		if err != nil {
			err = fmt.Errorf("nvimboat/Filter.ToggleRead: %w\n", err)
			return
		}
	outer1:
		for i, a := range f.Articles {
			for _, id := range ids {
				if a.Url == id && a.Unread == 1 {
					f.Articles[i].Unread = 0
					continue outer1
				}
			}
		}
	} else {
		err = dbh.Queries.SetArticlesUnread(dbh.Ctx, ids)
		if err != nil {
			err = fmt.Errorf("nvimboat/Filter.ToggleRead: %w\n", err)
			return
		}
	outer2:
		for i, a := range f.Articles {
			for _, id := range ids {
				if a.Url == id && a.Unread == 0 {
					f.Articles[i].Unread = 1
					continue outer2
				}
			}
		}
	}
	err = Pages.Show()
	if err != nil {
		err = fmt.Errorf("nvimboat/Filter.ToggleRead: %w\n", err)
		return
	}
	return
}

func updateFilters(dbh rssdb.DbHandle) (err error) {
	for _, filter := range Filters {
		filter.Articles, err = dbh.Queries.QueryFilter(dbh.Ctx, filter.QueryFilterParams)
		if err != nil {
			err = fmt.Errorf("nvimboat/MainMenu.UpdateFilters: %w\n", err)
			return
		}
	}
	return
}

func parseFilter(rawFilter map[string]any) (filter Filter, err error) {
	var (
		descriptionTags []string
		descriptionSql  []string
		descriptions    []string
	)
	configMapping := map[string]*string{
		"name":              &filter.Name,
		"guid":              &filter.Guid,
		"title":             &filter.Title,
		"author":            &filter.Author,
		"url":               &filter.Url,
		"content":           &filter.Content,
		"content_mime_type": &filter.ContentMimeType,
	}
	for luaValue, filterAttr := range configMapping {
		if ok := assignFilterVarcharAttr(filterAttr, rawFilter[luaValue]); ok && filterAttr != &filter.Name {
			descriptionSql = append(descriptionSql, luaValue+": "+*filterAttr)
		}
	}
	if filter.Name == "" {
		err = fmt.Errorf(
			"nvimboat/parseFilters: no name for filter in: %+v\n",
			prettyStruct(filter),
		)
		return
	}
	if value, ok := rawFilter["unread"].(int64); ok {
		filter.UnreadStates = []int{int(value)}
		descriptionSql = append(descriptionSql, "unread: "+strconv.Itoa(int(value)))
	} else {
		filter.UnreadStates = []int{0, 1}
	}
	filter.ExcludeTags = make(map[string]bool)
	filter.IncludeTags = make(map[string]bool)
	if tags, okTags := rawFilter["tags"].([]any); okTags {
		for _, tag := range tags {
			if t, ok := tag.(string); ok {
				if len(t) == 0 {
					err = fmt.Errorf("nvimboat/parseFilters: string length 0. cannot parse %+v\n", rawFilter)
					return
				} else if t[0] == '!' {
					filter.ExcludeTags[t[1:]] = true
				} else {
					filter.IncludeTags[t] = true
				}
				descriptionTags = append(descriptionTags, t)
			}
		}
	} else {
		err = fmt.Errorf("nvimboat/parseFilters: cannot parse %+v\n", rawFilter)
		return
	}
	if len(descriptionSql) > 0 {
		descriptions = append(descriptions, descriptionSql...)
	}
	if len(descriptionTags) > 0 {
		descriptions = append(descriptions, "tags: "+strings.Join(descriptionTags, ", "))
	}
	filter.FilterDescription = "filter: " + strings.Join(descriptions, ", ")
filterFeed:
	for _, feed := range Feeds {
		for excTag := range filter.ExcludeTags {
			if feed.Tags[excTag] == true {
				continue filterFeed
			}
		}
		for incTag := range filter.IncludeTags {
			if feed.Tags[incTag] == true {
				filter.Feedurls = append(filter.Feedurls, feed.Rssurl)
				continue filterFeed
			}
		}
	}
	return
}

func assignFilterVarcharAttr(attribute *string, luaValue any) (replaced bool) {
	if value, ok := luaValue.(string); ok {
		*attribute = value
		return true
	} else {
		*attribute = "%"
		return false
	}
}
