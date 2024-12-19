package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"

	"hw3/internal/handlers"
	"hw3/internal/nodes"
	"hw3/internal/protos"
)

var (
	GRPCNodesAddress = []string{
		"localhost:5030",
		"localhost:5031",
		"localhost:5032",
		"localhost:5033",
		"localhost:5034",
	}
)

func main() {
	id := flag.Int("id", -1, "Node index")
	port := flag.Int("port", -1, "Node HTTP port")
	flag.Parse()

	if id == nil || *id < 0 || port == nil || *port < 0 {
		fmt.Println("Usage: <command> --id <node index> --port <node HTTP port>")
		os.Exit(1)
	}

	node := nodes.New(int64(*id), GRPCNodesAddress)

	go startGRPCServer(*id, node)
	go startHTTPServer(fmt.Sprintf("localhost:%s", strconv.Itoa(*port)), node)
	select {}
}

func startGRPCServer(id int, node *nodes.Node) {
	server := grpc.NewServer()
	address := GRPCNodesAddress[id]

	protos.RegisterNodeServer(server, node)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", address, err)
	}

	log.Printf("gRPC server listening on %s", address)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

func startHTTPServer(address string, node *nodes.Node) {
	r := mux.NewRouter()

	r.HandleFunc("/get", handlers.MakeGetHandler(node))
	r.HandleFunc("/update", handlers.MakeUpdateHandler(node))

	log.Printf("HTTP server listening on %s", address)
	if err := http.ListenAndServe(address, r); err != nil {
		log.Fatalf("Failed to serve HTTP server: %v", err)
	}
}
