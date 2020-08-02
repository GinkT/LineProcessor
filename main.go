package main

import (
	"LineProcessor/db_storage"
	pb "LineProcessor/proto"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"
)

var httpserveraddr string = "http://127.0.0.1:8000/api/v1/lines/"

func RequestWorker(sportname string, timeout int, db *sql.DB) {
	var ratiovalue string

	for {
		response, err := http.Get(httpserveraddr + sportname)
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalln(err)
		}
		ratiovalue = strings.TrimFunc(string(body), func(r rune) bool {
			return !unicode.IsDigit(r)
		})

		log.Println(sportname+" worker:", ratiovalue)

		db_storage.PutSportLine(db, sportname, ratiovalue)
		time.Sleep(time.Duration(timeout) * time.Second)
	}
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "qwerty"
	dbname   = "LinesStorage"
)

type server struct {
	Db_Ptr *sql.DB
}

func (s server) SubscribeOnSportsLines(linesServer pb.SubscribeOnSportsLines_SubscribeOnSportsLinesServer) error {
	log.Println("Server started")
	ctx := linesServer.Context()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := linesServer.Recv()
		if err == io.EOF {
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}

		fmt.Printf("Got new request\nSports: %s", req.Sport)
		fmt.Printf("Time interval: %d", req.TimeInterval)

		for {
			for _, sport := range req.Sport {
				fmt.Printf("Getting database select for sport: %s", sport)
				resp := pb.Response{
					SportName:  sport,
					SportRatio: 321.3,
				}
				if err := linesServer.Send(&resp); err != nil {
					log.Printf("send error %v", err)
				}
				log.Println("send new response:")
				log.Println("Sport: ", sport)
			}
			time.Sleep(time.Duration(req.TimeInterval) * time.Second)
		}
	}
}

func main() {
	db := db_storage.Storage_Init(host, port, user, password, dbname)
	defer db.Close()

	go RequestWorker("BASEBALL", 3, db)
	go RequestWorker("SOCCER", 10, db)
	go RequestWorker("FOOTBALL", 10, db)

	//http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request){
	//	if db_storage.DbIsConnected(db) {
	//		w.WriteHeader(http.StatusOK)
	//		fmt.Fprint(w, "DB CONNECTED")
	//	} else {
	//		w.WriteHeader(http.StatusServiceUnavailable)
	//		fmt.Fprint(w, "DB NOT CONNECTED")
	//	}
	//
	//})
	//http.ListenAndServe("localhost:8181", nil)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSubscribeOnSportsLinesServer(s, server{Db_Ptr: db})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	var name string
	fmt.Fscan(os.Stdin, &name)
}
