package repository

import (
	"context"
	"database/sql"

	"proxysql-galera-app/internal/model"
	"proxysql-galera-app/internal/pool"
)

type Writer interface {
	InitSchema(ctx context.Context) error
	CreateTodo(ctx context.Context, t *model.Todo) error
}

type MySQLWriter struct {
	pool   *pool.RouterPool
	dbname string
}

func NewMySQLWriter(pool *pool.RouterPool, dbname string) *MySQLWriter {
	return &MySQLWriter{pool: pool, dbname: dbname}
}

func (m *MySQLWriter) InitSchema(ctx context.Context) error {
	const ddl = `CREATE TABLE IF NOT EXISTS todos (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR(255) NOT NULL,
		created_at TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
	)`
	return m.pool.WithConn(ctx, func(ctx context.Context, db *sql.DB) error { _, err := db.ExecContext(ctx, ddl); return err })
}

func (m *MySQLWriter) CreateTodo(ctx context.Context, t *model.Todo) error {
	q := "INSERT INTO todos(title, created_at) VALUES(?, ?)"
	return m.pool.WithTx(ctx, func(ctx context.Context, tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, q, t.Title, t.CreatedAt.UTC())
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		t.ID = id
		return nil
	})
}
