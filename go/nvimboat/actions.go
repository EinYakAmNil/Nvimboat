package nvimboat

func (nb *Nvimboat) ShowMain() error {
	mainmenu, err := QueryMain(nb.DB, nb.ConfigFeeds, nb.ConfigFilters)
	if err != nil {
		return err
	}
	nb.Pages.Pages = nil
	err = nb.Push(mainmenu)
	return err
}

func (nb *Nvimboat) Enable() error {
	err := nb.ShowMain()
	if err != nil {
		return err
	}
	err = nb.Nvim.Plugin.Nvim.ExecLua(nvimboatEnable, new(any))
	return err
}

func (nb *Nvimboat) Disable() error {
	err := nb.Nvim.Plugin.Nvim.ExecLua(nvimboatDisable, new(any))
	return err
}

func (nb *Nvimboat) ShowTags() error {
	tags, err := QueryTags(nb.ConfigFeeds)
	if err != nil {
		return err
	}
	nb.Push(tags)
	if err != nil {
		return err
	}
	return nil
}

func (nb *Nvimboat) Select(id string) error {
	defer nb.Nvim.Plugin.Nvim.SetWindowCursor(*nb.Nvim.Window, [2]int{0, 1})
	page, err := nb.Pages.Top().QuerySelect(nb.DB, id)
	if err != nil {
		return err
	}
	err = nb.Push(page)
	return err
}

func (nb *Nvimboat) Back() error {
	switch nb.Pages.Top().(type) {
	case *MainMenu:
		return nil
	default:
		err := nb.Pop()
		if err != nil {
			return err
		}
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
		nb.ToggleArticleRead(article.Url)
		idx, err := nb.Pages.Top().SubPageIdx(article)
		if err != nil {
			return err
		}
		switch page := nb.Pages.Top().(type) {
		case *Filter:
			page.Articles[idx].Unread = 1
			page.updateUnreadCount()
		case *Feed:
			page.Articles[idx].Unread = 1
			page.updateUnreadCount()
		}
		err = nb.Show(nb.Pages.Top())
		return err
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
