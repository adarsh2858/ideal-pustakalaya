package grpc

import (
	"context"
	"math"

	"github.com/loupe-co/{{.repoName}}/internal/config"
	"github.com/loupe-co/{{.repoName}}/internal/handlers"
	"github.com/loupe-co/go-common"
	"github.com/loupe-co/go-common/errors"
	logGRPC "github.com/loupe-co/go-loupe-logger/grpc"
	servicePb "github.com/loupe-co/protos/src/services/{{.repoName}}"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	config   config.Config
	handlers *handlers.Handlers
}

func New(cfg config.Config, handles *handlers.Handlers) *common.GRPCServer {
	grpcHandlers := &GRPCServer{
		config:   cfg,
		handlers: handles,
	}

	grpcServer := common.NewGRPCServer(
		cfg.GRPC.Host,
		cfg.GRPC.Port,
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.MaxSendMsgSize(math.MaxInt32),
		logGRPC.GetInterceptorOption(false),
	)

	grpcServer.Register(func(server *grpc.Server) {
		servicePb.Register{{.repoName}}Server(server, grpcHandlers)
	})

	return grpcServer
}

func (server *GRPCServer) Hello(ctx context.Context, in *servicePb.HelloRequest) (*servicePb.Empty, error) {
	res, err := server.handlers.Hello(ctx, in)
	if err != nil {
		if cErr, ok := err.(errors.CommonError); ok {
			return nil, cErr.AsGRPC()
		}
		return nil, errors.Wrap(err, "error calling Hello rpc").Clean().AsGRPC()
	}
	return res, nil
}
