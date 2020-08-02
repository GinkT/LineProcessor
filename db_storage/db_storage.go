package db_storage

import (
	"database/sql"
	"fmt"
	"log"
)

//type DBStorage struct {
//	*sql.DB
//}

func Storage_Init(host string, port int, user string, password string, dbname string) *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func DbIsConnected(db *sql.DB) bool {
	err := db.Ping()
	if err != nil {
		log.Fatalln(err)
		return false
	}
	log.Println("Database succesfully connected!")
	return true
}

func PutSportLine(db *sql.DB, sportname string, ratiovalue string) {
	sqlStatement := `
			INSERT INTO public.` + `"` + sportname + `"` + `
			VALUES($1, $2)
			ON CONFLICT ("SportName")
			DO
			UPDATE SET "SportRatio" = $2
		`
	_, err := db.Exec(sqlStatement, sportname, ratiovalue)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetSportRatio(db *sql.DB, sportname string) string {
	sqlStatement := `
			SELECT "SportRatio" FROM public.` + `"` + sportname + `"`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var ratio string
		if err := rows.Scan(&ratio); err != nil {
			log.Fatalln(err)
		}
		log.Printf("Got ratio for sport %s: %d", sportname, ratio)
		return ratio
	}
	return ""
}
