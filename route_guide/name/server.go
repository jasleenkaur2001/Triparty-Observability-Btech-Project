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

	pb "Btech_Project/route_guide/name/Proto"
)

var (
	jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port       = flag.Int("port", 50052, "The server port")
)

type nameServer struct {
	pb.UnimplementedNameServer
	savedFeatures []*pb.DbFeature // read-only after initialized

	mu sync.Mutex // protects routeNotes
}

func (s *nameServer) GetName(ctx context.Context, point *pb.Point) (*pb.Name, error) {
	for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, point) {
			return &pb.Name{Name: feature.Name}, nil
		}
	}
	// No feature was found, return an unnamed feature
	return &pb.Name{Name: "unnamed"}, nil
}

// loadFeatures loads features from a JSON file.
func (s *nameServer) loadFeatures(filePath string) {
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

func newServer() *nameServer {
	s := &nameServer{}
	s.loadFeatures(*jsonDBFile)
	return s
}

func main() {
	//flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("Name service listening on port 50052...")
	grpcServer := grpc.NewServer()
	pb.RegisterNameServer(grpcServer, newServer())
	grpcServer.Serve(lis)

}

var exampleData = []byte(`[{
    "location": {
        "latitude": 407838351,
        "longitude": -746143763
    },
    "name": "Patriots Path, Mendham, NJ 07945, USA"
},  {
    "location": {
        "latitude": 413843930,
        "longitude": -740501726
    },
    "name": "162 Merrill Road, Highland Mills, NY 10930, USA"
}, {
    "location": {
        "latitude": 406337092,
        "longitude": -740122226
    },
    "name": "6324 8th Avenue, Brooklyn, NY 11220, USA"
}, {
    "location": {
        "latitude": 406421967,
        "longitude": -747727624
    },
    "name": "1 Merck Access Road, Whitehouse Station, NJ 08889, USA"
}, {
    "location": {
        "latitude": 410248224,
        "longitude": -747127767
    },
    "name": "3 Hasta Way, Newton, NJ 07860, USA"
}]`)
