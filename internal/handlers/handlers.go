package handlers

import (
	"context"

	"github.com/loupe-co/{{.repoName}}/internal/config"
	"github.com/loupe-co/go-common/errors"
	"google.golang.org/grpc/codes"
)

var (
	ErrMethodNotImplemented = errors.New("method not implemented").Clean()
	ErrBadRequest           = errors.New("bad request").WithCode(codes.InvalidArgument).Clean()
)

// Handlers provide a generic interface for this service's various externally exposed functions (as defined in protos)
type Handlers struct {
	config          config.Config
}

// New returns a new instance of a Handlers struct
func New(cfg config.Config) *Handlers {
	return &Handlers{
		config:          cfg,
	}
}

func (h *Handlers) WithConfig(cfg config.Config) *Handlers {
	h.config = cfg
	return h
}
