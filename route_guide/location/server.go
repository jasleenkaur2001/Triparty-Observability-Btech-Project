package main

import (
	pbArea "Btech_Project/route_guide/area/Proto"
	pb "Btech_Project/route_guide/location/Proto"
	pbName "Btech_Project/route_guide/name/Proto"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var (
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 50051, "The server port")
)

type LocationServer struct {
	pb.UnimplementedLocationServer
	savedFeatures []*pb.Point

	mu sync.Mutex // protects routeNotes
}

func logFunc(caller string, strt time.Time, callee string) {
	log.Println(caller + "->" + pb.Location_ServiceDesc.ServiceName + "->" + callee + "Call completed in " + time.Now().UTC().Sub(strt).String())
}
func (s *LocationServer) ValidateLocation(ctx context.Context, point *pb.Point) bool {
	for _, val := range s.savedFeatures {
		if proto.Equal(val, point) {
			return true
		}
	}
	return false
}

func (s *LocationServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	if s.ValidateLocation(ctx, point) {
		// call 2 rpc's

		name := s.GetName(ctx, point)
		area := s.GetArea(ctx, point)
		return &pb.Feature{Name: name, Location: point, Area: area}, nil

	}
	return &pb.Feature{Location: point}, nil
}

// loadFeatures loads features from a JSON file.
func (s *LocationServer) loadPoint(filePath string) {
	var data []byte
	if filePath != "" {
		var err error
		data, err = os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Failed to load default features: %v", err)
		}
	} else {
		data = exampleData
	}
	if err := json.Unmarshal(data, &s.savedFeatures); err != nil {
		log.Fatalf("Failed to load default features: %v", err)
	}
}
func (s *LocationServer) GetName(ctx context.Context, point *pb.Point) string {

	caller := "Unknown"
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		caller = md["caller"][0]
	}
	defer logFunc(caller, time.Now().UTC(), pbName.Name_ServiceDesc.ServiceName)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("caller", pb.Location_ServiceDesc.ServiceName))
	serverAddress := "localhost:50052"
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	caller := "Unknown"
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		caller = md["caller"][0]
	}
	defer logFunc(caller, time.Now().UTC(), pbArea.Area_ServiceDesc.ServiceName)
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("caller", pb.Location_ServiceDesc.ServiceName))
	serverAddress := "localhost:50053"
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	s.loadPoint(*jsonDBFile)
	return s
}

func main() {
	//flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Location service listening on port 50051...")
	grpcServer := grpc.NewServer()
	pb.RegisterLocationServer(grpcServer, newServer())
	grpcServer.Serve(lis)

}

var exampleData = []byte(`[
{
    "latitude": 407838351,
    "longitude": -746143763
  },
  {
    "latitude": 413843930,
    "longitude": -740501726
  },
  {
    "latitude": 406337092,
    "longitude": -740122226
  },
  {
    "latitude": 406421967,
    "longitude": -747727624
  },
  {
    "latitude": 410248224,
    "longitude": -747127767
  }
]`)
