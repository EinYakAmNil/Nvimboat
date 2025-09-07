package rssdb

import (
	"context"
	"fmt"
)

const queryFilterOld = `-- name: QueryFilter :many
SELECT unread, pubDate, author, title, url FROM rss_item
WHERE %s
ORDER BY pubDate DESC
`

type QueryFilterRowOld struct {
	Unread  int
	Pubdate int64
	Author  string
	Title   string
	Url     string
}

func (q *Queries) QueryFilterOld(ctx context.Context, query string) ([]QueryFilterRowOld, error) {
	rows, err := q.db.QueryContext(ctx, fmt.Sprintf(queryFilterOld, query))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []QueryFilterRowOld
	for rows.Next() {
		var i QueryFilterRowOld
		if err := rows.Scan(
			&i.Unread,
			&i.Pubdate,
			&i.Author,
			&i.Title,
			&i.Url,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
