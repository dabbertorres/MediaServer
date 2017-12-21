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
	dbPassEnv    = "MYSQL_ROOT_PASSWORD"
)

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	appFileDir := os.Getenv(appDirEnv)
	urlGenDir := os.Getenv(urlGenDirEnv)
	dbPass := os.Getenv(dbPassEnv)

	// db startup may take a while, get it going now, while we do other setup
	dbC, errC := websrv.ConnectDB("root", dbPass, "db", "radio", 19*time.Second, 6*time.Second)

	if err := websrv.Init(appFileDir); err != nil {
		log.Panicln(err)
	}
	defer websrv.Close()

	err := urlgen.LoadDir(filepath.Clean(urlGenDir))
	if err != nil {
		log.Panicln(err)
	}

	var server http.Server

	waitShutdown := make(chan bool)

	// listen for termination signal
	interrupt.OnExit(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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
		defer websrv.CloseDB(db)

	case err := <-errC:
		log.Panicln("Error connecting to db:", err)
	}
	
	router := websrv.Routes()
	router = websrv.MiddlewareDB(db, router)
	
	server = http.Server{
		Addr:           ":8080",
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 8192,
		Handler:        router,
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
