package outbound

import (
	"context"

	"github.com/xtls/xray-core/common"
	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/proxy"
	"github.com/xtls/xray-core/proxy/reflex"
	"github.com/xtls/xray-core/transport"
	"github.com/xtls/xray-core/transport/internet"
)

func init() {
	common.Must(common.RegisterConfig((*reflex.OutboundConfig)(nil), func(ctx context.Context, cfg interface{}) (interface{}, error) {
		return New(ctx, cfg.(*reflex.OutboundConfig))
	}))
}

var _ proxy.Outbound = (*Handler)(nil)

// Handler is a Reflex outbound handler (Step1 stub).
type Handler struct {
	config *reflex.OutboundConfig
}

func New(ctx context.Context, cfg *reflex.OutboundConfig) (*Handler, error) {
	_ = ctx
	return &Handler{config: cfg}, nil
}

func (h *Handler) Process(ctx context.Context, link *transport.Link, dialer internet.Dialer) error {
	_ = ctx
	_ = link
	_ = dialer
	return errors.New("reflex outbound: not implemented (step1 stub)")
}
