package main

import (
	observability "Btech_Project/route_guide/Observability"
	"Btech_Project/route_guide/database"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"

	"github.com/golang/protobuf/proto"

	pb "Btech_Project/route_guide/area/Proto"
)

var (
	port = flag.Int("port", 50053, "The server port")
)

type areaServer struct {
	pb.UnimplementedAreaServer
	savedFeatures []*pb.DbFeature // read-only after initialized

	mu sync.Mutex // protects routeNotes
}

func (s *areaServer) GetArea(ctx context.Context, point *pb.Point) (*pb.Area, error) {
	for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return &pb.Area{Area: feature.Area}, nil
		}
	}
	// No feature was found, return an unnamed feature
	return &pb.Area{Area: "UnFound"}, nil
}

// loadFeatures loads features from a JSON file.
func (s *areaServer) loadArea() {
	var data []byte
	data = database.GetAreaData()
	if err := json.Unmarshal(data, &s.savedFeatures); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
}

func newServer() *areaServer {
	s := &areaServer{}
	s.loadArea()
	return s
}

func main() {
	//flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Area Service listening on port 50053...")
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(observability.UnaryServerInterceptor))
	pb.RegisterAreaServer(grpcServer, newServer())
	grpcServer.Serve(lis)

}
