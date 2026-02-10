package nvimboat

import (
	"errors"
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

func Enable(nv *nvim.Nvim, args ...string) (err error) {
	err = initNvimboat(nv)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/Enable"))
		return
	}
	err = Nvim.ExecLua(luaEnable, new(any))
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/Enable"))
		return
	}
	err = Nvim.Echo([]nvim.TextChunk{{
		Text: "Enabled Nvimboat",
	}},
		false,
		make(map[string]any),
	)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.NextUnread"))
		return
	}
	return
}

func Disable(nv *nvim.Nvim, args ...string) (err error) {
	err = Nvim.ExecLua(luaDisable, new(any))
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/Disable"))
		return
	}
	err = Nvim.Echo([]nvim.TextChunk{{
		Text: "Disabled Nvimboat",
	}},
		false,
		make(map[string]any),
	)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Feed.NextUnread"))
		return
	}
	return
}

func Reload(nv *nvim.Nvim, args ...string) (err error) {
	if len(args) < 1 {
		err = fmt.Errorf("Reload: expected at least one argument")
		return
	}
	// reload all feeds if no arguments are given to the subcommand
	var feedUrls []string
	if len(args) == 1 {
		for feedUrl := range FeedConfig {
			feedUrls = append(feedUrls, feedUrl)
		}
	} else {
		feedUrls = args[1:]
	}
	err = ReloadFeeds(feedUrls)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/Reload"))
		return
	}
	return
}

func ShowMain(nv *nvim.Nvim, args ...string) (err error) {
	mm := new(MainMenu)
	dbh, err := rssdb.ConnectDb(DbPath)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/ShowMain"))
		return
	}
	mainPageFeeds, err := dbh.Queries.QueryMainPage(dbh.Ctx)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/ShowMain"))
		return
	}
	for _, feed := range mainPageFeeds {
		tags := make(map[string]bool)
		for _, tag := range FeedConfig[feed.Rssurl] {
			tags[tag] = true
		}
		mm.Feeds = append(mm.Feeds, MainPageFeed{
			QueryMainPageRow: feed,
			Tags:             tags,
		})
	}
	err = updateFilters(dbh)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/ShowMain"))
		return
	}
	err = Pages.Reset()
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/ShowMain"))
		return
	}
	err = Pages.Push(mm, "")
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/ShowMain"))
		return
	}
	err = Nvim.SetWindowOption(*NvWindow, "wrap", false)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/ShowMain"))
		return
	}
	err = Pages.Show()
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/ShowMain"))
		return
	}
	return
}

func Select(nv *nvim.Nvim, args ...string) (err error) {
	if len(args) < 2 {
		err = fmt.Errorf(`not enough arguments. Given arguments:`)
		for _, arg := range args {
			err = errors.Join(err, errors.New(arg))
		}
		err = errors.Join(err, errors.New("nvimboat/Select"))
		return
	}
	dbh, err := rssdb.ConnectDb(DbPath)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Select"))
		return
	}
	p, err := Pages.Top().Select(dbh, args[1])
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Select"))
		return
	}
	err = Pages.Push(p, args[1])
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Select"))
		return
	}
	if _, ok := Pages.Top().(*Article); ok {
		err = Nvim.SetWindowOption(*NvWindow, "wrap", true)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Select"))
			return
		}
	}
	err = Pages.Show()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Select"))
		return
	}
	return
}

func ShowTags(nv *nvim.Nvim, args ...string) (err error) {
	cursorPosition, err := nv.WindowCursor(*NvWindow)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/ShowTags"))
		return
	}
	for i, page := range Pages {
		switch page.(type) {
		case *TagsOverview:
			Pages = Pages[:i+1]
			Pages.Show()
			return
		}
	}
	p := new(TagsOverview)
	p.PrevCursorPosition = cursorPosition
	p.Tags = make(map[string][]string)
	for url, tags := range FeedConfig {
		for _, t := range tags {
			p.Tags[t] = append(p.Tags[t], url)
		}
	}
	err = Pages.Push(p, "Tags")
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/ShowTags"))
		return
	}
	err = Pages.Show()
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/ShowTags"))
		return
	}
	return
}

func Back(nv *nvim.Nvim, args ...string) (err error) {
	switch Pages.Top().(type) {
	case *MainMenu:
		return nil
	default:
		if _, ok := Pages.Top().(*Article); ok {
			err = Nvim.SetWindowOption(*NvWindow, "wrap", false)
			if err != nil {
				err = errors.Join(err, errors.New("nvimboat/Back"))
				return
			}
		}
		var cursor_x int
		cursor_x, err = Pages.Top().Back()
		if err != nil {
			err = errors.Join(err, fmt.Errorf("nvimboat/Back"))
			return
		}
		defer Nvim.SetWindowCursor(*NvWindow, [2]int{cursor_x, 0})
		_, err = Pages.Pop()
		if err != nil {
			err = errors.Join(err, fmt.Errorf("nvimboat/Back"))
			return
		}
		err = Pages.Show()
		if err != nil {
			err = errors.Join(err, fmt.Errorf("nvimboat/Back"))
			return
		}
		return nil
	}
}

