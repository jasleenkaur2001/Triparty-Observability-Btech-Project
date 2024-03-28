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
	"time"

	"github.com/golang/protobuf/proto"

	pb "Btech_Project/route_guide/name/Proto"
)

var (
	port = flag.Int("port", 50052, "The server port")
)

type nameServer struct {
	pb.UnimplementedNameServer
	savedFeatures []*pb.DbFeature // read-only after initialized

	mu sync.Mutex // protects routeNotes
}

func (s *nameServer) GetName(ctx context.Context, point *pb.Point) (*pb.Name, error) {
	time.Sleep(15 * time.Second)
	for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return &pb.Name{Name: feature.Name}, nil
		}
	}
	// No feature was found, return an unnamed feature
	return &pb.Name{Name: "unnamed"}, nil
}

// loadFeatures loads features from a JSON file.
func (s *nameServer) loadFeatures() {
	var data []byte
	data = database.GetNameData()
	if err := json.Unmarshal(data, &s.savedFeatures); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
}

func newServer() *nameServer {
	s := &nameServer{}
	s.loadFeatures()
	return s
}

func main() {
	//flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Name service listening on port 50052...")
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(observability.UnaryServerInterceptor),
	)
	pb.RegisterNameServer(grpcServer, newServer())
	grpcServer.Serve(lis)

}
