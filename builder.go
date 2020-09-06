package simpsql

import "strings"

type builder struct {
	tableName string
	simpsql   *Simpsql
}

func (b *builder) Insert(columns map[string]interface{}) *Execute {
	//insert command sql
	q := strings.Builder{}
	//Pre-defined allocation
	q.Grow(256)
	q.WriteString("INSERT INTO `")
	q.WriteString(b.tableName)
	q.WriteString("` (")
	//columns length
	cl := len(columns)
	var (
		values = make([]interface{}, 0, cl)
		//placeholder
		ph = strings.Builder{}
	)
	ph.Grow(32)
	for column, value := range columns {
		values = append(values, value)
		q.WriteByte('`')
		q.WriteString(column)
		q.WriteByte('`')
		ph.WriteByte('?')
		if len(values) < cl {
			q.WriteByte(',')
			ph.WriteByte(',')
		}
	}
	q.WriteString(") VALUES (")
	q.WriteString(ph.String())
	q.WriteByte(')')
	return b.simpsql.PrepareExecute(q.String(), values...)
}

func (b *builder) BatchInsert(columns []string, data [][]interface{}) *Execute {
	const defaultLen = 23
	grow, cLen, dLen := 0, len(columns), len(data)
	grow += defaultLen + len(b.tableName) + 2*dLen*(cLen+1)
	for _, column := range columns {
		grow += len(column) + 3
	}
	q := strings.Builder{}
	q.Grow(grow)
	q.WriteString("INSERT INTO `")
	q.WriteString(b.tableName)
	q.WriteString("`(")
	for i, column := range columns {
		q.WriteByte('`')
		q.WriteString(column)
		q.WriteByte('`')
		if i != cLen-1 {
			q.WriteByte(',')
		}
	}
	q.WriteString(") VALUES ")
	for ri := 0; ri < dLen; ri++ {
		q.WriteByte('(')
		for ci := 0; ci < cLen; ci++ {
			q.WriteByte('?')
			if ci != cLen-1 {
				q.WriteByte(',')
			}
		}
		q.WriteByte(')')
		if ri != dLen-1 {
			q.WriteByte(',')
		}
	}
	values := make([]interface{}, 0, dLen*cLen)
	for _, datum := range data {
		values = append(values, datum...)
	}
	return b.simpsql.PrepareExecute(q.String(), values...)
}

func (b *builder) Delete(condition string, args ...interface{}) *Execute {
	q := strings.Builder{}
	q.Grow(len(b.tableName) + len(condition) + 19)
	q.WriteString("DELETE FROM `")
	q.WriteString(b.tableName)
	q.WriteByte('`')
	if condition != "" {
		q.WriteString(" WHERE ")
		q.WriteString(condition)
	}
	return b.simpsql.PrepareExecute(q.String(), args...)
}

func (b *builder) Update(data map[string]interface{}, condition string, args ...interface{}) *Execute {
	q := strings.Builder{}
	q.Grow(256)
	q.WriteString("UPDATE `")
	q.WriteString(b.tableName)
	q.WriteString("` SET ")
	values := make([]interface{}, 0, len(data)+len(args))
	sets := strings.Builder{}
	sets.Grow(128)
	for key, value := range data {
		sets.WriteByte('`')
		sets.WriteString(key)
		sets.WriteByte('`')
		sets.WriteString(" = ?,")
		values = append(values, value)
	}
	q.WriteString(sets.String()[:sets.Len()-1])
	if condition != "" {
		q.WriteString(" WHERE ")
		q.WriteString(condition)
		values = append(values, args...)
	}
	return b.simpsql.PrepareExecute(q.String(), values...)
}
