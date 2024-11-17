package rssdb

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "embed"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed sql/schema.sql
var schemaSql string

type DbHandle struct {
	DB      *sql.DB
	Ctx     context.Context
	Queries *Queries
}

func ConnectDb(dbPath string) (dbh DbHandle, err error) {
	dbh.Ctx = context.Background()
	dbh.DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		err = fmt.Errorf("ConnectDb: %w", err)
		return
	}
	// only create tables, if the database does not exist yet
	if _, noDbErr := os.Stat(dbPath); noDbErr != nil {
		if _, err = dbh.DB.ExecContext(dbh.Ctx, schemaSql); err != nil {
			err = fmt.Errorf("ConnectDb: %w", err)
			return
		}
	}
	dbh.Queries = New(dbh.DB)
	return
}
