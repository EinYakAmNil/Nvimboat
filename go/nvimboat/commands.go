package nvimboat

import "log"

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
		err = nb.Select(args[1])
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
	nb.PageStack.top = Main{}
	err = nb.SetPageType("Main")
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

func (nb *Nvimboat) Select(item string) error {
	var err error
	switch nb.PageStack.top.(type) {
	case Main:
		nb.PageStack.top = Feed{}
		err = nb.SetPageType("Feed")
		if err != nil {
			return err
		}
	}
	return nil
}

func (nb *Nvimboat) Back() error {
	var err error
	switch nb.PageStack.top.(type) {
	case Main:
	case Feed:
		nb.PageStack.top = Main{}
		err = nb.SetPageType("Main")
		if err != nil {
			return err
		}
	}
	return nil
}
