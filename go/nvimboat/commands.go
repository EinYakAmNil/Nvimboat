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
	err = nb.SyncState(Main{})
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
	case Main:
		if id[:4] == "http" {
			feed, err := nb.QueryFeed(id)
			if err != nil {
				return err
			}
			err = nb.SyncState(feed)
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
			err = nb.SyncState(filter)
			if err != nil {
				return err
			}
		}
	case Filter:
		article, err := nb.QueryArticle(id)
		if err != nil {
			return err
		}
		nb.SyncState(article)
		if err != nil {
			return err
		}
	case Feed:
		article, err := nb.QueryArticle(id)
		if err != nil {
			return err
		}
		nb.SyncState(article)
		if err != nil {
			return err
		}
	case TagsPage:
		tags, err := nb.QueryTags()
		if err != nil {
			return err
		}
		nb.SyncState(tags)
		if err != nil {
			return err
		}
	case TagFeeds:
		feeds, err := nb.QueryTagFeeds(id)
		if err != nil {
			return err
		}
		nb.SyncState(feeds)
		if err != nil {
			return err
		}
	}
	nb.Log(nb.PageStack.top)
	return nil
}

func (nb *Nvimboat) Back() error {
	switch nb.PageStack.top.(type) {
	case Main:
		return nil
	default:
		nb.PageStack.Pop()
		nb.setPageType(nb.PageStack.top)
		nb.Log(nb.PageStack.top)
	}

	return nil
}
