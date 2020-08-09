package db_storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// Подключение к БД
func StorageInit(host string, port int, user string, password string, dbname string) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logrus.Fatalln(err)
	}
	logrus.Infoln("Database was succesfully connected!")
	createTables(db, "SOCCER", "BASEBALL", "FOOTBALL")
	return db
}

func createTables(db *sql.DB, args ...string) {
	for _, sport := range args {
		sqlStatement := `
			CREATE TABLE IF NOT EXISTS ` + sport + ` (
				SportName   text PRIMARY KEY,
				SportRatio  double precision NOT NULL
			);
		`
		_, err := db.Exec(sqlStatement)
		if err != nil {
			logrus.Errorln(err)
			return
		}
		logrus.Infoln("Created table for: ", sport)
	}
	logrus.Infoln("Created tables: ", args)
}

// Проверка подключения к БД
func IsConnected(db *sql.DB) bool {
	err := db.Ping()
	if err != nil {
		logrus.Fatalln(err)
		return false
	}
	logrus.Infoln("Database check successfully gone!")
	return true
}

// Пуш линии спорта в бд
func PutSportLine(db *sql.DB, sportName string, ratioValue string) {
	sqlStatement := `
			INSERT INTO "` + sportName + `"
			VALUES($1, $2)
			ON CONFLICT ("SportName")
			DO
			UPDATE SET "SportRatio" = $2
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
			SELECT "SportRatio" FROM public.` + `"` + sportName + `"`

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