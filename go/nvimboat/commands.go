package nvimboat

import "fmt"

func (nb *Nvimboat) Command(args []string) error {
	err := nb.init()
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
	case "next-unread":
		err = nb.NextUnread()
	case "prev-unread":
		err = nb.PrevUnread()
	case "next-article":
		err = nb.NextArticle()
	case "prev-article":
		err = nb.PrevArticle()
	case "toggle-read":
		err = nb.ToggleArticleRead(args[1:]...)
	case "back":
		err = nb.Back()
	case "select":
		if len(args) > 1 {
			err = nb.Select(args[1])
			return nil
		}
		return fmt.Errorf("No arguments for select command.")
	default:
		nb.Log("command not yet implemented: ", args)
	}
	if err != nil {
		nb.Log(err)
		return err
	}
	return nil
}

func (nb *Nvimboat) Push(p Page) error {
	err := nb.Show(p)
	if err != nil {
		return err
	}
	nb.Pages.Push(p)
	return err
}

func (nb *Nvimboat) Pop() error {
	currentPage := nb.Pages.Top()
	nb.Pages.Pop()
	pos, err := nb.Pages.Top().SubPageIdx(currentPage)
	if err != nil {
		return err
	}
	nb.Show(nb.Pages.Top())
	err = nb.Nvim.Plugin.Nvim.SetWindowCursor(*nb.Nvim.Window, [2]int{pos + 1, 0})
	return nil
}

func (nb *Nvimboat) Show(p Page) error {
	nb.SetLines([]string{})
	defer nb.trimTrail()

	switch p.(type) {
	case *Article:
		lines, err := p.Render(false)
		if err != nil {
			return err
		}
		nb.SetLines(lines[0])
	case *TagsPage:
		lines, err := p.Render(false)
		if err != nil {
			return err
		}
		nb.SetLines(lines[0])
	default:
		cols, err := p.Render(false)
		if err != nil {
			return err
		}
		for _, c := range cols {
			err = nb.addColumn(c, nb.Config["separator"].(string))
			if err != nil {
				return err
			}
		}
	}
	nb.setPageType(p)
	return nil
}

func (nb *Nvimboat) SetLines(lines []string) error {
	err := nb.Nvim.Plugin.Nvim.SetBufferLines(*nb.Nvim.Buffer, 0, -1, false, strings2bytes(lines))
	if err != nil {
		return err
	}
	return nil
}
