package nvimboat

import (
	"fmt"
	"log"

	"github.com/neovim/go-client/nvim"
	nvimPlugin "github.com/neovim/go-client/nvim/plugin"
)

func (nb *Nvimboat) Init(p *nvimPlugin.Plugin) error {
	var err error
	nb.plugin = p
	nb.batch = p.Nvim.NewBatch()
	nb.buffer = new(nvim.Buffer)
	nb.window = new(nvim.Window)
	nb.batch.CurrentBuffer(nb.buffer)
	nb.batch.CurrentWindow(nb.window)
	nb.Config = make(map[string]any)
	nb.batch.ExecLua(nvimboatConfig, &nb.Config)
	nb.batch.ExecLua(nvimboatFeeds, &nb.ConfigFeeds)
	nb.batch.ExecLua(nvimboatFilters, &nb.ConfigFilters)
	nb.batch.SetBufferOption(*nb.buffer, "filetype", "nvimboat")
	nb.batch.SetBufferOption(*nb.buffer, "buftype", "nofile")
	nb.batch.SetWindowOption(*nb.window, "wrap", false)
	if err != nil {
		return err
	}
	return nil
}

func (nb *Nvimboat) Show(p Page) error {
	nb.SetLines([]string{})
	defer nb.trimTrail()

	switch p.(type) {
	case *Article:
		lines, err := p.Render()
		if err != nil {
			return err
		}
		nb.SetLines(lines[0])
	case *TagsPage:
		lines, err := p.Render()
		if err != nil {
			return err
		}
		nb.SetLines(lines[0])
	default:
		cols, err := p.Render()
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
	return nil
}

func (nb *Nvimboat) Log(val ...any) {
	log.Println(val...)
	nb.plugin.Nvim.Exec(fmt.Sprintf(`echo "%v"`, val), false)
}

func (nb *Nvimboat) Push(p Page) error {
	err := nb.Show(p)
	if err != nil {
		return err
	}
	err = nb.setPageType(p)
	nb.PageStack.Push(p)
	return err
}

func (nb *Nvimboat) Pop() error {
	var err error
	oldPage := nb.PageStack.top
	nb.PageStack.Pop()
	switch p := nb.PageStack.top.(type) {
	case *MainMenu:
		nb.PageStack.top, err = nb.showMain()
		if err != nil {
			return err
		}
	case *Feed:
		newPage, err := nb.QueryFeed(p.RssUrl)
		if err != nil {
			return err
		}
		nb.PageStack.top = &newPage
	case *TagFeeds:
		newPage, err := nb.QueryTagFeeds(p.Tag)
		if err != nil {
			return err
		}
		nb.PageStack.top = &newPage
	}
	pos, err := nb.PageStack.top.ElementIdx(oldPage)
	if err != nil {
		return err
	}
	err = nb.Show(nb.PageStack.top)
	if err != nil {
		return err
	}
	err = nb.setPageType(nb.PageStack.top)
	if err != nil {
		return err
	}
	err = nb.plugin.Nvim.SetWindowCursor(*nb.window, [2]int{pos + 1, 0})
	if err != nil {
		return err
	}
	return nil
}

func (nb *Nvimboat) SetLines(lines []string) error {
	err := nb.plugin.Nvim.SetBufferLines(*nb.buffer, 0, -1, false, strings2bytes(lines))
	if err != nil {
		return err
	}
	return nil
}
