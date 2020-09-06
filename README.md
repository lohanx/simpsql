## Installation
```
go get -u github.com/lohanx/simpsql
```
## Example
```go
import (
    "database/sql"
    "github.com/go-sql-driver/mysql"
    "github.com/lohanx/simpsql"
)

func main() {
    cfg := mysql.Config{
        User:   "user",
        Passwd: "passwd",
        Net:    "tcp",
        Addr:   "127.0.0.1:3306",
        DBName: "db_name",
    }
    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        panic(err)
    }
    
    sdb := simpsql.New(db)
    //prepare query
    s1 := "SELECT `id`,`username`,`email` FROM `users` WHERE `state` = ?"
    users, err := sdb.PrepareQuery(s1, 1).FetchAll()
    
    //query
    users, err := sdb.Query(s1, 1).FetchAll()
    
    //create
    s2 := "INSERT INTO `users` (`username`,`email`,`state`) VALUES (?,?,?)"
    lastId,err := sdb.Execute(s2, "lohanx", "example@lohanx.cn",1).LastInsertID()
    //or
    lastId,err := sdb.Table("users").Insert(map[string]interface{}{
        "username":"lohanx",
        "email":"example@lohanx.cn",
        "state":1,
    }).LastInsertID()    

    //update
    sdb.Table("users").Update(map[string]interface{}{
        "username":"test",
        "email":"test@lohanx.cn"
    },"id = ?",1).RowsAffected()

    //delete
    sdb.Table("users").Delete("id = ?",1).RowsAffected()
    //prepare create
    lastId,err := sdb.PrepareExecute(s2, "lohanx", "example@lohanx.cn",1).LastInsertID()
    
    //transaction
    tx := sdb.BeginTransaction()
    id,err := tx.PrepareExecute(s2).LastInsertID()
    ...
    if err != nil {
        tx.RollBack()
        return
    }
    tx.Commit()
}
```
