package nvimboat

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/neovim/go-client/nvim"
)

func setCursorUnread(row, col, maxRows int, matched any) (err error) {
	newCursorPosition := [2]int{
		row,
		col,
	}
	err = Nvim.SetWindowCursor(*NvWindow,
		newCursorPosition,
	)
	if err != nil {
		errPosition := fmt.Errorf(
			`Got: %+v, max row count: %d, %d`,
			newCursorPosition,
			row,
			maxRows,
		)
		errArticle := fmt.Errorf(
			`Matched: %+v`,
			prettyStruct(matched),
		)
		err = errors.Join(
			err,
			errPosition,
			errArticle,
			errors.New("nvimboat/setCursorNextUnread"),
		)
		return
	}
	return
}

func updateFilters(dbh rssdb.DbHandle) (err error) {
	for _, filter := range Filters {
		filter.Articles, err = dbh.Queries.QueryFilter(dbh.Ctx, filter.QueryFilterParams)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/updateFilters"))
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
			"No name for filter in: %+v",
			prettyStruct(filter),
		)
		err = errors.Join(err, errors.New("nvimboat/parseFilter"))
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
					err = fmt.Errorf(`Length of string is 0. cannot parse %+v`, rawFilter)
					err = errors.Join(err, errors.New("nvimboat/parseFilter"))
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
		err = fmt.Errorf(`Can't parse %+v`, rawFilter)
		err = errors.Join(err, errors.New("nvimboat/parseFilter"))
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

func selectFeed(dbh rssdb.DbHandle, feedurl string) (feed *Feed, err error) {
	feed = new(Feed)
	feed.Articles, err = dbh.Queries.GetFeedPage(dbh.Ctx, feedurl)
	feed.Rssurl = feedurl
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/selectFeed"))
		return
	}
	feedInfo, err := dbh.Queries.GetFeed(dbh.Ctx, feedurl)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/selectFeed"))
		return
	}
	feed.RssFeed = feedInfo
	return
}

func extracUrls(content string) (links []string) {
	re := regexp.MustCompile(
		`\b((?:https?|ftp|file):\/\/[-a-zA-Z0-9+&@#\/%?=~_|!:,.;]*[-a-zA-Z0-9+&@#\/%=~_|])`,
	)
	matches := re.FindAll([]byte(content), -1)
	for _, l := range matches {
		links = append(links, string(l))
	}
	return links
}

func renderHTML(content string) ([]string, error) {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(content)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/renderHTML"))
		return nil, err
	}
	return strings.Split(markdown, "\n"), nil
}

func unixToDate(unixTime int64) (string, error) {
	tz, err := time.LoadLocation("Local")
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/unixToDate"))
		return "", err
	}
	ut := time.Unix(unixTime, 0)
	dateString := ut.In(tz).Format("02 Jan 06")

	return dateString, nil
}

func trimTrail(nv *nvim.Nvim, buffer nvim.Buffer) (err error) {
	currentLines, err := nv.BufferLines(buffer, 0, -1, false)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/trimTrail"))
		return
	}
	var lines []string
	for _, l := range currentLines {
		lines = append(lines, strings.TrimRight(string(l), " "))
	}
	err = setLines(nv, buffer, lines)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/trimTrail"))
		return
	}
	return
}

func addColumn(nv *nvim.Nvim, buf nvim.Buffer, col []string) (err error) {
	currentLines, err := nv.BufferLines(buf, 0, -1, false)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/addColumn"))
		return
	}
	var (
		diff  int
		lines = []string{}
	)
	diff = (len(col) - len(currentLines))
	for range diff {
		currentLines = append(currentLines, []byte{})
	}
	for i, c := range col {
		if len(currentLines[i]) == 0 {
			lines = append(lines, c)
		} else {
			lines = append(lines, string(currentLines[i])+ColumnSeparator+c)
		}
	}
	err = setLines(nv, buf, lines)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/addColumn"))
		return
	}
	vcl, err := virtColLens(nv)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/addColumn"))
		return
	}
	maxLineLen := slices.Max(vcl)

	for i, l := range lines {
		diff = maxLineLen - vcl[i]
		lines[i] = l + strings.Repeat(" ", diff)
	}
	err = setLines(nv, buf, lines)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/addColumn"))
		return
	}
	return err
}

func setLines(nv *nvim.Nvim, buffer nvim.Buffer, lines []string) (err error) {
	err = nv.SetBufferLines(buffer, 0, -1, false, strings2bytes(lines))
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/setLines"))
		return
	}
	return
}

func virtColLens(nv *nvim.Nvim) (evalResult []int, err error) {
	virtCols := "map(range(1, line('$')), \"virtcol([v:val, '$'])\")"
	err = nv.Eval(virtCols, &evalResult)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/virtColLens"))
		return
	}
	return
}

func makeUnreadRatio(unreadCount, articleCount int64) (prefix string) {
	if unreadCount > 0 {
		prefix = "N (" + strconv.FormatInt(unreadCount, 10) + "/" +
			strconv.FormatInt(articleCount, 10) + ")"
		return
	}
	prefix = "  (" + strconv.FormatInt(unreadCount, 10) + "/" +
		strconv.FormatInt(articleCount, 10) + ")"
	return
}

func strings2bytes(stringSlice []string) (byteSlices [][]byte) {
	for _, s := range stringSlice {
		byteSlices = append(byteSlices, []byte(s))
	}
	return
}

func sortMapKeys(m any) (keyList []string) {
	keys := reflect.ValueOf(m).MapKeys()
	for _, key := range keys {
		keyList = append(keyList, key.Interface().(string))
	}
	sort.Strings(keyList)
	return
}

func prettyStruct(s any) string {
	if err, ok := s.(error); ok {
		return fmt.Sprintf("%+v", err)
	}
	marshal, _ := json.MarshalIndent(s, "", "	")
	return string(marshal)
}
