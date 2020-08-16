package api

import (
	"LineProcessor/db_storage"
	pb "LineProcessor/proto"

	"database/sql"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
	"net"
	"strconv"
	"time"
)

type Server struct {
	pb.UnimplementedGRPCApiServer
	dbPtr *sql.DB
}

func (s *Server) SubscribeOnSportsLines(stream pb.GRPCApi_SubscribeOnSportsLinesServer) error {
	reqHandler := make(chan *pb.Request)
	ctx := stream.Context()
	logrus.Infoln("Stream started!")

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
						sportRatio_str, err := strconv.ParseFloat(db_storage.GetSportRatio(s.dbPtr, sport), 32)
						if err != nil {
							logrus.Warningln("Convert error:", err)
						}
						resp := pb.Response{
							SportName:  sport,
							SportRatio: float32(sportRatio_str),
						}
						if err := stream.Send(&resp); err != nil {
							logrus.Errorln("GRPC stream send error:", err)
							return
						}
						logrus.Tracef("Sent response for sport: %s | value: %f", resp.SportName, resp.SportRatio)
					}
				}
			}
			time.Sleep(time.Second * 3)
		}
	} ()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		logrus.Infoln("Got a request for sports:", req.Sport)
		reqHandler <- req
	}
}

func GrpcInit(db *sql.DB, grpcServAddr string) {
	logrus.Infoln("Starting GRPC Server!")
	lis, err := net.Listen("tcp", grpcServAddr + ":50051")
	if err != nil {
		logrus.Fatalln("Failed to listen:", err)
	}
	grpcServer := grpc.NewServer()
	instance := new(Server)
	instance.dbPtr = db
	pb.RegisterGRPCApiServer(grpcServer, instance)
	grpcServer.Serve(lis)


	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalln("Failed to serve GRPC server:", err)
	}
}
