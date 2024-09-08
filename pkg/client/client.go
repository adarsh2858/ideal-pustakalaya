package client

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/protobuf/proto"
	"github.com/loupe-co/go-common/errors"
	commonPub "github.com/loupe-co/go-common/pubsub"
	servicePb "github.com/loupe-co/protos/src/services/{{.repoName}}"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/grpc/credentials/insecure"
)

// ClientProtocol defines the type of messaging protocol to use for requests
type ClientProtocol uint8

const (
	ClientProtocolGRPC ClientProtocol = iota + 1
	ClientProtocolPubSub
)

// Options are the options for configuring the geyser connection
type Options struct {
	// Protocol is the mechanism by which the requests are sent
	// Supports 'grpc' and 'pubsub'
	Protocol ClientProtocol `json:"protocol"`

	// GRPCHost is the host address of the grpc service implementing the geyser grpc service
	GRPCHost string `json:"grpcHost"`

	// DialOptions provide grpc specific options to the grpc client
	DialOptions []grpc.DialOption `json:"dialOptions"`

	// RefreshTableTopic is the PubSub topic handling refresh table requests
	HelloTopic string `json:"helloTopic"`
}

// Op is a function which modifies the Client's options instance
type Op func(opts *Options)

// SetOptions allows you to pass in a full Client Options struct to configure the Client
func SetOptions(options Options) Op {
	return Op(func(opts *Options) {
		opts.Protocol = options.Protocol
		opts.GRPCHost = options.GRPCHost
		opts.DialOptions = options.DialOptions
		opts.HelloTopic = options.HelloTopic
	})
}

// SetProtocol sets which messaging protocol to use for client requests
// Supports Either grpc or pubsub
func SetProtocol(protocol ClientProtocol) Op {
	return Op(func(opts *Options) {
		opts.Protocol = protocol
	})
}

// SetGRPCHost sets the geyser grpc host for the grpc connection
func SetGRPCHost(host string) Op {
	return Op(func(opts *Options) {
		opts.GRPCHost = host
	})
}

// SetDialOptions sets the grpc DialOptions used to connect to the geyser grpc host
func SetDialOptions(dialOps ...grpc.DialOption) Op {
	return Op(func(opts *Options) {
		opts.DialOptions = dialOps
	})
}

// SetHelloTopic sets the RefreshTable topic string on the client options instance
func SetHelloTopic(topic string) Op {
	return Op(func(opts *Options) {
		opts.HelloTopic = topic
	})
}

// Client exposes geyser service methods to externally calling services via either grpc or pubsub
type Client struct {
	options      Options
	conn         *grpc.ClientConn
	grpcClient   servicePb.GeyserClient
	pubsubClient *commonPub.Client
}

// New returns a new instance of a geyser Client with the provided options
func New(ops ...Op) (*Client, error) {
	// Create default options instance
	options := Options{
		Protocol: ClientProtocolGRPC,
	}

	// Apply any user provided client ops
	for _, op := range ops {
		op(&options)
	}

	// Create client instance with options
	client := &Client{options: options}

	// Connect to geyser via the configured protocol
	var err error
	switch options.Protocol {
	case ClientProtocolGRPC:
		err = client.connectGRPC(context.Background())
	case ClientProtocolPubSub:
		err = client.connectPubSub(context.Background())
	default:
		err = errors.New("can't connect to unknown protocol")
	}

	// Handle any connect errors
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to {{.repoName}} via configured protocol")
	}

	return client, nil
}

// Hello calls the Hello method on the {{.repoName}} service, using the configured protocol
func (c *Client) Hello(ctx context.Context) error {
	return errors.New("method Hello not implemented").WithCode(codes.Unimplemented).Clean().AsGRPC()
}

func (c *Client) connectGRPC(ctx context.Context) error {
	// Make sure we don't already have a connection (may call this at times when we have a connection established already but don't know for sure)
	if c.conn != nil {
		return nil
	}

	// Get all grpc dial options, including our standard `grpc.WithInsecure()`
	dialOps := make([]grpc.DialOption, 0, len(c.options.DialOptions)+1)
	dialOps = append(dialOps, grpc.WithTransportCredentials(insecure.NewCredentials()))
	dialOps = append(dialOps, c.options.DialOptions...)

	// Establish grpc connection with configured grpc host and dial options
	conn, err := grpc.DialContext(ctx, c.options.GRPCHost, dialOps...)
	if err != nil {
		return errors.Wrap(err, "error dialing {{.repoName}} grpc host")
	}

	// Get {{.repoName}} client with grpc connection
	grpcClient := servicePb.New{{.repoName}}Client(conn)

	// Set the connection and client on the client wrapper
	c.conn = conn
	c.grpcClient = grpcClient

	return nil
}

func (c *Client) connectPubSub(ctx context.Context) error {
	// Make sure we don't already have a connection (may call this at times when we have a connection established already but don't know for sure)
	if c.pubsubClient != nil {
		return nil
	}

	// Create new pubsub client
	pubsubClient, err := commonPub.NewClient()
	if err != nil {
		return errors.Wrap(err, "error creating pubsub client")
	}

	// Set the pubsub client on the client wrapper
	c.pubsubClient = pubsubClient

	return nil
}
