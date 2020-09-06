package simpsql

import (
	"context"
	"database/sql"
)

type Transaction struct {
	err  error
	conn *sql.Tx
}

func (t *Transaction) Execute(query string, args ...interface{}) *Execute {
	e := new(Execute)
	if t.err != nil {
		e.err = t.err
		return e
	}
	e.result, e.err = t.conn.ExecContext(context.Background(), query, args...)
	return e
}

func (t *Transaction) PrepareExecute(query string, args ...interface{}) *Execute {
	e := new(Execute)
	if t.err != nil {
		e.err = t.err
		return e
	}
	stmt, err := t.conn.PrepareContext(context.Background(), query)
	if err != nil {
		e.err = err
		return e
	}
	defer stmt.Close()
	e.result, e.err = stmt.ExecContext(context.Background(), args...)
	return e
}

func (t *Transaction) Query(query string, args ...interface{}) *Query {
	q := new(Query)
	q.rows, q.err = t.conn.QueryContext(context.Background(), query, args...)
	return q
}

func (t *Transaction) PrepareQuery(query string, args ...interface{}) *Query {
	q := new(Query)
	stmt, err := t.conn.PrepareContext(context.Background(), query)
	if err != nil {
		q.err = err
		return q
	}
	defer stmt.Close()
	q.rows, q.err = stmt.QueryContext(context.Background(), args...)
	return q
}

func (t *Transaction) Commit() error {
	return t.conn.Commit()
}

func (t *Transaction) RollBack() error {
	return t.conn.Rollback()
}
