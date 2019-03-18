package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"google.golang.org/grpc"

	gw "github.com/dannyrsu/league-api/leagueservice"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

var (
	leagueEndpoint = flag.String("league_endpoint", "localhost:50051", "Endpoint for League Service")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterLeagueApiHandlerFromEndpoint(ctx, mux, *leagueEndpoint, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(":8080", mux)
}

func main() {
	flag.Parse()
	defer glog.Flush()
	log.Println("Starting proxy server")
	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
