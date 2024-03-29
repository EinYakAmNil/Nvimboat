package nvimboat

import (
	"fmt"

	"github.com/neovim/go-client/nvim"
)

func (nb *Nvimboat) ShowMain(nv *nvim.Nvim, args ...string) (err error) {
	mainMenu, err := QueryMain(nb.DBHandler, nb.Feeds, nb.Filters)
	if err != nil {
		return
	}
	nb.Pages.Pages = nil
	err = nb.Push(mainMenu)
	return
}

func (nb *Nvimboat) Enable(nv *nvim.Nvim, args ...string) (err error) {
	err = nb.init(nv)
	if err != nil {
		return
	}
	err = nv.ExecLua(nvimboatEnable, new(any))
	if err != nil {
		return
	}
	err = nb.ShowMain(nv, args...)
	return
}

func (nb *Nvimboat) Disable(nv *nvim.Nvim, args ...string) (err error) {
	err = nv.ExecLua(nvimboatDisable, new(any))
	return
}

func (nb *Nvimboat) ShowTags(nv *nvim.Nvim, args ...string) (err error) {
	tags, err := QueryTags(nb.Feeds)
	if err != nil {
		return
	}
	err = nb.Push(tags)
	return
}

func (nb *Nvimboat) Select(nv *nvim.Nvim, args ...string) (err error) {
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments to call 'select'")
	}
	defer nb.Nvim.SetWindowCursor(*nb.Window, [2]int{0, 1})
	err = nb.Pages.Top().Select(nb, args[1])
	return
}

func (nb *Nvimboat) Back(nv *nvim.Nvim, args ...string) (err error) {
	switch nb.Pages.Top().(type) {
	case *MainMenu:
		return
	default:
		err = nb.Pop()
		return
	}
}

func (nb *Nvimboat) NextUnread(nv *nvim.Nvim, args ...string) error {
	if article, isArticle := nb.Pages.Top().(*Article); isArticle {
		nb.Pages.Pop()
		if top, isCollection := nb.Pages.Top().(ArticlesPage); isCollection {
			newArticle, err := top.FindUnread("next", *article)
			if err != nil {
				return err
			}
			err = top.(Page).Select(nb, newArticle.Url)
			return err
		} else {
			nb.Pages.Push(article)
			return fmt.Errorf("next unread not implemented for: %v", top)
		}
	}
	return fmt.Errorf("not inside an article")
}

func (nb *Nvimboat) PrevUnread(nv *nvim.Nvim, args ...string) error {
	if article, isArticle := nb.Pages.Top().(*Article); isArticle {
		nb.Pages.Pop()
		if top, isCollection := nb.Pages.Top().(ArticlesPage); isCollection {
			newArticle, err := top.FindUnread("prev", *article)
			if err != nil {
				return err
			}
			err = top.(Page).Select(nb, newArticle.Url)
			return err
		} else {
			nb.Pages.Push(article)
			return fmt.Errorf("next unread not implemented for: %v", top)
		}
	}
	return fmt.Errorf("not inside an article")
}

func (nb *Nvimboat) NextArticle(nv *nvim.Nvim, args ...string) error {
	if article, isArticle := nb.Pages.Top().(*Article); isArticle {
		nb.Pages.Pop()
		top := nb.Pages.Top()
		idx, err := top.ChildIdx(article)
		switch feed := top.(type) {
		case *Feed:
			if idx+1 < len(feed.Articles) {
				err = feed.Select(nb, feed.Articles[idx+1].Url)
				return err
			} else {
				err = feed.Select(nb, feed.Articles[idx].Url)
				return nil
			}
		case *Filter:
			if idx+1 < len(feed.Articles) {
				err = feed.Select(nb, feed.Articles[idx+1].Url)
				return err
			} else {
				err = feed.Select(nb, feed.Articles[idx].Url)
				return nil
			}
		}
	}
	return fmt.Errorf("not inside an article")
}

func (nb *Nvimboat) PrevArticle(nv *nvim.Nvim, args ...string) error {
	if article, isArticle := nb.Pages.Top().(*Article); isArticle {
		nb.Pages.Pop()
		top := nb.Pages.Top()
		idx, err := top.ChildIdx(article)
		switch feed := top.(type) {
		case *Feed:
			if idx-1 >= 0 {
				err = feed.Select(nb, feed.Articles[idx-1].Url)
				return err
			} else {
				err = feed.Select(nb, feed.Articles[idx].Url)
				return nil
			}
		case *Filter:
			if idx-1 >= 0 {
				err = feed.Select(nb, feed.Articles[idx-1].Url)
				return err
			} else {
				err = feed.Select(nb, feed.Articles[idx].Url)
				return nil
			}
		}
	}
	return fmt.Errorf("not inside an article")
}

func (nb *Nvimboat) ToggleArticleRead(nv *nvim.Nvim, args ...string) (err error) {
	defer trimTrail(nb.Nvim, *nb.Buffer)
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments to call 'toggle-unread'")
	}
	err = nb.Pages.Top().ToggleUnread(nb, args[1:]...)
	return
}

func (nb *Nvimboat) Delete(nv *nvim.Nvim, args ...string) (err error) {
	defer trimTrail(nb.Nvim, *nb.Buffer)
	if len(args) < 2 {
		return fmt.Errorf("not enough arguments to call 'toggle-unread'")
	}
	err = nb.Pages.Top().Delete(nb, args[1:]...)
	return
}
