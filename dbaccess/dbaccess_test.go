package dbaccess

import (
    "database/sql"
    "fmt"
    "testing"
)

func TestDB(t *testing.T) {
    var (
        err   error
        stmt  *sql.Stmt
        db    *DB
        res   sql.Result
        id, a int64
        rows  *sql.Rows
    )

    // db = NewDB("mysql", "root:root@tcp(localhost:3306)/mysql")
    CreateDatabase("mysql", "root", "root", "localhost", "apitest", 3306)

    db = NewDB("mysql", "root", "root", "localhost", "apitest", 3306)
    // close database after all work is done
    defer db.Close()

    db.Ping()
    db.Exec("create table posts(id INT, Name TEXT, Text TEXT)")

    // INSERT INTO DB
    // prepare
    stmt = db.Prepare("insert into posts(id, Name, Text) values (?, ?, ?)")

    // execute
    res, err = stmt.Exec("5", "Post five", "Contents of post 5")
    ErrorCheck(err)

    id, err = res.LastInsertId()
    ErrorCheck(err)

    fmt.Println("Insert id", id)

    // Update db
    stmt = db.Prepare("update posts set Text=? where id=?")

    // execute
    res, err = stmt.Exec("This is post five", "5")
    ErrorCheck(err)

    a, err = res.RowsAffected()
    ErrorCheck(err)

    fmt.Println(a)

    // query all data
    rows = db.Query("select * from posts")
    ErrorCheck(err)

    var post = Post{}

    for rows.Next() {
        err = rows.Scan(&post.Id, &post.Name, &post.Text)
        ErrorCheck(err)
        fmt.Println(post)
    }

    // delete data
    stmt = db.Prepare("delete from posts where id=?")
    ErrorCheck(err)

    // delete 5th post
    res, err = stmt.Exec("5")
    ErrorCheck(err)

    // affected rows
    a, err = res.RowsAffected()
    ErrorCheck(err)

    fmt.Println(a) // 1
}

func ErrorCheck(err error) {
    if err != nil {
        panic(err.Error())
    }
}
