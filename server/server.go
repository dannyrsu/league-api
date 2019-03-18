package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc/reflection"

	pb "github.com/dannyrsu/league-api/leagueservice"
	"github.com/dannyrsu/league-api/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	tls        = flag.Bool("tls", false, "Coonection uses TLS if true")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS Key file")
	port       = flag.Int("port", 50051, "The server port")
	riotAPIKey = flag.String("riot_api_key", "", "Riot API Key")
)

type leagueServer struct{}

func (*leagueServer) GetSummonerProfileUnary(ctx context.Context, req *pb.GetSummonerProfileRequest) (*pb.GetSummonerProfileResponse, error) {
	summonerProfile := models.GetSummonerProfile(req.GetSummonerName(), req.GetRegion())

	res := &pb.GetSummonerProfileResponse{
		ProfileIconId: int32(summonerProfile.ProfileIconID),
		Name:          summonerProfile.Name,
		Puuid:         summonerProfile.PUUID,
		SummonerLevel: int64(summonerProfile.SummonerLevel),
		RevisionDate:  int64(summonerProfile.RevisionDate),
		Id:            summonerProfile.ID,
		AccountId:     summonerProfile.AccountID,
	}

	return res, nil
}

func (*leagueServer) GetSummonerProfileBiDirectional(stream pb.LeagueApi_GetSummonerProfileBiDirectionalServer) error {
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
		sendErr := stream.Send(&pb.GetSummonerProfileResponse{
			ProfileIconId: int32(summonerProfile.ProfileIconID),
			Name:          summonerProfile.Name,
			Puuid:         summonerProfile.PUUID,
			SummonerLevel: int64(summonerProfile.SummonerLevel),
			RevisionDate:  int64(summonerProfile.RevisionDate),
			Id:            summonerProfile.ID,
			AccountId:     summonerProfile.AccountID,
		})

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
