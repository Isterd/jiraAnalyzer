package database

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type JiraPostgres struct {
	db *sqlx.DB
}

func NewJiraPostgres(db *sqlx.DB) *JiraPostgres {
	return &JiraPostgres{db: db}
}

func (r *JiraPostgres) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
