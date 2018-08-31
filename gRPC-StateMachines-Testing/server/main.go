package main

import (
	"context"
	"log"
	"net"

	proto "github.com/amy/gophercon/version"
	"google.golang.org/grpc"
)

func main() {
	grpc := grpc.NewServer()

	server := NewVersionServer()

	proto.RegisterVersionServiceServer(grpc, server)

	ln, err := net.Listen("tcp", "some address")
	if err != nil {
		log.Fatal(err)
	}

	go grpc.Serve(ln)
}

type versionServer struct {
}

func NewVersionServer() versionServer {
	return versionServer{}
}

func (v versionServer) GetVersion(context.Context, *proto.GetVersionRequest) (*proto.GetVersionResponse, error) {
	// This is built into the server
	// but another option is reading it froma configmap
	return &proto.GetVersionResponse{
		Package: &proto.Package{
			Name:    "client manifest",
			Version: "2",
			Config:  "some yaml",
		},
	}, nil

}
