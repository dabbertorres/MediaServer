package websrv

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type dbStmtKey int

const (
	dbStmtUserNew      dbStmtKey = iota
	dbStmtUserLogin
	dbStmtSessionNew
	dbStmtSessionValid
	dbStmtSessionGet
	dbStmtSessionEnd
	dbStmtSearch
	
	dbStmtStrUserNew      = "insert into users (name, password) values ($1, $2)"
	dbStmtStrUserLogin    = "select password from users where name = $1"
	dbStmtStrSessionNew   = "insert into sessions (id, user, ipAddr, userAgent, expires, tunedTo) values ($1, $2, $3, $4, $5, $6)"
	dbStmtStrSessionValid = "select expires from sessions where id = $1"
	dbStmtStrSessionGet   = "select * from sessions where id = $1"
	dbStmtStrSessionEnd   = "delete from sessions where id = $1"
	dbStmtStrSearch       = "select (title, artist) from songs where title like \"%$1%\" or artist like \"%$1%\""
)

var stmts = map[dbStmtKey]*sql.Stmt{}

func PrepareDB(db *sql.DB) error {
	s, err := db.Prepare(dbStmtStrUserNew)
	if err != nil {
		CloseDB(db)
		return err
	}
	stmts[dbStmtUserNew] = s
	
	s, err = db.Prepare(dbStmtStrUserLogin)
	if err != nil {
		CloseDB(db)
		return err
	}
	stmts[dbStmtUserLogin] = s
	
	s, err = db.Prepare(dbStmtStrSessionNew)
	if err != nil {
		CloseDB(db)
		return err
	}
	stmts[dbStmtSessionNew] = s
	
	s, err = db.Prepare(dbStmtStrSessionValid)
	if err != nil {
		CloseDB(db)
		return err
	}
	stmts[dbStmtSessionValid] = s
	
	s, err = db.Prepare(dbStmtStrSessionGet)
	if err != nil {
		CloseDB(db)
		return err
	}
	stmts[dbStmtSessionGet] = s
	
	s, err = db.Prepare(dbStmtStrSessionEnd)
	if err != nil {
		CloseDB(db)
		return err
	}
	stmts[dbStmtSessionEnd] = s
	
	s, err = db.Prepare(dbStmtStrSearch)
	if err != nil {
		CloseDB(db)
		return err
	}
	stmts[dbStmtSearch] = s
	
	return nil
}

func ExecDB(key dbStmtKey, args... interface{}) (sql.Result, error) {
	return stmts[key].Exec(args)
}

func QueryDB(key dbStmtKey, args... interface{}) (*sql.Rows, error) {
	return stmts[key].Query(args)
}

func QueryRowDB(key dbStmtKey, args... interface{}) *sql.Row {
	return stmts[key].QueryRow(args)
}

func CloseDB(db *sql.DB) error {
	for _, s := range stmts {
		s.Close()
	}
	
	return db.Close()
}

func ConnectDB(user, password, address, dbName string, timeout, tryAgainPeriod time.Duration) (<-chan *sql.DB, <-chan error) {
	dbC := make(chan *sql.DB)
	errC := make(chan error)
	
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		
		ticker := time.NewTicker(tryAgainPeriod)
		defer ticker.Stop()
		
		for {
			select {
			case <-ctx.Done():
				errC <- ctx.Err()
				return
			
			case <-ticker.C:
				db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", user, password, address, dbName))
				if err != nil {
					errC <- err
					return
				}
				
				err = db.Ping()
				if err != nil {
					errC <- err
					db.Close()
					return
				}
				
				dbC <- db
			}
		}
	}()
	
	return dbC, errC
}
