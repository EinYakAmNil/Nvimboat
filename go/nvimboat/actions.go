package nvimboat

import (
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

var Actions = map[string]NvimboatAction{
	"back":         (*Nvimboat).Back,
	"delete":       (*Nvimboat).Delete,
	"disable":      (*Nvimboat).Disable,
	"enable":       (*Nvimboat).Enable,
	"next-article": (*Nvimboat).NextArticle,
	"next-unread":  (*Nvimboat).NextUnread,
	"prev-article": (*Nvimboat).PrevArticle,
	"prev-unread":  (*Nvimboat).PrevUnread,
	"reload":       (*Nvimboat).Reload,
	"select":       (*Nvimboat).Select,
	"show-main":    (*Nvimboat).ShowMain,
	"show-tags":    (*Nvimboat).ShowTags,
	"toggle-read":  (*Nvimboat).ToggleRead,
}

func (nb *Nvimboat) Enable(nv *nvim.Nvim, args ...string) (err error) {
	err = nb.init(nv)
	if err != nil {
		err = fmt.Errorf("Nvimboat enable: %w", err)
		return
	}
	err = nb.Nvim.ExecLua(luaEnable, new(any))
	if err != nil {
		err = fmt.Errorf("Nvimboat enable: %w", err)
		return
	}
	nb.Log("enabled Nvimboat")
	return
}

func (nb *Nvimboat) Disable(nv *nvim.Nvim, args ...string) (err error) {
	err = nb.Nvim.ExecLua(luaDisable, new(any))
	if err != nil {
		err = fmt.Errorf("Nvimboat disable: %w", err)
		return
	}
	return
}

func (nb *Nvimboat) Reload(nv *nvim.Nvim, args ...string) (err error) {
	if len(args) < 1 {
		err = fmt.Errorf("Reload: expected at least one argument")
		return
	}
	// reload all feeds if no arguments are given to the subcommand
	var feedUrls []string
	if len(args) == 1 {
		for feedUrl := range nb.FeedConfig {
			feedUrls = append(feedUrls, feedUrl)
		}
	} else {
		feedUrls = args[1:]
	}
	err = ReloadFeeds(nb, feedUrls)
	if err != nil {
		err = fmt.Errorf("Reload: %w", err)
		return
	}
	return
}

func (nb *Nvimboat) ShowMain(nv *nvim.Nvim, args ...string) (err error) {
	mm := new(MainMenu)
	dbh, err := rssdb.ConnectDb(nb.DbPath)
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
		for _, tag := range nb.FeedConfig[feed.Feedurl] {
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
	err = nb.ResetPages()
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ShowMain: %w\n", err)
		return
	}
	err = nb.PushPage(mm, "")
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ShowMain: %w\n", err)
		return
	}
	err = nb.Show(mm)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ShowMain: %w\n", err)
		return
	}
	return
}

func (nb *Nvimboat) Select(nv *nvim.Nvim, args ...string) (err error) {
	if len(args) < 2 {
		err = fmt.Errorf("nvimboat/Nvimboat.Select: no arguments")
		return
	}
	dbh, err := rssdb.ConnectDb(nb.DbPath)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Select: %w\n", err)
		return
	}
	p, err := nb.Top().Select(dbh, args[1])
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Select: %w\n", err)
		return
	}
	err = nb.PushPage(p, args[1])
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Select: %w\n", err)
		return
	}
	err = nb.Show(p)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Select: %w\n", err)
		return
	}
	return
}

func (nb *Nvimboat) ShowTags(nv *nvim.Nvim, args ...string) (err error) {
	cursorPosition, err := nv.WindowCursor(*nb.Window)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ShowTags: %w\n", err)
		return
	}
	for i, page := range nb.Pages {
		switch page.(type) {
		case *TagsOverviewPage:
			nb.Pages = nb.Pages[:i+1]
			nb.Show(page)
			return
		}
	}
	p := new(TagsOverviewPage)
	p.PrevCursorPosition = cursorPosition
	p.Tags = make(map[string][]string)
	for url, tags := range nb.FeedConfig {
		for _, t := range tags {
			p.Tags[t] = append(p.Tags[t], url)
		}
	}
	err = nb.PushPage(p, "Tags")
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Select: %w\n", err)
		return
	}
	err = nb.Show(p)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Select: %w\n", err)
		return
	}
	return
}

func (nb *Nvimboat) Back(nv *nvim.Nvim, args ...string) error {
	switch nb.Top().(type) {
	case *MainMenu:
		return nil
	default:
		cursor_x, err := nb.Top().Back(nb)
		if err != nil {
			return fmt.Errorf("nvimboat/Nvimboat.Back: %w\n", err)
		}
		defer nb.Nvim.SetWindowCursor(*nb.Window, [2]int{cursor_x, 0})
		parentPage, err := nb.PopPage()
		if err != nil {
			return fmt.Errorf("nvimboat/Nvimboat.Back: %w\n", err)
		}
		err = nb.Show(parentPage)
		if err != nil {
			return fmt.Errorf("nvimboat/Nvimboat.Back: %w\n", err)
		}
		return nil
	}
}

func (nb *Nvimboat) NextUnread(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func (nb *Nvimboat) PrevUnread(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func (nb *Nvimboat) NextArticle(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func (nb *Nvimboat) PrevArticle(nv *nvim.Nvim, args ...string) (err error) {
	return
}

func (nb *Nvimboat) ToggleRead(nv *nvim.Nvim, args ...string) (err error) {
	if len(args) < 2 {
		err = fmt.Errorf("nvimboat/Nvimboat.ToggleRead: no arguments")
		return
	}
	dbh, err := rssdb.ConnectDb(nb.DbPath)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ToggleRead: %w\n", err)
		return
	}
	err = nb.Top().ToggleRead(dbh, args[1])
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ToggleRead: %w\n", err)
		return
	}
	return
}

func (nb *Nvimboat) Delete(nv *nvim.Nvim, args ...string) (err error) {
	return
}
