package main

import (
	"LineProcessor/api"
	"LineProcessor/db_storage"
	"LineProcessor/http_workers"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

const (
	host     = "db"
	port     = 5432
	user     = "postgres"
	password = "qwerty"
	dbname   = "LinesStorage"
)

//var (
//	grpcAddress = flag.String("gaddr", "", "Адрес GRPC сервера")
//	httpAddress = flag.String("haddr", "linesprovider", "Адрес HTTP сервера")
//	//linesproviderTimeInterval = flag.Int("time", 5, "Интервал синхронизации хранилища с Lines Provider")
//	logLevel = flag.String("log", "trace", "Уровень логирования")
//	baseballTimeInterval = flag.Int("b", 15, "Интервал синхронизации линии спорта BASEBALL")
//	soccerTimeInterval = flag.Int("s", 5, "Интервал синхронизации линии спорта  SOCCER")
//	footballTimeInterval = flag.Int("f", 5, "Интервал синхронизации линии спорта FOOTBALL")
//)

type configuration struct {
	httpServerAddr string
	grpcServerAddr string
	logLevel string
	baseballInterval int
	soccerInterval int
	footballInterval int
}

var Conf configuration

func main() {
	logrus.SetFormatter(formLogger())

	setEnv(&Conf)

	lvl, err := logrus.ParseLevel(Conf.logLevel)
	if err != nil {
		logrus.Fatalln(err)
	}
	logrus.SetLevel(lvl)

	logrus.Infoln("Program started!")

	// Инициализация БД
	db := db_storage.StorageInit(host, port, user, password, dbname)
	defer db.Close()

	// Запуск воркеров по спортам
	go http_workers.RequestWorker("BASEBALL", Conf.baseballInterval, db)
	go http_workers.RequestWorker("SOCCER", Conf.soccerInterval, db)
	go http_workers.RequestWorker("FOOTBALL", Conf.footballInterval, db)

	// Подключение API
	//api.StatusCheckInit(db, Conf.httpServerAddr)
	api.GrpcInit(db, Conf.grpcServerAddr)
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

func setEnv(conf *configuration) {
	var err error
	logrus.Infoln("Setting configuration variables...")
	conf.httpServerAddr = os.Getenv("HTTP_SERVER_ADDR")
	conf.grpcServerAddr = os.Getenv("GRPC_SERVER_ADDR")
	if conf.baseballInterval, err = strconv.Atoi(os.Getenv("BASEBALL_INT")); err != nil {
		logrus.Fatalln("Uncorrect BASEBALL_INT variable!")
	}
	if conf.footballInterval, err = strconv.Atoi(os.Getenv("FOOTBALL_INT")); err != nil {
		logrus.Fatalln("Uncorrect FOOTBALL_INT variable!")
	}
	if conf.soccerInterval, err = strconv.Atoi(os.Getenv("SOCCER_INT")); err != nil {
		logrus.Fatalln("Uncorrect SOCCER_INT variable!")
	}
	conf.logLevel = os.Getenv("LOG_LEVEL")
}

