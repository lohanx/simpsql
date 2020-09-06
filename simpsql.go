package simpsql

import (
	"context"
	"database/sql"
)

type Simpsql struct {
	db *sql.DB
}

func New(db *sql.DB) *Simpsql {
	return &Simpsql{db: db}
}

func (s *Simpsql) Query(query string, args ...interface{}) *Query {
	q := new(Query)
	q.rows, q.err = s.db.QueryContext(context.Background(), query, args...)
	return q
}

func (s *Simpsql) PrepareQuery(query string, args ...interface{}) *Query {
	q := new(Query)
	stmt, err := s.db.PrepareContext(context.Background(), query)
	if err != nil {
		q.err = err
		return q
	}
	defer stmt.Close()
	q.rows, q.err = stmt.QueryContext(context.Background(), args...)
	return q
}

func (s *Simpsql) Execute(query string, args ...interface{}) *Execute {
	e := new(Execute)
	e.result, e.err = s.db.ExecContext(context.Background(), query, args...)
	return e
}

func (s *Simpsql) BeginTransaction() *Transaction {
	t := new(Transaction)
	t.conn, t.err = s.db.Begin()
	return t
}

func (s *Simpsql) PrepareExecute(query string, args ...interface{}) *Execute {
	e := new(Execute)
	stmt, err := s.db.PrepareContext(context.Background(), query)
	if err != nil {
		e.err = err
		return e
	}
	defer stmt.Close()
	e.result, e.err = stmt.ExecContext(context.Background(), args...)
	return e
}

func (s *Simpsql) Table(tableName string) *builder {
	return &builder{tableName: tableName, simpsql: s}
}

func (s *Simpsql) DB() *sql.DB {
	return s.db
}
