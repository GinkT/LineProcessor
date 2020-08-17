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

type pairNameRatio struct {
	sportName string
	sportRatio float32
}

func (s *Server) SubscribeOnSportsLines(stream pb.GRPCApi_SubscribeOnSportsLinesServer) error {
	reqHandler := make(chan *pb.Request)
	ctx := stream.Context()
	logrus.Infoln("Stream started!")

	go func () {
		request := pb.Request{
			Sport:        nil,
			TimeInterval: 3,
		}
		var initialResponse []pairNameRatio
		isStarted := false
		for {
			select {
			case newreq := <-reqHandler:
				if !isStarted {
					isStarted = true
				}
				request = *newreq
				for _, sport := range request.Sport {
					sportRatio_float, err := strconv.ParseFloat(db_storage.GetSportRatio(s.dbPtr, sport), 32)
					if err != nil {
						logrus.Warningln("Convert error:", err)
					}
					resp := pb.Response{
						SportName:  sport,
						SportRatio: float32(sportRatio_float),
					}
					if err := stream.Send(&resp); err != nil {
						logrus.Errorln("GRPC stream send error:", err)
						return
					}
					initialResponse = append(initialResponse, pairNameRatio{
						sportName:  sport,
						sportRatio: float32(sportRatio_float),
					})
					logrus.Tracef("Sent initial response for sport: %s | value: %f", resp.SportName, resp.SportRatio)
				}

			default :
				if isStarted {
					for i, sport := range request.Sport {
						sportRatio_float, err := strconv.ParseFloat(db_storage.GetSportRatio(s.dbPtr, sport), 32)
						if err != nil {
							logrus.Warningln("Convert error:", err)
						}
						deltaValue := initialResponse[i].sportRatio - float32(sportRatio_float)
						resp := pb.Response{
							SportName:  sport,
							SportRatio: deltaValue,
						}
						if err := stream.Send(&resp); err != nil {
							logrus.Errorln("GRPC stream send error:", err)
							return
						}
						logrus.Tracef("Sent delta response for sport: %s | value: %f | delva value %f", resp.SportName, sportRatio_float, deltaValue)
					}
				}
			}
			time.Sleep(time.Second * time.Duration(request.TimeInterval))
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
	grpcServer := grpc.NewServer()
	instance := new(Server)
	instance.dbPtr = db
	pb.RegisterGRPCApiServer(grpcServer, instance)

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		logrus.Fatalln("Failed to listen:", err)
	}

	logrus.Infoln("Starting to serve GRPC Server!")
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalln("Failed to serve GRPC server:", err)
	}
}
