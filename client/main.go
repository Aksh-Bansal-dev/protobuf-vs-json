package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/Aksh-Bansal-dev/grpc-vs-json/greet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr           = flag.String("addr", "localhost:50051", "the address to connect to")
	encodingFormat = flag.String("enc", "grpc", "encoding method (grpc/json)")
	dur            = flag.Int("t", 1, "time duration")
)

func main() {
	flag.Parse()
	timer := time.NewTimer(time.Second * time.Duration(*dur))
	startT := time.Now()
	successReq := 0
	if *encodingFormat == "grpc" {
		// Set up a connection to the server.
		conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewGreeterClient(conn)

		go func() {
			for {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				grpcReq(ctx, c)
				successReq++
				cancel()
			}
		}()
	} else {
		tcpAddr, err := net.ResolveTCPAddr("tcp", *addr)
		if err != nil {
			log.Fatal("ResolveTCPAddr failed:", err.Error())
		}

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		defer conn.Close()
		if err != nil {
			log.Fatal("Dial failed:", err.Error())
		}
		decoder := json.NewDecoder(conn)
		encoder := json.NewEncoder(conn)
		go func() {
			for {
				_, cancel := context.WithTimeout(context.Background(), time.Second)
				tcpReq(decoder, encoder)
				successReq++
				cancel()
			}
		}()
	}
	<-timer.C
	fmt.Println("Time: ", time.Since(startT).Milliseconds(), "ms")
	fmt.Println("Successful request: ", successReq)
}

func grpcReq(ctx context.Context, c pb.GreeterClient) {
	_, err := c.Auth(ctx, &pb.User{Name: "guy", Pass: "pas", Email: "acb@email.com", Age: 33})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
}

func tcpReq(decoder *json.Decoder, encoder *json.Encoder) {
	encoder.Encode(pb.Greet{Name: "guy", Pass: "pas", Email: "acb@email.com", Age: 33})
	var respBytes pb.Resp
	err := decoder.Decode(&respBytes)
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(respBytes)
}
