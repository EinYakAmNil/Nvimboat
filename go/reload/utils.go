package reload

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EinYakAmNil/Nvimboat/go/engine/rssdb"
	_ "github.com/mattn/go-sqlite3"
)
func ConnectDb(dbPath string) (queries *rssdb.Queries, ctx context.Context, err error) {
	ctx = context.Background()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		err = fmt.Errorf("ConnectDb: %w", err)
		return
	}
	queries = rssdb.New(db)
	return
}
