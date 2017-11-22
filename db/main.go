package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"radio/internal/connect"
	"radio/internal/interrupt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := openDB()
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/search", searchHandler)
	mux.HandleFunc("/stream", streamHandler)

	server := http.Server{
		Addr:           ":" + connect.DatabasePort,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 32,
		Handler:        mux,
	}

	// listen for termination signal!
	interrupt.OnExit(func() { server.Close() })

	fmt.Println("Serving...")
	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		fmt.Println("Done")
	} else {
		log.Println("Server shutdown unexpectedly:", err)
	}
}

func openDB() *sql.DB {
	db, err := sql.Open("mysql", "filedb@/songs")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
