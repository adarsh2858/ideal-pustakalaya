package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"google.golang.org/protobuf/proto"
	"github.com/loupe-co/{{.repoName}}/internal/config"
	"github.com/loupe-co/{{.repoName}}/internal/handlers"
	"github.com/loupe-co/go-common/errors"
	commonPub "github.com/loupe-co/go-common/pubsub"
	"github.com/loupe-co/go-loupe-logger/log"
	servicePb "github.com/loupe-co/protos/src/services/{{.repoName}}"
)

func New(ctx context.Context, cfg config.Config, handles *handlers.Handlers) (*commonPub.Client, error) {
	pubsubServer, err := commonPub.NewClient(
		commonPub.MaxMessages(cfg.PubSub.MaxMessages),
		commonPub.MaxRetries(cfg.PubSub.MaxRetries),
		commonPub.AckTimeout(cfg.PubSub.AckTimeout),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating pubsub server")
	}

	if cfg.ProjectID != "local" {
		if err := pubsubServer.Handle(ctx, cfg.PubSub.HelloSub, cfg.PubSub.HelloTopic, Hello(handles)); err != nil {
			return nil, errors.Wrap(err, "error setting up Hello pubsub handler")
		}
	}

	return pubsubServer, nil
}

func Hello(handles *handlers.Handlers) func(ctx context.Context, msg *pubsub.Message) error {
	return func(ctx context.Context, msg *pubsub.Message) error {
		req := &servicePb.HelloRequest{}
		if err := proto.Unmarshal(msg.Data, req); err != nil {
			err := errors.Wrap(err, "error reading hello request")
			log.Error(err)
			return err
		}
		if _, err := handles.Hello(ctx, req); err != nil {
			log.Error(err)
			return err
		}
		return nil
	}
}
