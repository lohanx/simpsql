package simpsql

import (
	"database/sql"
)

type Execute struct {
	result sql.Result
	err    error
}

func (e *Execute) LastInsertID() (int64, error) {
	if e.err != nil {
		return 0, e.err
	}
	return e.result.LastInsertId()
}

func (e *Execute) RowsAffected() (int64, error) {
	if e.err != nil {
		return 0, e.err
	}
	return e.result.RowsAffected()
}
