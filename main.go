package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"radio/interrupt"
	"radio/urlgen"
	"radio/websrv"

	_ "github.com/go-sql-driver/mysql"
)

const (
	appDirEnv    = "WEBSRV_APP_DIR"
	urlGenDirEnv = "WEBSRV_URLGEN_DIR"
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	appFileDir := os.Getenv(appDirEnv)
	urlGenDir := os.Getenv(urlGenDirEnv)

	// db startup may take a while, get it going now, while we do other setup
	dbC, errC := connectToDB(60 * time.Second, 6 * time.Second)
	
	if err := websrv.Init(appFileDir); err != nil {
		panic(err)
	}
	defer websrv.Close()

	err := urlgen.LoadDir(filepath.Clean(urlGenDir))
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Addr:           ":8080",
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 8192,
		Handler:        websrv.Routes(),
	}

	waitShutdown := make(chan bool)

	// listen for termination signal
	interrupt.OnExit(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			log.Println("Server shutdown error:", err)
		}
		waitShutdown <- true
	})
	
	// now lets grab (or wait for) the db connection
	var db *sql.DB
	select {
	case db = <-dbC:
		defer db.Close()
		
	case err := <-errC:
		panic(err)
	}

	// hey now we can do what we want
	log.Println("Serving...")
	err = server.ListenAndServe()

	<-waitShutdown
	if err == http.ErrServerClosed {
		log.Println("Done")
	} else {
		log.Println("Server shutdown unexpectedly:", err)
	}
}

func connectToDB(timeout, tryAgainPeriod time.Duration) (chan *sql.DB, chan error) {
	dbC := make(chan *sql.DB, 2)
	errC := make(chan error, 2)
	
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
				db, err := sql.Open("mysql", "database@/radio")
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
