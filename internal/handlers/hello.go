package handlers

import (
	"context"
	"fmt"

	servicePb "github.com/loupe-co/protos/src/services/{{.repoName}}"
	"github.com/loupe-co/{{.repoName}}/internal/config"
	"github.com/loupe-co/go-common/errors"
	"github.com/loupe-co/go-loupe-logger/log"
)

func (h *Handlers) Hello(ctx context.Context, in *servicePb.HelloRequest) (*servicePb.Empty, error) {
	logger := log.WithCustom("name", in.GetName())

	if len(in.GetName()) == 0 {
		return &servicePb.Empty{}, ErrBadRequest
	}

	helloMsg := fmt.Sprintf("hello %s", in.GetName())

	logger.Debug(helloMsg)
	fmt.Println(helloMsg)

	return &servicePb.Empty{}, nil
}
