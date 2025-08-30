package rssdb

import (
	"context"
	"fmt"
)

const queryFilter = `-- name: QueryFilter :many
SELECT unread, pubDate, author, title, url FROM rss_item
WHERE %s
ORDER BY pubDate DESC
`

type QueryFilterRow struct {
	Unread  int
	Pubdate int64
	Author  string
	Title   string
	Url     string
}

func (q *Queries) QueryFilterOld(ctx context.Context, query string) ([]QueryFilterRow, error) {
	rows, err := q.db.QueryContext(ctx, fmt.Sprintf(queryFilter, query))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []QueryFilterRow
	for rows.Next() {
		var i QueryFilterRow
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
