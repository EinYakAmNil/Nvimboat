package nvimboat

import (
	"errors"
	"fmt"
	"strings"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	"github.com/neovim/go-client/nvim"
)

type (
	Page interface {
		ID() string
		Select(dbh rssdb.DbHandle, id string) (p Page, err error)
		Open(urls ...string) (err error)
		Render(nv *nvim.Nvim, buf nvim.Buffer) (err error)
		ChildIdx(p Page) (idx int, err error)
		Back() (cursor_x int, err error)
		ToggleRead(dbh rssdb.DbHandle, ids []string) (pos [2][2]int, err error)
		NextUnread(dbh rssdb.DbHandle) (err error)
		PrevUnread(dbh rssdb.DbHandle) (err error)
		Delete(dbh rssdb.DbHandle, ids []string) (err error)
	}
	PageStack []Page
)

func (ps *PageStack) Show() (err error) {
	err = setLines(Nvim, *NvBuffer, []string{""})
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/PageStack.Show"))
		return
	}
	defer trimTrail(Nvim, *NvBuffer)
	if _, ok := Pages.Top().(*Article); ok {
		err = Nvim.SetWindowOption(*NvWindow, "wrap", true)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/PageStack.Show"))
			return
		}
	} else {
		err = Nvim.SetWindowOption(*NvWindow, "wrap", false)
		if err != nil {
			err = errors.Join(err, errors.New("nvimboat/PageStack.Show"))
			return
		}
	}
	err = ps.Top().Render(Nvim, *NvBuffer)
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/PageStack.Show"))
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
		err = errors.Join(err, errors.New("nvimboat/PageStack.Push"))
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
		err = errors.Join(err, errors.New("nvimboat/PageStack.Pop"))
		return
	}
	return ps.Top(), nil
}

func (ps *PageStack) Reset() (err error) {
	currentPages := *ps // Save current state in case of error
	*ps = []Page{}
	err = Nvim.ExecLua(luaResetPages, new(any))
	if err != nil {
		err = errors.Join(err, errors.New("nvimboat/PageStack.Reset"))
		*ps = currentPages
		return
	}
	return
}
