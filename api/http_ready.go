package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"LineProcessor/db_storage"
	log "github.com/sirupsen/logrus"
)

func Status_Api_Init (db *sql.DB) {
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request){
		if db_storage.DbIsConnected(db) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "DB CONNECTED")
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprint(w, "DB NOT CONNECTED")
		}
	})
	go log.Fatal(http.ListenAndServe("localhost:8181", nil))
}

