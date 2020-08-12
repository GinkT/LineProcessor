package http_workers

import (
	"LineProcessor/db_storage"
	"database/sql"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"unicode"
)

// Воркер для чтения HTTP сервера
func RequestWorker(sportName string, timeoutInterval int, db *sql.DB) {
	var sportRatio string
	httpAddress := "http://linesprovider:8000/api/v1/lines/"

	logrus.Infoln(sportName, " worker initialized!")
	for {
		time.Sleep(time.Duration(timeoutInterval) * time.Second)

		// Получаем страницу
		response, err := http.Get(httpAddress + sportName)
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		// Читаем тело
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			logrus.Errorln(err)
			continue
		}
		// Обрезаем ненужное
		sportRatio = strings.TrimFunc(string(body), func(r rune) bool {
			return !unicode.IsDigit(r)
		})
		logrus.Tracef("Got a body from page %s: %s\tTrimmed it to: %s", httpAddress + sportName, string(body), sportRatio)

		db_storage.PutSportLine(db, sportName, sportRatio)	// Сохраняем в БД
	}
}