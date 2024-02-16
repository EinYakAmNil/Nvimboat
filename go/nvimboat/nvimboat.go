package nvimboat

import (
	"database/sql"
	"os"

	"github.com/neovim/go-client/nvim"
)

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
	pos, err := nb.Pages.Top().ChildIdx(currentPage)
	if err != nil {
		return err
	}
	page, err := nb.Pages.Top().QuerySelf(nb.DBHandler)
	if err != nil {
		return err
	}
	err = nb.Show(page)
	if err != nil {
		return err
	}
	err = nb.Nvim.SetWindowCursor(*nb.Window, [2]int{pos + 1, 0})
	return err
}

func (nb *Nvimboat) Show(page Page) (err error) {
	defer trimTrail(nb.Nvim, *nb.Buffer)
	err = setLines(nb.Nvim, *nb.Buffer, []string{""})
	if err != nil {
		return
	}
	err = page.Render(nb.Nvim, *nb.Buffer, nb.UnreadOnly, nb.Config["separator"].(string))
	if err != nil {
		return
	}
	switch p := page.(type) {
	case *Article:
		nb.SyncDBchan <- SyncDB{Unread: 0, ArticleUrls: []string{p.Url}}
	}
	err = nb.setPageType(page)
	return
}

type (
	Nvimboat struct {
		Config     map[string]any
		Pages      PageStack
		Feeds      []map[string]any
		Filters    []map[string]any
		LogFile    *os.File
		DBHandler  *sql.DB
		SyncDBchan chan SyncDB
		Nvim       *nvim.Nvim
		Window     *nvim.Window
		Buffer     *nvim.Buffer
		UnreadOnly bool
	}
	Action func(*Nvimboat, *nvim.Nvim, ...string) error
)
