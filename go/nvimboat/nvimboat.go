package nvimboat

import (
	"database/sql"
	"os"

	"github.com/neovim/go-client/nvim"
)

func (nb *Nvimboat) Push(newPage Page) (err error) {
	err = nb.Show(newPage)
	if err != nil {
		return
	}
	nb.Pages.Push(newPage)
	return
}

func (nb *Nvimboat) Pop() (err error) {
	currentPage := nb.Pages.Top()
	nb.Pages.Pop()
	pos, err := nb.Pages.Top().ChildIdx(currentPage)
	if err != nil {
		return
	}
	page, err := nb.Pages.Top().QuerySelf(nb.DBHandler)
	if err != nil {
		return
	}
	err = nb.Show(page)
	if err != nil {
		return
	}
	err = nb.Nvim.SetWindowCursor(*nb.Window, [2]int{pos + 1, 0})
	return
}

func (nb *Nvimboat) Show(newPage Page) (err error) {
	defer trimTrail(nb.Nvim, *nb.Buffer)
	err = setLines(nb.Nvim, *nb.Buffer, []string{""})
	if err != nil {
		return
	}
	err = newPage.Render(nb.Nvim, *nb.Buffer, nb.UnreadOnly, nb.Config["separator"].(string))
	if err != nil {
		return
	}
	err = nb.setPageType(newPage)
	return
}

func (nb *Nvimboat) init(nv *nvim.Nvim) (err error) {
	nb.Nvim = nv
	nb.SyncDBchan = make(chan SyncDB)
	nb.Config = make(map[string]any)
	nb.Window = new(nvim.Window)
	nb.Buffer = new(nvim.Buffer)
	nb.UnreadOnly = false
	execBatch := nv.NewBatch()
	execBatch.CurrentWindow(nb.Window)
	execBatch.CurrentBuffer(nb.Buffer)
	execBatch.ExecLua(nvimboatConfig, &nb.Config)
	execBatch.ExecLua(nvimboatFeeds, &nb.Feeds)
	execBatch.ExecLua(nvimboatFilters, &nb.Filters)
	execBatch.SetBufferOption(*nb.Buffer, "filetype", "nvimboat")
	execBatch.SetBufferOption(*nb.Buffer, "buftype", "nofile")
	execBatch.SetWindowOption(*nb.Window, "wrap", false)
	err = execBatch.Execute()
	if err != nil {
		return
	}
	nb.DBHandler, err = InitDB(nb.Config["dbpath"].(string))
	if err != nil {
		return
	}
	err = SetupLogging(nb.Config["log"].(string))
	return
}

func (nb *Nvimboat) setPageType(p Page) (err error) {
	t := pageTypeString(p)
	err = nb.Nvim.ExecLua(nvimboatSetPageType, new(any), t)
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
