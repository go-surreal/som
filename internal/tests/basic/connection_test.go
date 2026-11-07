package basic

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"som.test/gen/som/query"
	"som.test/gen/som/repo"
	"som.test/model"
)

func TestHealth(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	if err := client.Health(ctx); err != nil {
		t.Fatalf("expected healthy server, got: %v", err)
	}
}

func TestReconnect(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	proxy := newDropProxy(t, sharedEndpoint)
	statusEvents := make(chan repo.ConnectionEvent, 64)

	client, cleanup := prepareDatabaseWith(ctx, t, func(config *repo.Config) {
		config.Address = "ws://" + proxy.Address()
		config.Reconnect = repo.ReconnectConfig{
			MinDelay: 10 * time.Millisecond,
			MaxDelay: 100 * time.Millisecond,
		}
		config.OnStatus = func(event repo.ConnectionEvent) {
			statusEvents <- event
		}
	})
	defer cleanup()

	liveCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	liveChan, err := client.LocationRepo().Query().Live(liveCtx)
	if err != nil {
		t.Fatal(err)
	}

	proxy.DropConnections()

	expectLiveDisconnected(t, liveChan)
	expectStatus(t, statusEvents, repo.StatusDisconnected)
	expectStatus(t, statusEvents, repo.StatusConnected)

	// INFO FOR DB fails unless namespace, database and authentication were all
	// restored on the new connection.
	if _, err := client.Raw(ctx, "INFO FOR DB", nil); err != nil {
		t.Fatalf("expected client to be usable after reconnect, got: %v", err)
	}
}

func expectLiveDisconnected(t *testing.T, liveChan <-chan query.LiveResult[*model.Location]) {
	t.Helper()

	select {

	case result, more := <-liveChan:
		if !more {
			t.Fatal("live channel closed without a disconnected event")
		}

		if _, ok := result.(query.LiveDisconnected[*model.Location]); !ok {
			t.Fatalf("expected LiveDisconnected event, got %T", result)
		}

	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for disconnected event after dropping the connection")
	}
}

func expectStatus(t *testing.T, statusEvents <-chan repo.ConnectionEvent, want repo.ConnectionStatus) {
	t.Helper()

	timeout := time.After(5 * time.Second)

	for {
		select {

		case event := <-statusEvents:
			if event.Status == want {
				return
			}

		case <-timeout:
			t.Fatalf("timeout waiting for connection status %s", want)
		}
	}
}

// dropProxy forwards TCP connections to a target address and can cut all of
// them at once, simulating a lost connection without closing the client.
type dropProxy struct {
	listener net.Listener
	target   string

	mu    sync.Mutex
	conns []net.Conn
}

func newDropProxy(tb testing.TB, target string) *dropProxy {
	tb.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		tb.Fatalf("failed to start proxy listener: %v", err)
	}

	proxy := &dropProxy{
		listener: listener,
		target:   target,
	}

	tb.Cleanup(func() {
		_ = listener.Close()
		proxy.DropConnections()
	})

	go proxy.acceptLoop()

	return proxy
}

func (p *dropProxy) Address() string {
	return p.listener.Addr().String()
}

func (p *dropProxy) DropConnections() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, conn := range p.conns {
		_ = conn.Close()
	}

	p.conns = nil
}

func (p *dropProxy) acceptLoop() {
	for {
		client, err := p.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			continue
		}

		go p.forward(client)
	}
}

func (p *dropProxy) forward(client net.Conn) {
	upstream, err := net.Dial("tcp", p.target)
	if err != nil {
		_ = client.Close()
		return
	}

	p.mu.Lock()
	p.conns = append(p.conns, client, upstream)
	p.mu.Unlock()

	go pipe(upstream, client)
	pipe(client, upstream)
}

func pipe(destination, source net.Conn) {
	_, _ = io.Copy(destination, source)
	_ = destination.Close()
	_ = source.Close()
}
