package nvimboat

import (
	"fmt"
	"strings"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type (
	Page interface {
		Select(dbh rssdb.DbHandle, id string) (p Page, err error)
		Render(nv *nvim.Nvim, buf nvim.Buffer) (err error)
		ChildIdx(p Page) (idx int, err error)
		Back() (cursor_x int, err error)
		ToggleRead(dbh rssdb.DbHandle, ids []string) (err error)
	}
	PageStack []Page
)

func (ps *PageStack) Show() (err error) {
	err = setLines(Nvim, *NvBuffer, []string{""})
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Show: %w\n", err)
		return
	}
	defer trimTrail(Nvim, *NvBuffer)
	defer Nvim.SetWindowCursor(*NvWindow, [2]int{0, 1})
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Show: %w\n", err)
		return
	}
	ps.Top().Render(Nvim, *NvBuffer)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.Show: %w\n", err)
		return
	}
	return
}

func (ps *PageStack) Top() (p Page) {
	if pageCount := len(*ps); pageCount > 0 {
		return (*ps)[pageCount-1]
	}
	return nil
}

func (ps *PageStack) Push(p Page, id string) (err error) {
	*ps = append(*ps, p)
	pageType := fmt.Sprintf("%T", p)
	_, pageType, _ = strings.Cut(pageType, "nvimboat.")
	err = Nvim.ExecLua(luaPushPage, new(any), pageType, id)
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ShowMain: %w\n", err)
		return
	}
	return
}

func (ps *PageStack) Pop() (p Page, err error) {
	if len(*ps) > 1 {
		*ps = (*ps)[:len(*ps)-1]
	}
	err = Nvim.ExecLua(luaPopPage, new(any))
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.PopPage: %w\n", err)
		return
	}
	return ps.Top(), nil
}

func (ps *PageStack) Reset() (err error) {
	currentPages := *ps // Save current state in case of error
	*ps = []Page{}
	err = Nvim.ExecLua(luaResetPages, new(any))
	if err != nil {
		err = fmt.Errorf("nvimboat/Nvimboat.ResetPages: %w\n", err)
		*ps = currentPages
		return
	}
	return
}
