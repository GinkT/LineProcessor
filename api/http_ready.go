package api

import (
	"LineProcessor/db_storage"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

var db_ptr *sql.DB

// Запуск Listener'а, передача указателя на БД, инициализация хэндлеров
func StatusCheckInit(db *sql.DB) {
	port := "8181"
	db_ptr = db

	logrus.Infoln("Trying to initialize Status Check HTTP API on port:", port)
	http.HandleFunc("/ready", CheckConnection)

	go http.ListenAndServe("localhost:" + port, nil)
}

// Описание /ready API
func CheckConnection(w http.ResponseWriter, r *http.Request) {
	if db_storage.IsConnected(db_ptr) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "DB CONNECTED")
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "DB NOT CONNECTED")
	}
}