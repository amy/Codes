package main

import (
	"testing"

	"github.com/amy/gophercon/mock_version"
	proto "github.com/amy/gophercon/version"
	"github.com/golang/mock/gomock"
)

func TestMockingServerStub(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	gRPCclient := mock_version.NewMockVersionServiceClient(ctrl)
	gRPCclient.EXPECT().GetVersion(
		gomock.Any(),
		gomock.Any(),
	).Return(&proto.GetVersionResponse{})

	client := Client{
		version: "1",
		newVersion: &proto.Package{
			Name:    "Gophercon",
			Version: "2",
			Config:  "Some Kubernetes Config",
		},
	}
	if err := client.stateMachine(gRPCclient, PackageClient{}); err != nil {
		t.Error("unexpected error")
	}
}

// Unit Testing Transitions
func TestStartToProcessed(t *testing.T) {
	client := Client{
		version:    "1",
		newVersion: &proto.Package{},
		localState: []Package{},
	}

	if state := client.startState(); state != "PROCESSED" {
		t.Errorf("Expected state PROCESSED. Got state %s.", state)
	}
}
