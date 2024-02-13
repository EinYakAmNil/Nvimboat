package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"nvimboat"

	"github.com/neovim/go-client/nvim"
	nvimPlugin "github.com/neovim/go-client/nvim/plugin"
)

func main() {
	chanNvimboat := make(chan *nvimboat.Nvimboat)
	go nvimboatLoop(chanNvimboat)
	nb := <-chanNvimboat
	unreadUpdate(nb)
}

func nvimboatLoop(cnb chan *nvimboat.Nvimboat) {
	nb := new(nvimboat.Nvimboat)
	defer nb.LogFile.Close()
	if nb.DB != nil {
		defer nb.DB.Close()
	}
	cnb <- nb
	nb.ChanExecDB = make(chan nvimboat.DBsync)
	nvimPlugin.Main(func(p *nvimPlugin.Plugin) error {
		nb.Prepare(p)
		p.HandleCommand(&nvimPlugin.CommandOptions{Name: "Nvimboat", NArgs: "+", Complete: "customlist,CompleteNvimboat"}, nb.Command)
		p.HandleFunction(&nvimPlugin.FunctionOptions{Name: "CompleteNvimboat"}, func(c *nvim.CommandCompletionArgs) ([]string, error) {
			var suggestions []string

			if c.ArgLead != "" {
				for _, a := range nvimboat.Actions {
					lcd := min(len(a), len(c.ArgLead))
					if c.ArgLead[:lcd] == a[:lcd] {
						suggestions = append(suggestions, a)
					}
				}
				return suggestions, nil
			}
			return nvimboat.Actions, nil
		})
		return nil
	})
}

func unreadUpdate(nb *nvimboat.Nvimboat) {
	for {
		err := handleExec(nb)
		if err != nil {
			time.Sleep(time.Millisecond)
		}
	}
}

func handleExec(nb *nvimboat.Nvimboat) error {
	select {
	case exec, ok := <-nb.ChanExecDB:
		if ok {
			if len(exec.ArticleUrls) > 0 {
				articleReadState(nb.DB, exec.Unread, exec.ArticleUrls...)
			}
			if len(exec.FeedUrls) > 0 {
				articleReadState(nb.DB, exec.Unread, exec.FeedUrls...)
			}
		}
	default:
		return fmt.Errorf("channel closed")
	}
	return nil
}

func articleReadUpdate(n int) string {
	if n == 0 {
		return ""
	}
	const (
		prefix = `UPDATE rss_item SET unread = ? WHERE url IN (?`
		suffix = `)`
	)
	if n < 2 {
		return prefix + suffix
	}
	articleCount := strings.Repeat(", ?", n-1)

	return prefix + articleCount + suffix
}

func articleReadState(db *sql.DB, read int, url ...string) error {
	sqlArgs := []any{read}
	for _, u := range url {
		sqlArgs = append(sqlArgs, u)
	}
	update := articleReadUpdate(len(url))
	_, err := db.Exec(update, sqlArgs...)
	if err != nil {
		return errors.New("ArticleReadState -> db.open: " + fmt.Sprintln(err))
	}
	return nil
}
