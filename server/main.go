package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Aksh-Bansal-dev/grpc-vs-json/greet"
	pb "github.com/Aksh-Bansal-dev/grpc-vs-json/greet"
	"google.golang.org/grpc"
)

var (
	port           = flag.Int("port", 50051, "the server port")
	encodingFormat = flag.String("enc", "grpc", "encoding method (grpc/json)")
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) Auth(ctx context.Context, in *pb.User) (*pb.Response, error) {
	return &pb.Response{Message: "Hello " + in.GetName()}, nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if *encodingFormat == "grpc" {
		s := grpc.NewServer()
		pb.RegisterGreeterServer(s, &server{})
		log.Printf("server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	} else {
		defer lis.Close()
		for {
			conn, err := lis.Accept()
			if err != nil {
				log.Fatal(err)
			}
			// close conn
			defer conn.Close()
			decoder := json.NewDecoder(conn)
			encoder := json.NewEncoder(conn)
			go handleIncomingRequest(decoder, encoder)
		}
	}
}

func handleIncomingRequest(decoder *json.Decoder, encoder *json.Encoder) {
	for {
		// store incoming data
		var data pb.Greet
		err := decoder.Decode(&data)
		if err != nil {
			log.Println(err)
			return
		}
		// fmt.Println(data)
		// respond
		encoder.Encode(greet.Resp{Message: "ok"})
	}

}
