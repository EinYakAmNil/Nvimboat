package nvimboat

import (
	"errors"
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

func (nb *Nvimboat) Enable(nv *nvim.Nvim, args ...string) (err error) {
	err = initNvimboat(nb, nv)
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
		err = errors.Join(err, errors.New("nvimboat/Enable"))
		return
	}
	return
}

func (nb *Nvimboat) Disable(nv *nvim.Nvim, args ...string) (err error) {
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

func (nb *Nvimboat) Reload(nv *nvim.Nvim, args ...string) (err error) {
	if len(args) < 1 {
		err = fmt.Errorf("Reload: expected at least one argument")
		return
	}
	// reload all feeds if no arguments are given to the subcommand
	var feedUrls []string
	if len(args) == 1 {
		for feedUrl := range Feeds {
			feedUrls = append(feedUrls, feedUrl)
		}
	} else {
		feedUrls = args[1:]
	}
	Global.ChanAsync <- Async{func(a ...any) (err error) {
		err = ReloadFeeds(feedUrls)
		if err != nil {
			err = fmt.Errorf("ReloadFeeds: %w\n"+
				"nvimboat/Nvimboat.Reload", err,
			)
			return
		}
		return
	}, nil}
	return
}

func (nb *Nvimboat) ShowMain(nv *nvim.Nvim, args ...string) (err error) {
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
		if _, ok := Feeds[feed.Rssurl]; !ok {
			continue
		}
		Feeds[feed.Rssurl].Title = feed.Title
		Feeds[feed.Rssurl].ArticleCount = feed.ArticleCount
		Feeds[feed.Rssurl].UnreadCount = feed.UnreadCount
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
	err = Pages.Show()
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/ShowMain"))
		return
	}
	return
}

func (nb *Nvimboat) Select(nv *nvim.Nvim, args ...string) (err error) {
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
	err = Pages.Show()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Select"))
		return
	}
	return
}

func (nb *Nvimboat) Open(nv *nvim.Nvim, args ...string) (err error) {
	err = Pages.Top().Open(args[1:]...)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Nvimboat.Open"))
		return
	}
	return
}

func (nb *Nvimboat) ShowTags(nv *nvim.Nvim, args ...string) (err error) {
	cursorPosition, err := nv.WindowCursor(*NvWindow)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("nvimboat/ShowTags"))
		return
	}
	oldPageStack := Pages
	err = Pages.Reset()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Nvimboat.ShowTags"))
		return
	}
loopPages:
	for _, page := range oldPageStack {
		switch page.(type) {
		case *TagsOverview:
			break loopPages
		default:
			err = Pages.Push(page, page.ID())
			if err != nil {
				err = errors.Join(err, errors.New("nvimboat/Nvimboat.ShowTags"))
				return
			}
		}
	}
	p := new(TagsOverview)
	p.PrevCursorPosition = cursorPosition
	err = Pages.Push(p, p.ID())
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

func (nb *Nvimboat) Back(nv *nvim.Nvim, args ...string) (err error) {
	switch Pages.Top().(type) {
	case *MainMenu:
		return nil
	default:
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

func (nb *Nvimboat) NextUnread(nv *nvim.Nvim, args ...string) (err error) {
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

func (nb *Nvimboat) PrevUnread(nv *nvim.Nvim, args ...string) (err error) {
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

func (nb *Nvimboat) NextArticle(nv *nvim.Nvim, args ...string) (err error) {
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
			err = nb.Select(nv, "", f.Articles[idx].Url)
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
			err = nb.Select(nv, "", f.Articles[idx].Url)
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
				`Only use this func (nb *Nvimboat)tion for pages of type "Article" not "%T".`,
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

func (nb *Nvimboat) PrevArticle(nv *nvim.Nvim, args ...string) (err error) {
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
			err = nb.Select(nv, "", f.Articles[idx].Url)
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
			err = nb.Select(nv, "", f.Articles[idx].Url)
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
				`Only use this func (nb *Nvimboat)tion for pages of type "Article" not "%T".`,
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

func (nb *Nvimboat) ToggleRead(nv *nvim.Nvim, args ...string) (err error) {
	if len(args) < 2 {
		err = fmt.Errorf(`not enough arguments. Given arguments:`)
		for _, arg := range args {
			err = errors.Join(err, errors.New(arg))
		}
		err = errors.Join(err, errors.New("nvimboat/Nvimboat.ToggleRead"))
		return
	}
	dbh, err := rssdb.ConnectDb(DbPath)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ToggleRead: %w\n", err)
		return
	}
	cursorPosition, err := Pages.Top().ToggleRead(dbh, args[1:])
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ToggleRead: %w\n", err)
		return
	}
	err = Pages.Show()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Nvimboat.ToggleRead"))
		return
	}
	err = Nvim.SetWindowCursor(*NvWindow, cursorPosition[0])
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Nvimboat.ToggleRead"))
		return
	}
	mode := new(nvim.Mode)
	batch := Nvim.NewBatch()
	batch.Mode(mode)
	err = batch.Execute()
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/Nvimboat.ToggleRead"))
		return
	}
	switch mode.Mode {
	case "v", "V":
		batch.Command("normal! o")
		batch.SetWindowCursor(*NvWindow, cursorPosition[1])
		err = batch.Execute()
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/Nvimboat.ToggleRead"))
			return
		}
		return
	default:
		return
	}
}

func (nb *Nvimboat) Delete(nv *nvim.Nvim, args ...string) (err error) {
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
