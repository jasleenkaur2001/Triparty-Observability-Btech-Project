package main

import (
	observability "Btech_Project/route_guide/Observability"
	pb "Btech_Project/route_guide/location/Proto"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
)

const (
	serverAddress = "localhost:50051"
)

func getFeatureHandler(w http.ResponseWriter, r *http.Request) {

	var point pb.Point
	err := json.NewDecoder(r.Body).Decode(&point)
	if err != nil {
		log.Fatalf("could not connect parse input: %v", err)
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithUnaryInterceptor(observability.UnaryClientInterceptor))
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewLocationClient(conn)
	ctx := context.Background()
	md := metadata.New(map[string]string{"fullmethod": "gateway.gateway/getFeatureHandler", "caller": "client"})
	ctx = metadata.NewIncomingContext(ctx, md)
	feature, err := client.GetFeature(ctx, &point)
	if err != nil {
		log.Fatalf("could not get feature: %v", err)
	}

	response, err := json.Marshal(feature)
	if err != nil {
		log.Fatalf("could not marshal response: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func main() {
	http.HandleFunc("/getFeature", getFeatureHandler)
	fmt.Println("Gateway listening on port 8081...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("failed to start gateway server: %v", err)
	}
}