func NextUnread(nv *nvim.Nvim, args ...string) (err error) {
	dbh, err := rssdb.ConnectDb(DbPath)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/NextUnread"))
		return
	}
	err = Pages.Top().NextUnread(dbh)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/NextUnread"))
		return
	}
	return
}

func PrevUnread(nv *nvim.Nvim, args ...string) (err error) {
	dbh, err := rssdb.ConnectDb(DbPath)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/PrevUnread"))
		return
	}
	err = Pages.Top().PrevUnread(dbh)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/PrevUnread"))
		return
	}
	return
}

func NextArticle(nv *nvim.Nvim, args ...string) (err error) {
	switch p := Pages.Top().(type) {
	case *Article:
		var (
			f   Page
			idx int
		)
		f, err = Pages.Pop()
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/NextArticle"))
			return
		}
		switch f := f.(type) {
		case *Feed:
			idx, err = f.ChildIdx(p)
			if err != nil {
				err = errors.Join(err, errors.New("nvimboat/NextArticle"))
				return
			}
			idx = (idx + 1) % len(f.Articles)
			err = Select(nv, "", f.Articles[idx].Url)
			if err != nil {
				err = errors.Join(err, errors.New("nvimboat/NextArticle"))
				return
			}
		case *Filter:
			idx, err = f.ChildIdx(p)
			if err != nil {
				err = errors.Join(err, errors.New("nvimboat/NextArticle"))
				return
			}
			idx = (idx + 1) % len(f.Articles)
			err = Select(nv, "", f.Articles[idx].Url)
			if err != nil {
				err = errors.Join(err, errors.New("nvimboat/NextArticle"))
				return
			}
		default:
			err = fmt.Errorf(`Unexpected page type "%T".`, f)
			err = errors.Join(err, errors.New("nvimboat/NextArticle"))
			return
		}
	default:
		Nvim.Echo([]nvim.TextChunk{{
			Text: fmt.Sprintf(
				`Only use this function for pages of type "Article" not "%T".`,
				p,
			),
		}},
			true,
			make(map[string]any),
		)
		return
	}
	return
}

func PrevArticle(nv *nvim.Nvim, args ...string) (err error) {
	switch p := Pages.Top().(type) {
	case *Article:
		var (
			f   Page
			idx int
		)
		f, err = Pages.Pop()
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/PrevArticle"))
			return
		}
		switch f := f.(type) {
		case *Feed:
			idx, err = f.ChildIdx(p)
			if err != nil {
				err = errors.Join(err, errors.New("nvimboat/PrevArticle"))
				return
			}
			idx = idx - 1
			if idx < 0 {
				idx = len(f.Articles) - 1
			}
			err = Select(nv, "", f.Articles[idx].Url)
			if err != nil {
				err = errors.Join(err, errors.New("nvimboat/PrevArticle"))
				return
			}
		case *Filter:
			idx, err = f.ChildIdx(p)
			if err != nil {
				err = errors.Join(err, errors.New("nvimboat/PrevArticle"))
				return
			}
			idx = idx - 1
			if idx < 0 {
				idx = len(f.Articles) - 1
			}
			err = Select(nv, "", f.Articles[idx].Url)
			if err != nil {
				err = errors.Join(err, errors.New("nvimboat/PrevArticle"))
				return
			}
		default:
			err = fmt.Errorf(`Unexpected page type "%T".`, f)
			err = errors.Join(err, errors.New("nvimboat/PrevArticle"))
			return
		}
	default:
		Nvim.Echo([]nvim.TextChunk{{
			Text: fmt.Sprintf(
				`Only use this function for pages of type "Article" not "%T".`,
				p,
			),
		}},
			true,
			make(map[string]any),
		)
		return
	}
	return
}

func ToggleRead(nv *nvim.Nvim, args ...string) (err error) {
	if len(args) < 2 {
		err = fmt.Errorf(`not enough arguments. Given arguments:`)
		for _, arg := range args {
			err = errors.Join(err, errors.New(arg))
		}
		err = errors.Join(err, errors.New("nvimboat/ToggleRead"))
		return
	}
	dbh, err := rssdb.ConnectDb(DbPath)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ToggleRead: %w\n", err)
		return
	}
	err = Pages.Top().ToggleRead(dbh, args[1:])
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ToggleRead: %w\n", err)
		return
	}
	return
}

func Delete(nv *nvim.Nvim, args ...string) (err error) {
	if len(args) < 2 {
		err = fmt.Errorf(`not enough arguments. Given arguments:`)
		for _, arg := range args {
			err = errors.Join(err, errors.New(arg))
		}
		err = errors.Join(err, errors.New("nvimboat/Delete"))
		return
	}
	dbh, err := rssdb.ConnectDb(DbPath)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Delete"))
		return
	}
	err = Pages.Top().Delete(dbh, args[1:])
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Delete"))
		return
	}
	return
}
