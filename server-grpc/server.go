package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"

	"google.golang.org/grpc/reflection"

	pb "github.com/dannyrsu/league-api/leagueservice"
	"github.com/dannyrsu/league-api/models"
	"github.com/golang/protobuf/jsonpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	tls      = flag.Bool("tls", false, "Coonection uses TLS if true")
	certFile = flag.String("cert_file", "", "The TLS cert file")
	keyFile  = flag.String("key_file", "", "The TLS Key file")
	port     = flag.Int("port", 50051, "The server port")
)

type leagueServer struct{}

type summonerStatsResponse struct{}

func constructSummonerStatsResponse(summonerProfile models.SummonerProfile, matchHistory models.MatchHistory) *pb.GetSummonerStatsResponse {

	m := map[string]interface{}{
		"summonerProfile": summonerProfile,
		"matchHistory":    matchHistory,
		"realmData":       models.GetRealmData(),
	}

	jbytes, err := json.Marshal(m)

	if err != nil {
		log.Fatalf("Error marshaling data: %v", err)
	}

	result := &pb.GetSummonerStatsResponse{}

	r := strings.NewReader(string(jbytes))
	if err := jsonpb.Unmarshal(r, result); err != nil {
		panic(err)
	}

	return result
}

func (*leagueServer) GetSummonerStatsUnary(ctx context.Context, req *pb.GetSummonerStatsRequest) (*pb.GetSummonerStatsResponse, error) {
	summonerProfile := models.GetSummonerProfile(req.GetSummonerName(), req.GetRegion())
	matchHistory := models.GetMatchHistory(summonerProfile.AccountID, req.GetRegion(), 0, 5)

	res := constructSummonerStatsResponse(summonerProfile, matchHistory)

	return res, nil
}

func (*leagueServer) GetSummonerStatsBiDirectional(stream pb.LeagueApi_GetSummonerStatsBiDirectionalServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		summonerProfile := models.GetSummonerProfile(req.GetSummonerName(), req.GetRegion())
		matchHistory := models.GetMatchHistory(summonerProfile.AccountID, req.GetRegion(), 0, 5)
		sendErr := stream.Send(constructSummonerStatsResponse(summonerProfile, matchHistory))

		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", err)
		}
	}
}

func main() {
	// if we crash the code, we get the filename and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)

		if err != nil {
			log.Fatalf("Failed to generate credentials: %v", err)
		}

		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterLeagueApiServer(grpcServer, &leagueServer{})
	reflection.Register(grpcServer)

	go func() {
		fmt.Println("Starting League Server ...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until signal is received
	<-ch
	grpcServer.Stop()
	lis.Close()
	fmt.Println("League Server Stopped...")
}
