package webserver

import (
    "database/sql"
    "fmt"
    "github.com/notargets/api-tester/dbaccess"
    "io"
    "net/http"
    "strconv"
)

var (
    accessCount int
    db          *dbaccess.DB
)

func countAccess(w http.ResponseWriter, r *http.Request) {
    var (
        rows *sql.Rows
        err  error
    )
    accessCount++
    // stmt.Exec(accessCount, "count")
    db.Exec("UPDATE webtable SET a = " + strconv.Itoa(accessCount))
    rows = db.Query("SELECT a FROM webtable")
    var result int
    rows.Next()
    if err = rows.Scan(&result); err != nil {
        panic(err)
    }
    msg := fmt.Sprintf("Internal count: %d From DB: %d", accessCount, result)
    _, _ = io.WriteString(w, msg)
}

func Start(port int) {
    var err error
    dbaccess.CreateDatabase("mysql", "root", "root", "localhost", "apitest", 3306)
    db = dbaccess.NewDB("mysql", "root", "root", "localhost", "apitest", 3306)
    db.Exec(`CREATE TABLE webtable (a int, b text)`)
    db.Exec(`INSERT INTO webtable VALUES (0, "initial")`)
    mux := http.NewServeMux()

    mux.HandleFunc("/countAccess", countAccess)

    aport := strconv.Itoa(port)
    if err = http.ListenAndServe("localhost:"+aport, mux); err != nil {
        panic(err)
    }
}
