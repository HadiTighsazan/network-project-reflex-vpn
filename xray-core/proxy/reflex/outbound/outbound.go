package outbound

import (
	"context"

	"github.com/xtls/xray-core/common"
	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/uuid"
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

// Handler is a Reflex outbound handler (Step2: handshake only).
type Handler struct {
	config   *reflex.OutboundConfig
	dest     net.Destination
	clientID [16]byte
	engine   *reflex.ClientHandshakeEngine
}

func New(ctx context.Context, cfg *reflex.OutboundConfig) (*Handler, error) {
	_ = ctx

	// Parse UUID once (fast + safer).
	id, err := uuid.ParseString(cfg.GetId())
	if err != nil {
		return nil, errors.New("reflex outbound: invalid id uuid").Base(err)
	}
	var idBytes [16]byte
	copy(idBytes[:], id.Bytes())

	dest := net.TCPDestination(net.ParseAddress(cfg.GetAddress()), net.Port(cfg.GetPort()))
	eng := reflex.NewClientHandshakeEngine(idBytes, cfg.GetAddress())

	return &Handler{
		config:   cfg,
		dest:     dest,
		clientID: idBytes,
		engine:   eng,
	}, nil
}

func (h *Handler) Process(ctx context.Context, link *transport.Link, dialer internet.Dialer) error {
	// Step2 focuses on handshake + auth. We don't forward payload yet (Step3).
	_ = link

	conn, err := dialer.Dial(ctx, h.dest)
	if err != nil {
		return errors.New("reflex outbound: dial failed").Base(err).AtWarning()
	}
	defer conn.Close()

	_, err = h.engine.DoHandshakeHTTP(conn)
	if err != nil {
		return errors.New("reflex outbound: handshake failed").Base(err).AtInfo()
	}

	// Step2 success: handshake completed.
	return nil
}
