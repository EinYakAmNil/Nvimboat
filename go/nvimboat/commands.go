package nvimboat

import (
	"errors"
	"log"
)

func (nb *Nvimboat) Command(args []string) error {
	err := nb.batch.Execute()
	if err != nil {
		log.Println(err)
		return err
	}
	if nb.LogFile == nil {
		nb.setupLogging()
	}
	nb.Log(args)
	action := args[0]
	switch action {
	case "enable":
		err = nb.Enable()
	case "disable":
		err = nb.Disable()
	case "show-main":
		err = nb.ShowMain()
	case "show-tags":
		err = nb.ShowTags()
	case "select":
		if len(args) > 1 {
			err = nb.Select(args[1])
			return nil
		}
		return errors.New("No arguments for select command.")
	case "back":
		err = nb.Back()
	}
	if err != nil {
		nb.Log(err)
		return err
	}

	// pt, err := nb.PageType()
	// if err != nil {
	// 	return err
	// }
	// nb.Log(pt)
	return nil
}

func (nb *Nvimboat) Enable() error {
	var err error
	// nb.Feeds, err = nb.QueryFeeds()
	err = nb.plugin.Nvim.ExecLua(nvimboatEnable, new(any))
	if err != nil {
		return err
	}
	err = nb.Push(&MainMenu{})
	if err != nil {
		return err
	}

	return nil
}

func (nb *Nvimboat) Disable() error {
	err := nb.plugin.Nvim.ExecLua(nvimboatDisable, new(any))
	if err != nil {
		return err
	}

	return nil
}

func (nb *Nvimboat) Select(id string) error {
	switch nb.PageStack.top.(type) {
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
			if err != nil {
				return err
			}
			err = nb.Push(&filter)
			if err != nil {
				return err
			}
		}
	case *Filter:
		article, err := nb.QueryArticle(id)
		if err != nil {
			return err
		}
		nb.Push(&article)
		if err != nil {
			return err
		}
	case *Feed:
		article, err := nb.QueryArticle(id)
		if err != nil {
			return err
		}
		nb.Push(&article)
		if err != nil {
			return err
		}
	case *TagsPage:
		feeds, err := nb.QueryTagFeeds(id)
		if err != nil {
			return err
		}
		nb.Push(&feeds)
		if err != nil {
			return err
		}
	case *TagFeeds:
		article, err := nb.QueryArticle(id)
		if err != nil {
			return err
		}
		nb.Push(&article)
		if err != nil {
			return err
		}
	case *Article:
		return nil
	}
	return nil
}

func (nb *Nvimboat) Back() error {
	switch nb.PageStack.top.(type) {
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

func (nb *Nvimboat) ShowMain() error {
	switch nb.PageStack.top.(type) {
	case *MainMenu:
		return nil
	default:
		nb.PageStack.Pages = nb.PageStack.Pages[:1]
		nb.PageStack.top = nb.PageStack.Pages[0]
		lines, err := nb.PageStack.top.Render()
		if err != nil {
			return err
		}
		err = nb.SetLines(lines)
		if err != nil {
			return err
		}
		nb.setPageType(nb.PageStack.top)
	}

	return nil
}

func (nb *Nvimboat) ShowTags() error {
	nb.Log("Showing tags.")
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
