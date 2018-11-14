package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("sqlite3", "/tmp/crates-mirror/crates.db")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println(r.URL.Path)
		if !strings.HasSuffix(r.URL.Path, "download") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		reqSlice := strings.Split(r.URL.Path, "/")
		l := len(reqSlice)
		if l > 3 {
			version := reqSlice[l-2]
			name := reqSlice[l-3]
			var cratePath string
			err := db.QueryRow("select path from crate_version where name=? and version=?", name, version).Scan(&cratePath)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			f, err := os.Open(cratePath)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err = io.Copy(w, f)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
