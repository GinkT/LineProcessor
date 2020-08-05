package http_workers

import (
	"database/sql"
	"github.com/GinkT/LineProcessor/db_storage"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"
)

func RequestWorker(sportName string, timeoutInterval int, db *sql.DB, httpAddress string) {
	var sportRatio string

	for {
		// Получаем страницу
		response, err := http.Get(httpAddress + sportName)
		if err != nil {
			log.Fatalln(err)
		}
		// Читаем тело
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}
		// Обрезаем ненужное
		sportRatio = strings.TrimFunc(string(body), func(r rune) bool {
			return !unicode.IsDigit(r)
		})

		db_storage.PutSportLine(db, sportName, sportRatio)	// Сохраняем в БД
		time.Sleep(time.Duration(timeoutInterval) * time.Second)
	}
}