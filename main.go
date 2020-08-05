package main

import (
	"LineProcessor/api"
	"LineProcessor/db_storage"
	"LineProcessor/http_workers"
	pb "LineProcessor/proto"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"io"
	//"log"
	"net"
	"os"
	"time"
	log "github.com/sirupsen/logrus"
)

var httpserveraddr string = "http://127.0.0.1:8000/api/v1/lines/"

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "qwerty"
	dbname   = "LinesStorage"
)

type server struct {
	pb.UnimplementedGRPCApiServer
	Db_Ptr *sql.DB
}

func (s *server) SubscribeOnSportsLines(stream pb.GRPCApi_SubscribeOnSportsLinesServer) error {
		reqHandler := make(chan *pb.Request)

		go func () {
			request := pb.Request{
				Sport:        nil,
				TimeInterval: 0,
			}
			isStarted := false
			for {
				select {
				case newreq := <- reqHandler:
					if !isStarted {
						isStarted = true
					}
					request = *newreq
				default :
					if isStarted {
						for _, sport := range request.Sport {
							log.Println("DB select for: ", sport)
							resp := pb.Response{
								SportName:  sport,
								SportRatio: float32(db_storage.GetSportRatio(s.Db_Ptr, sport)),
							}
							if err := stream.Send(&resp); err != nil {
								log.Printf("send error %v", err)
							}
							log.Printf("Response for %s was sent!", sport)
						}
					}
				}
				time.Sleep(time.Second * 3)
			}
		} ()

		for {
			req, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			log.Println("Got a request for sports:", req.Sport)
			reqHandler <- req
		}
}

func main() {
	db := db_storage.Storage_Init(host, port, user, password, dbname)
	defer db.Close()

	go http_workers.RequestWorker("BASEBALL", 3, db, httpserveraddr)
	go http_workers.RequestWorker("SOCCER", 10, db, httpserveraddr)
	go http_workers.RequestWorker("FOOTBALL", 10, db, httpserveraddr)

	api.Status_Api_Init(db)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	instance := new(server)
	instance.Db_Ptr = db
	pb.RegisterGRPCApiServer(grpcServer, instance)
	grpcServer.Serve(lis)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	var name string
	fmt.Fscan(os.Stdin, &name)
}
