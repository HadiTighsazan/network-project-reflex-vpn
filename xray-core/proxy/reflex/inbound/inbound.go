package inbound

import (
	"bufio"
	"context"

	"github.com/xtls/xray-core/common"
	"github.com/xtls/xray-core/common/errors"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/features/routing"
	"github.com/xtls/xray-core/proxy"
	"github.com/xtls/xray-core/proxy/reflex"
	"github.com/xtls/xray-core/proxy/reflex/codec"
	"github.com/xtls/xray-core/proxy/reflex/handshake"
	"github.com/xtls/xray-core/transport/internet/stat"
)

func init() {
	common.Must(common.RegisterConfig((*reflex.InboundConfig)(nil), func(ctx context.Context, cfg interface{}) (interface{}, error) {
		return New(ctx, cfg.(*reflex.InboundConfig))
	}))
}

var _ proxy.Inbound = (*Handler)(nil)

// Handler is a Reflex inbound handler.
type Handler struct {
	config    *reflex.InboundConfig
	validator reflex.Validator
	engine    *reflex.HandshakeEngine
}

func New(ctx context.Context, cfg *reflex.InboundConfig) (*Handler, error) {
	_ = ctx

	// Build an in-memory validator from config.
	mv := reflex.NewMemoryValidator()
	for _, c := range cfg.GetClients() {
		if c == nil {
			continue
		}
		if err := mv.AddFromConfig(c.GetId(), c.GetPolicy()); err != nil {
			return nil, errors.New("reflex inbound: invalid client").Base(err)
		}
	}

	eng := reflex.NewHandshakeEngine(mv)

	return &Handler{
		config:    cfg,
		validator: mv,
		engine:    eng,
	}, nil
}

func (h *Handler) Network() []net.Network {
	return []net.Network{net.Network_TCP}
}

func (h *Handler) Process(ctx context.Context, network net.Network, conn stat.Connection, dispatcher routing.Dispatcher) error {
	_ = dispatcher

	if network != net.Network_TCP {
		_ = conn.Close()
		return errors.New("reflex inbound: only supports TCP")
	}

	// Wrap connection for peek & parsing.
	reader := bufio.NewReader(conn)

	// Peek helps us decide if it "looks like" HTTP for error responses.
	peeked, _ := reader.Peek(64)
	looksHTTP := codec.LooksLikeHTTPPost(peeked)

	// Run Step2 handshake (magic or HTTP-like).
	_, err := h.engine.ServerDoHandshake(reader, conn)
	if err != nil {
		// If it's HTTP-like, try to respond with a normal-looking HTTP error.
		if looksHTTP {
			switch {
			case handshake.IsKind(err, handshake.KindUnauthenticated),
				handshake.IsKind(err, handshake.KindReplay):
				_ = reflex.WriteHTTPForbidden(conn)
			case handshake.IsKind(err, handshake.KindInvalidHandshake):
				_ = reflex.WriteHTTPBadRequest(conn)
			default:
				// Minimal 500-like response (avoid leaking details).
				_, _ = conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\nContent-Type: text/plain\r\nContent-Length: 0\r\n\r\n"))
			}
		}

		_ = conn.Close()

		// For Step2 we keep logs clean: return an info-level error.
		return errors.New("reflex inbound: handshake failed").Base(err).AtInfo()
	}

	// Step2 success: handshake completed and response sent.
	// Step3 will keep the connection and start encrypted transport.
	_ = conn.Close()
	return nil
}
