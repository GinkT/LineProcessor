package main

import (
	"LineProcessor/api"
	"LineProcessor/db_storage"
	"LineProcessor/http_workers"
	"flag"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "qwerty"
	dbname   = "LinesStorage"
)

var (
	grpcAddress = flag.String("gaddr", "", "Адрес GRPC сервера")
	httpAddress = flag.String("haddr", "127.0.0.1", "Адрес HTTP сервера")
	//linesproviderTimeInterval = flag.Int("time", 5, "Интервал синхронизации хранилища с Lines Provider")
	logLevel = flag.String("log", "trace", "Уровень логирования")
	baseballTimeInterval = flag.Int("b", 5, "Интервал синхронизации линии спорта BASEBALL")
	soccerTimeInterval = flag.Int("s", 5, "Интервал синхронизации линии спорта  SOCCER")
	footballTimeInterval = flag.Int("f", 5, "Интервал синхронизации линии спорта FOOTBALL")
)

func main() {
	flag.Parse()

	logrus.SetFormatter(formLogger())
	lvl, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logrus.Fatalln(err)
	}
	logrus.SetLevel(lvl)

	logrus.Infoln("Program started!")

	logrus.Infoln("This is a new build 12:10")

	// Инициализация БД
	db := db_storage.StorageInit(host, port, user, password, dbname)
	defer db.Close()

	// Запуск воркеров по спортам
	go http_workers.RequestWorker("BASEBALL", *baseballTimeInterval, db, *httpAddress)
	//go http_workers.RequestWorker("SOCCER", *soccerTimeInterval, db, *httpAddress)
	// go http_workers.RequestWorker("FOOTBALL", *footballTimeInterval, db, *httpAddress)

	// Подключение API
	api.StatusCheckInit(db)
	api.GrpcInit(db, *grpcAddress)
}

func formLogger () *logrus.TextFormatter {
	formatter := new(logrus.TextFormatter)
	formatter = &logrus.TextFormatter{
		ForceColors:               true,
		DisableColors:             false,
		ForceQuote:                true,
		DisableQuote:              false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02 15:04:05",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		PadLevelText:              false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier:          nil,
	}
	return formatter
}

