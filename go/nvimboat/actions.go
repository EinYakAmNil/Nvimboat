package nvimboat

func (nb *Nvimboat) ShowMain() error {
	mainmenu, err := nb.QueryMain()
	if err != nil {
		return err
	}
	err = nb.Push(mainmenu)
	if err != nil {
		return err
	}
	nb.Pages.Pages = nb.Pages.Pages[:1]
	return nil
}

func (nb *Nvimboat) Enable() error {
	mainmenu, err := nb.QueryMain()
	if err != nil {
		return err
	}
	err = nb.Push(mainmenu)
	if err != nil {
		return err
	}
	err = nb.Nvim.Plugin.Nvim.ExecLua(nvimboatEnable, new(any))
	if err != nil {
		return err
	}
	return nil
}

func (nb *Nvimboat) Disable() error {
	return nil
}

func (nb *Nvimboat) ShowTags() error {
	tags, err := nb.QueryTags()
	if err != nil {
		return err
	}
	nb.Push(&tags)
	if err != nil {
		return err
	}
	return nil
}

func (nb *Nvimboat) Select(id string) error {
	defer nb.Nvim.Plugin.Nvim.SetWindowCursor(*nb.Nvim.Window, [2]int{0, 1})
	switch page := nb.Pages.Top().(type) {
	case *MainMenu:
		if id[:4] == "http" {
			feed, err := nb.QueryFeed(id)
			if err != nil {
				return err
			}
			err = nb.Push(&feed)
			if err != nil {
				return err
			}
		}
		if id[:6] == "query:" {
			query, inTags, exTags, err := parseFilterID(id)
			filter, err := nb.QueryFilter(query, inTags, exTags)
			filter.FilterID = id
			if err != nil {
				return err
			}
			err = nb.Push(&filter)
			if err != nil {
				return err
			}
		}
	case *Filter:
		articles := page.Articles
		for _, a := range articles {
			if a.Url == id {
				a.Unread = 0
				nb.ChanExecDB <- DBsync{Unread: 0, ArticleUrls: []string{a.Url}}
				err := nb.Push(a)
				return err
			}
		}
	case *Feed:
		articles := page.Articles
		for _, a := range articles {
			if a.Url == id {
				a.Unread = 0
				nb.ChanExecDB <- DBsync{Unread: 0, ArticleUrls: []string{a.Url}}
				err := nb.Push(a)
				return err
			}
		}
		// article, err := nb.QueryArticle(id)
		// if err != nil {
		// 	return err
		// }
		// nb.Push(&article)
		// if err != nil {
		// 	return err
		// }
		// nb.ChanExecDB <- DBsync{Unread: 0, ArticleUrls: []string{article.Url}}
	case *TagsPage:
		feeds, err := nb.QueryTagFeeds(id)
		if err != nil {
			return err
		}
		err = nb.Push(&feeds)
		if err != nil {
			return err
		}
	case *TagFeeds:
		feed, err := nb.QueryFeed(id)
		if err != nil {
			return err
		}
		err = nb.Push(&feed)
		if err != nil {
			return err
		}
	case *Article:
		return nil
	}
	return nil
}

func (nb *Nvimboat) Back() error {
	switch nb.Pages.Top().(type) {
	case *MainMenu:
		return nil
	default:
		nb.Pop()
	}
	return nil
}

func (nb *Nvimboat) NextUnread() error {
	return nil
}

func (nb *Nvimboat) PrevUnread() error {
	return nil
}

func (nb *Nvimboat) NextArticle() error {
	return nil
}

func (nb *Nvimboat) PrevArticle() error {
	return nil
}

func (nb *Nvimboat) ToggleArticleRead(urls ...string) error {
	var (
		err  error
		sync DBsync
	)
	if urls[0] == "Article" {
		article := nb.Pages.Top().(*Article)
		nb.Pages.Pop()
		nb.Log(nb.Pages.Top())
		nb.ToggleArticleRead(article.Url)
		idx, err := nb.Pages.Top().SubPageIdx(article)
		if err != nil {
			return err
		}
		switch page := nb.Pages.Top().(type) {
		case *Filter:
			page.Articles[idx].Unread = 1
		case *Feed:
			page.Articles[idx].Unread = 1
		}
		nb.Show(nb.Pages.Top())
	}
	anyUnread, err := nb.anyArticleUnread(urls...)
	if err != nil {
		return err
	}
	if anyUnread {
		sync.Unread = 0
	} else {
		sync.Unread = 1
	}
	switch nb.Pages.Top().(type) {
	case *Filter:
		sync.ArticleUrls = urls
	case *Feed:
		sync.ArticleUrls = urls
	}
	nb.ChanExecDB <- sync
	return err
}

var Actions = []string{
	"enable",
	"disable",
	"show-main",
	"show-tags",
	"select",
	"back",
	"next-unread",
	"prev-unread",
	"next-article",
	"prev-article",
	"toggle-read",
}
