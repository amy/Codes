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
	server := NewMockVersionServer()
	proto.RegisterVersionServiceServer(grpc, server)

	ln, err := net.Listen("tcp", "some address")
	if err != nil {
		log.Fatal(err)
	}

	go grpc.Serve(ln)
}

type mockVersionServer struct {
	testCases map[string]map[int][]proto.Package // test case : # server hits : packages to return
}

func NewMockVersionServer() mockVersionServer {
	tests := make(map[string]map[int][]proto.Package)

	tests["update"] = map[int][]proto.Package{
		-1: []proto.Package{},
		0:  []proto.Package{},
	}

	tests["fail"] = map[int][]proto.Package{
		-1: []proto.Package{
			proto.Package{
				Name:    "a failed package",
				Version: "2",
				Config:  "some config",
			},
		},
		0: []proto.Package{},
	}

	return mockVersionServer{}
}

func (v mockVersionServer) GetVersion(context.Context, *proto.GetVersionRequest) (*proto.GetVersionResponse, error) {

	return &proto.GetVersionResponse{}, nil
}
