package dbaccess

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "strconv"
)

type Post struct {
    Id   int
    Name string
    Text string
}
type DB struct {
    DriverName string
    db         *sql.DB
    baseUrl    string
    dbName     string
}

func NewDB(DriverName string, user, password, hostname, dbname string, port int) (db *DB) {
    if DriverName != "mysql" {
        panic(fmt.Errorf("Only mysql is supported right now"))
    }
    db = &DB{
        DriverName: DriverName,
        baseUrl:    CreateMySQLBaseDSN(user, password, hostname, port),
        dbName:     dbname,
    }
    db.db = openByDSN(DriverName, db.baseUrl+db.dbName)
    return
}

func openByDSN(driverName, DSN string) (db *sql.DB) {
    var err error
    if db, err = sql.Open(driverName, DSN); err != nil {
        panic(err)
    }
    return
}

func CreateMySQLBaseDSN(user, password, hostname string, port int) (DSN string) {
    DSN = user + ":" + password + "@tcp(" + hostname + ":" + strconv.Itoa(port) + ")/"
    return
}

func CreateDatabase(DriverName string, user, password, hostname, dbname string, port int) {
    var err error
    dbAdmin := openByDSN(DriverName, CreateMySQLBaseDSN(user, password, hostname, port)+"mysql")
    if _, err = dbAdmin.Exec("drop database if exists " + dbname); err != nil {
        panic(err)
    }
    if _, err = dbAdmin.Exec("create database " + dbname); err != nil {
        panic(err)
    }
    _ = dbAdmin.Close()
}

func (db *DB) Exec(stmt string) {
    var err error
    if _, err = db.db.Exec(stmt); err != nil {
        panic(err)
    }
}

func (db *DB) Query(stmt string) (rows *sql.Rows) {
    var err error
    if rows, err = db.db.Query(stmt); err != nil {
        panic(err)
    }
    return
}

func (db *DB) Close() {
    _ = db.db.Close()
}

func (db *DB) Ping() {
    _ = db.db.Ping()
}

func (db *DB) Prepare(stmt string) (s *sql.Stmt) {
    var err error
    if s, err = db.db.Prepare(stmt); err != nil {
        panic(err)
    }
    return
}
