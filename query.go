package simpsql

import (
	"database/sql"
)

type Query struct {
	rows *sql.Rows
	err  error
}

func (q *Query) Rows() (*sql.Rows, error) {
	return q.rows, q.err
}

func (q *Query) FetchOne() (map[string]interface{}, error) {
	if q.err != nil {
		return nil, q.err
	}
	defer q.rows.Close()
	columns, dest, err := q.parseColumnTypes()
	if err != nil {
		return nil, err
	}
	if !q.rows.Next() {
		if err := q.rows.Err(); err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}
	if err := q.rows.Scan(dest...); err != nil {
		return nil, err
	}
	datum := make(map[string]interface{})
	for i, d := range dest {
		datum[columns[i]] = ptrval(d)
	}
	return datum, q.rows.Close()
}

func (q *Query) FetchAll() ([]map[string]interface{}, error) {
	if q.err != nil {
		return nil, q.err
	}
	defer q.rows.Close()
	columns, dest, err := q.parseColumnTypes()
	if err != nil {
		return nil, err
	}
	var data []map[string]interface{}
	for q.rows.Next() {
		if err := q.rows.Scan(dest...); err != nil {
			return nil, err
		}
		datum := make(map[string]interface{})
		for i, d := range dest {
			datum[columns[i]] = ptrval(d)
		}
		data = append(data, datum)
	}
	if err := q.rows.Err(); err != nil {
		return nil, err
	}
	return data, q.rows.Close()
}

func (q *Query) parseColumnTypes() ([]string, []interface{}, error) {
	types, err := q.rows.ColumnTypes()
	if err != nil {
		return nil, nil, err
	}
	tl := len(types)
	dest := make([]interface{}, tl)
	columns := make([]string, tl)
	for i, t := range types {
		columns[i] = t.Name()
		switch t.DatabaseTypeName() {
		case "BOOLEAN", "BOOL":
			var v bool
			dest[i] = &v
		case "BIT":
			var v uint8
			dest[i] = &v
		case "TINYINT", "SMALLINT", "MEDIUMINT", "INT":
			var v int
			dest[i] = &v
		case "FLOAT", "REAL":
			var v float32
			dest[i] = &v
		case "DOUBLE", "DECIMAL":
			var v float64
			dest[i] = &v
		case "CHAR", "VARCHAR", "NVARCHAR", "TEXT", "TINYTEXT",
			"MEDIUMTEXT", "LONGTEXT", "VARBINARY", "ENUM",
			"JSON", "SET", "DATE", "DATETIME", "YEAR", "TIME":
			var v string
			dest[i] = &v
		case "BIGINT", "TIMESTAMP":
			var v int64
			dest[i] = &v
		case "BINARY", "TINYBLOB", "BLOB", "MEDIUMBLOB", "LONGBLOB":
			fallthrough
		default:
			var v []byte
			dest[i] = &v
		}
	}
	return columns, dest, nil
}

func ptrval(val interface{}) interface{} {
	switch val.(type) {
	case *uint8:
		return *val.(*uint8)
	case *int:
		return *val.(*int)
	case *string:
		return *val.(*string)
	case *float32:
		return *val.(*float32)
	case *float64:
		return *val.(*float64)
	case *bool:
		return *val.(*bool)
	case *[]byte:
		return *val.(*[]byte)
	default:
		return *val.(*[]byte)
	}
}
