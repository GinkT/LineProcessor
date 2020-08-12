package db_storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"strings"
)

// Подключение к БД
func StorageInit(host string, port int, user string, password string, dbname string) *sql.DB {
	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user, password, host, port, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logrus.Fatalln(err)
	}
	logrus.Infoln("Database was successfully connected!")
	return db
}

// Проверка подключения к БД
func IsConnected(db *sql.DB) bool {
	err := db.Ping()
	if err != nil {
		logrus.Errorln(err)
		return false
	}
	logrus.Infoln("Database check successfully gone!")
	return true
}

// Пуш линии спорта в бд
func PutSportLine(db *sql.DB, sportName string, ratioValue string) {
	sqlStatement := `
			INSERT INTO public."` + strings.ToLower(sportName) + `"
			VALUES($1, $2)
			ON CONFLICT ("sportname")
			DO
			UPDATE SET "sportratio" = $2
		`
	_, err := db.Exec(sqlStatement, sportName, ratioValue)
	if err != nil {
		logrus.Errorln(err)
	}
	logrus.Tracef("Pushed sport %s and value %s to DB", sportName, ratioValue)
}

// Пулл линии спорта из БД
func GetSportRatio(db *sql.DB, sportName string) string {
	sqlStatement := `
			SELECT "sportratio" FROM public.` + `"` + strings.ToLower(sportName) + `"`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		logrus.Errorln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var ratio string
		if err := rows.Scan(&ratio); err != nil {
			logrus.Errorln(err)
		}
		logrus.Tracef("Got ratio from DB for sport %s: %s", sportName, ratio)

		return ratio
	}
	return ""
}