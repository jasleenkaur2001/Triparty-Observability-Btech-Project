package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"sync"

	"github.com/golang/protobuf/proto"

	pb "Btech_Project/route_guide/area/Proto"
)

var (
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 50053, "The server port")
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
func (s *areaServer) loadArea(filePath string) {
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

func newServer() *areaServer {
	s := &areaServer{}
	s.loadArea(*jsonDBFile)
	return s
}

func main() {
	//flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Area Service listening on port 50053...")
	grpcServer := grpc.NewServer()
	pb.RegisterAreaServer(grpcServer, newServer())
	grpcServer.Serve(lis)

}

var exampleData = []byte(`[{
    "location": {
        "latitude": 407838351,
        "longitude": -746143763
    },
    "area": "Hilly"
},  {
    "location": {
        "latitude": 413843930,
        "longitude": -740501726
    },
    "area": "Snowy"
}, {
    "location": {
        "latitude": 406337092,
        "longitude": -740122226
    },
    "area": "Plateau"
}, {
    "location": {
        "latitude": 406421967,
        "longitude": -747727624
    },
    "area": "Plains"
}, {
    "location": {
        "latitude": 410248224,
        "longitude": -747127767
    },
    "area": "Coastal"
}]`)
