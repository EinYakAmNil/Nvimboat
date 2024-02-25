package main

import (
	"database/sql"
	"fmt"
	"log"
	"nvimboat"
	"strings"
	"time"
)

func dbUpdate(nb *nvimboat.Nvimboat) {
	for {
		err := readDBSyncChan(nb)
		if err != nil {
			time.Sleep(time.Millisecond)
		}
	}
}

func readDBSyncChan(nb *nvimboat.Nvimboat) error {
	select {
	case exec, ok := <-nb.SyncDBchan:
		if ok {
			if len(exec.ArticleUrls) > 0 {
				if exec.Delete {
					deleteArticle(nb.DBHandler, exec.ArticleUrls...)
				}
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

func deleteArticle(db *sql.DB, urls ...string) (err error) {
	log.Println("Delete", urls)
	var (
		deleteStmt = `UPDATE rss_item SET deleted = 1 WHERE url IN `
		sqlArgs    []any
	)
	for _, u := range urls {
		sqlArgs = append(sqlArgs, u)
	}
	placeholder := `(` + strings.Repeat("?, ", len(urls)) + `)`
	deleteStmt += placeholder
	_, err = db.Exec(deleteStmt, sqlArgs...)
	if err != nil {
		return fmt.Errorf("error deleting articles %v:\n%v\n", urls, err)
	}
	return
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
