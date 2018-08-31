package main

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	proto "github.com/amy/gophercon/version"
	"google.golang.org/grpc"
)

const (
	START     = "start"
	PROCESSED = "processed"
	UPDATING  = "updating"
	CONDEMNED = "condemned"
	FAILED    = "failed"
)

func main() {
	// Create gRPC client
	conn, _ := grpc.Dial(
		"server-connection:1234",
		grpc.WithInsecure(),
	)
	defer conn.Close()
	gRPCclient := proto.NewVersionServiceClient(conn)

	// Create k8s CRD client
	crdClient := NewPackageClient()

	// Create Client
	client := Client{
		version: "1.0", // Pull this version from --version flag in client yaml
	}

	// Call polling function every 5 seconds until stop channel closed
	stop := make(chan struct{})
	wait.JitterUntil(
		func() {
			var err error
			err = client.stateMachine(gRPCclient, crdClient)
			if err != nil {
				close(stop)
			}
		},
		time.Second*5,
		1.1,
		true,
		stop,
	)
}

type Client struct {
	version    string
	newVersion *proto.Package
	localState []Package
}

func (c *Client) stateMachine(client proto.VersionServiceClient, crdClient PackageClient) error {
	// Get new versions from server
	c.newVersion = fetchNewVersion(client)

	// Get local versions from CRD
	c.localState = fetchLocalState(crdClient)

	// Determine current state of Client
	var state string
	var name string
	var version string
	for _, pkg := range c.localState {
		if pkg.version == c.version {
			state = pkg.state
			name = pkg.name
			version = pkg.version
		}
	}

	switch state {
	case START:
		state = c.startState()
	case PROCESSED:
		state = c.processedState()
	case UPDATING:
		state = c.updatingState()
	case CONDEMNED:
		state = c.condemnedState()
	case FAILED:
		state = c.failedState()
	}

	// set look up CRD by name & version and update State
	crdClient.SetPackageState(name, version, state)

	return nil
}

// Call to Server
func fetchNewVersion(client proto.VersionServiceClient) *proto.Package {
	response, _ := client.GetVersion(context.TODO(), &proto.GetVersionRequest{})

	return response.Package
}

type Package struct {
	name    string
	version string
	state   string
	config  string
}

// Local state is read from CRDs
func fetchLocalState(packageClient PackageClient) []Package {
	return packageClient.GetAllPackages()
}

// handle transitions here
func (c *Client) startState() string {
	// This is the first client
	if len(c.localState) == 0 {
		// Create CRD
		// State of client is PROCESSED
		return PROCESSED
	}

	return UPDATING
}

func (c *Client) processedState() string {
	return ""
}

func (c *Client) updatingState() string {
	return ""
}

func (c *Client) condemnedState() string {
	return ""
}

func (c *Client) failedState() string {
	return ""
}

// something else deletes condemned & CRD
// something else deletes failed

type PackageClient struct{}

func (p PackageClient) GetAllPackages() []Package                   { return nil }
func (p PackageClient) SetPackageState(name, version, state string) {}

func NewPackageClient() PackageClient { return nil }
