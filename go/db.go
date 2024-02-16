package main

import (
	"database/sql"
	"fmt"
	"nvimboat"
	"strings"
	"time"
)

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
	case exec, ok := <-nb.SyncDBchan:
		if ok {
			if len(exec.ArticleUrls) > 0 {
				articleReadState(nb.DBHandler, exec.Unread, exec.ArticleUrls...)
			}
			if len(exec.FeedUrls) > 0 {
				articleReadState(nb.DBHandler, exec.Unread, exec.FeedUrls...)
			}
		}
	default:
		return fmt.Errorf("channel closed")
	}
	return nil
}

func articleReadUpdate(count int) string {
	if count == 0 {
		return ""
	}
	update := `UPDATE rss_item SET unread = ? WHERE url IN (`
	qmarks := make([]string, count)
	for i := range qmarks {
		qmarks[i] = "?"
	}
	articleCount := strings.Join(qmarks, ", ")

	return update + articleCount + `)`
}

func articleReadState(db *sql.DB, read int, url ...string) error {
	sqlArgs := []any{read}
	for _, u := range url {
		sqlArgs = append(sqlArgs, u)
	}
	update := articleReadUpdate(len(url))
	_, err := db.Exec(update, sqlArgs...)
	if err != nil {
		return fmt.Errorf("ArticleReadState -> db.open: " + fmt.Sprintln(err))
	}
	return nil
}
