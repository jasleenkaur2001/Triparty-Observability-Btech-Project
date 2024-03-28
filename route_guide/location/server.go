package main

import (
	observability "Btech_Project/route_guide/Observability"
	pbArea "Btech_Project/route_guide/area/Proto"
	"Btech_Project/route_guide/database"
	pb "Btech_Project/route_guide/location/Proto"
	pbName "Btech_Project/route_guide/name/Proto"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"sync"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type LocationServer struct {
	pb.UnimplementedLocationServer
	savedFeatures []*pb.Point

	mu sync.Mutex // protects routeNotes
}

func (s *LocationServer) ValidateLocation(point *pb.Point) bool {
	for _, val := range s.savedFeatures {
		if proto.Equal(val, point) {
			return true
		}
	}
	return false
}

func (s *LocationServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	if s.ValidateLocation(point) {
		// call 2 rpc's

		name := s.GetName(ctx, point)
		area := s.GetArea(ctx, point)
		return &pb.Feature{Name: name, Location: point, Area: area}, nil

	}
	return &pb.Feature{Location: point}, nil
}

// loadFeatures loads features from a JSON file.
func (s *LocationServer) loadPoint() {
	var data []byte
	data = database.GetLocData()
	if err := json.Unmarshal(data, &s.savedFeatures); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
}
func (s *LocationServer) GetName(ctx context.Context, point *pb.Point) string {
	serverAddress := "localhost:50052"
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithUnaryInterceptor(observability.UnaryClientInterceptor))
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()
	client := pbName.NewNameClient(conn)
	namePoint := &pbName.Point{Latitude: point.Latitude, Longitude: point.Longitude}

	name, err := client.GetName(ctx, namePoint)
	if err != nil {
		log.Fatalf("could not get feature: %v", err)
	}
	return name.Name
}

func (s *LocationServer) GetArea(ctx context.Context, point *pb.Point) string {
	serverAddress := "localhost:50053"
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithUnaryInterceptor(observability.UnaryClientInterceptor))
	conn, err := grpc.Dial(serverAddress, opts...)
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()
	client := pbArea.NewAreaClient(conn)
	areaPoint := &pbArea.Point{Latitude: point.Latitude, Longitude: point.Longitude}

	area, err := client.GetArea(ctx, areaPoint)
	if err != nil {
		log.Fatalf("could not get feature: %v", err)
	}
	return area.Area
}

func newServer() *LocationServer {
	s := &LocationServer{}
	s.loadPoint()
	return s
}

func main() {
	//flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Location service listening on port 50051...")
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(observability.UnaryServerInterceptor),
	)
	pb.RegisterLocationServer(grpcServer, newServer())
	grpcServer.Serve(lis)

}
