package nvimboat

import (
	"errors"
	"fmt"
	"log"
	"os"

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
	nb.batch.SetBufferOption(*nb.buffer, "filetype", "nvimboat")
	nb.batch.SetBufferOption(*nb.buffer, "buftype", "nofile")
	nb.batch.SetWindowOption(*nb.window, "wrap", false)
	if err != nil {
		return err
	}

	return nil
}

func (nb *Nvimboat) Log(val ...any) {
	log.Println(val...)
	nb.plugin.Nvim.Exec(fmt.Sprintf(`echo "%v"`, val), false)
}

func (nb *Nvimboat) Push(p Page) error {
	lines, err := p.Render()
	nb.SetLines(lines)
	nb.PageStack.Push(p)
	err = nb.setPageType(p)
	return err
}

func (nb *Nvimboat) Pop() error {
	nb.PageStack.Pop()
	lines, err := nb.PageStack.top.Render()
	nb.SetLines(lines)
	err = nb.setPageType(nb.PageStack.top)
	return err
}

func (nb *Nvimboat) setupLogging() {
	var err error

	nb.LogFile, err = os.OpenFile(nb.Config["log"].(string), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Println(err)
	}

	log.SetOutput(nb.LogFile)
	log.SetFlags(0)
}

func (nb *Nvimboat) setPageType(p Page) error {
	t := pageTypeString(p)
	err := nb.plugin.Nvim.ExecLua(nvimboatSetPageType, new(any), t)
	if err != nil {
		return err
	}
	return nil
}

func (nb *Nvimboat) PageType() (any, error) {
	var page_type any
	err := nb.plugin.Nvim.ExecLua(nvimboatPage, &page_type)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Can't get page type: %v", err))
	}
	return page_type, nil
}

func (nb *Nvimboat) SetLines(lines []string) error {
	err := nb.plugin.Nvim.SetBufferLines(*nb.buffer, 0, -1, false, strings2bytes(lines))
	if err != nil {
		return err
	}

	return nil

}
