package nvimboat

import (
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

func Enable(nv *nvim.Nvim, args ...string) (err error) {
	err = initNvimboat(nv)
	if err != nil {
		err = fmt.Errorf("Nvimboat enable: %w", err)
		return
	}
	err = Nvim.ExecLua(luaEnable, new(any))
	if err != nil {
		err = fmt.Errorf("Nvimboat enable: %w", err)
		return
	}
	Log("enabled Nvimboat")
	return
}

func Disable(nv *nvim.Nvim, args ...string) (err error) {
	err = Nvim.ExecLua(luaDisable, new(any))
	if err != nil {
		err = fmt.Errorf("Nvimboat disable: %w", err)
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
		err = fmt.Errorf("Reload: %w", err)
		return
	}
	return
}

func ShowMain(nv *nvim.Nvim, args ...string) (err error) {
	mm := new(MainMenu)
	dbh, err := rssdb.ConnectDb(DbPath)
	if err != nil {
		err = fmt.Errorf("ShowMain: %w", err)
		return
	}
	mainPageFeeds, err := dbh.Queries.QueryMainPage(dbh.Ctx)
	if err != nil {
		err = fmt.Errorf("ShowMain: %w", err)
		return
	}
	for _, feed := range mainPageFeeds {
		tags := make(map[string]bool)
		for _, tag := range FeedConfig[feed.Feedurl] {
			tags[tag] = true
		}
		mm.Feeds = append(mm.Feeds, MainPageFeed{
			QueryMainPageRow: feed,
			Tags:             tags,
		})
	}
	err = updateFilters(dbh)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ShowMain: %w\n", err)
		return
	}
	err = Pages.Reset()
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ShowMain: %w\n", err)
		return
	}
	err = Pages.Push(mm, "")
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ShowMain: %w\n", err)
		return
	}
	err = Pages.Show()
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ShowMain: %w\n", err)
		return
	}
	return
}

func Select(nv *nvim.Nvim, args ...string) (err error) {
	if len(args) < 2 {
		err = fmt.Errorf("nvimboat/Nvimboat.Select: no arguments")
		return
	}
	dbh, err := rssdb.ConnectDb(DbPath)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Select: %w\n", err)
		return
	}
	p, err := Pages.Top().Select(dbh, args[1])
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Select: %w\n", err)
		return
	}
	err = Pages.Push(p, args[1])
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Select: %w\n", err)
		return
	}
	err = Pages.Show()
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Select: %w\n", err)
		return
	}
	return
}

func ShowTags(nv *nvim.Nvim, args ...string) (err error) {
	cursorPosition, err := nv.WindowCursor(*NvWindow)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ShowTags: %w\n", err)
		return
	}
	for i, page := range Pages {
		switch page.(type) {
		case *TagsOverviewPage:
			Pages = Pages[:i+1]
			Pages.Show()
			return
		}
	}
	p := new(TagsOverviewPage)
	p.PrevCursorPosition = cursorPosition
	p.Tags = make(map[string][]string)
	for url, tags := range FeedConfig {
		for _, t := range tags {
			p.Tags[t] = append(p.Tags[t], url)
		}
	}
	err = Pages.Push(p, "Tags")
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Select: %w\n", err)
		return
	}
	err = Pages.Show()
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Select: %w\n", err)
		return
	}
	return
}

func Back(nv *nvim.Nvim, args ...string) error {
	switch Pages.Top().(type) {
	case *MainMenu:
		return nil
	default:
		cursor_x, err := Pages.Top().Back()
		if err != nil {
			return fmt.Errorf("nvimboat/Nvimboat.Back: %w\n", err)
		}
		defer Nvim.SetWindowCursor(*NvWindow, [2]int{cursor_x, 0})
		_, err = Pages.Pop()
		if err != nil {
			return fmt.Errorf("nvimboat/Nvimboat.Back: %w\n", err)
		}
		err = Pages.Show()
		if err != nil {
			return fmt.Errorf("nvimboat/Nvimboat.Back: %w\n", err)
		}
		return nil
	}
}

func NextUnread(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func PrevUnread(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func NextArticle(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func PrevArticle(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func ToggleRead(nv *nvim.Nvim, args ...string) (err error) {
	if len(args) < 2 {
		err = fmt.Errorf("nvimboat/Nvimboat.ToggleRead: no arguments")
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
	return
}
