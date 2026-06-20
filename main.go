package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	port := flag.Int("port", 7293, "server port")
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dbPath := filepath.Join(dir, "waxmuseum.db")

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = initDB(db); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, filepath.Join(dir, "static", "index.html"))
			return
		}
		http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
	})

	mux.HandleFunc("/api/halls", hallsHandler)
	mux.HandleFunc("/api/guides", guidesHandler)
	mux.HandleFunc("/api/reservations", reservationsHandler)
	mux.HandleFunc("/api/assign", assignHandler)
	mux.HandleFunc("/api/stats/monthly-theme", monthlyThemeStatsHandler)
	mux.HandleFunc("/api/reservations/today", todayReservationsHandler)

	addr := ":" + strconv.Itoa(*port)
	fmt.Printf("Wax Museum server starting on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
