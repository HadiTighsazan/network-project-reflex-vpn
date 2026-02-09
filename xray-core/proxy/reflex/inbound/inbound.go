package inbound

import (
	"context"

	"github.com/xtls/xray-core/common"
	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/features/routing"
	"github.com/xtls/xray-core/proxy"
	"github.com/xtls/xray-core/proxy/reflex"
	"github.com/xtls/xray-core/transport/internet/stat"
)

func init() {
	common.Must(common.RegisterConfig((*reflex.InboundConfig)(nil), func(ctx context.Context, cfg interface{}) (interface{}, error) {
		return New(ctx, cfg.(*reflex.InboundConfig))
	}))
}

var _ proxy.Inbound = (*Handler)(nil)

// Handler is a Reflex inbound handler (Step1 stub).
type Handler struct {
	config *reflex.InboundConfig
}

func New(ctx context.Context, cfg *reflex.InboundConfig) (*Handler, error) {
	_ = ctx
	return &Handler{config: cfg}, nil
}

func (h *Handler) Network() []net.Network {
	return []net.Network{net.Network_TCP}
}

func (h *Handler) Process(ctx context.Context, network net.Network, conn stat.Connection, dispatcher routing.Dispatcher) error {
	_ = ctx
	_ = network
	_ = dispatcher

	// Step1 stub: close immediately to avoid leaking resources.
	_ = conn.Close()
	return errors.New("reflex inbound: not implemented (step1 stub)")
}
