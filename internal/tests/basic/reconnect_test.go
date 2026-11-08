package basic

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"som.test/gen/som/query"
	"som.test/gen/som/repo"
	"som.test/model"
)

const (
	reconnectInterval = 100 * time.Millisecond
	reconnectTimeout  = 10 * time.Second
)

func TestPing(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, cleanup := prepareDatabase(ctx, t)
	defer cleanup()

	if err := client.Ping(ctx); err != nil {
		t.Fatalf("expected ping to succeed, got: %v", err)
	}
}

func TestPingAfterClose(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	client, _ := prepareDatabase(ctx, t)
	client.Close()

	if err := client.Ping(ctx); err == nil {
		t.Fatal("expected ping to fail on a closed client")
	}
}

func TestReconnectRequiresWebSocket(t *testing.T) {
	t.Parallel()

	_, err := repo.NewClient(context.Background(), repo.Config{
		Address:           "http://" + sharedEndpoint,
		ReconnectInterval: reconnectInterval,
	})
	if err == nil {
		t.Fatal("expected an error for reconnect over http")
	}
}

func TestReconnect(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	proxy := startFlakyProxy(t, sharedEndpoint)

	client, cleanup := prepareDatabaseWith(ctx, t, func(config *repo.Config) {
		config.Address = "ws://" + proxy.address()
		config.ReconnectInterval = reconnectInterval
	})
	defer cleanup()

	proxy.setDown(true)

	waitFor(t, "ping to fail while the database is unreachable", func() bool {
		return pingWithTimeout(ctx, client) != nil
	})

	proxy.setDown(false)

	waitFor(t, "ping to succeed after reconnecting", func() bool {
		return pingWithTimeout(ctx, client) == nil
	})

	record := &model.AllTypes{FieldString: "after reconnect", FieldMonth: time.January}

	if err := client.AllTypesRepo().Create(ctx, record); err != nil {
		t.Fatalf("failed to create record after reconnecting: %v", err)
	}
}

func TestReconnectFailsInFlightRequest(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	proxy := startFlakyProxy(t, sharedEndpoint)

	client, cleanup := prepareDatabaseWith(ctx, t, func(config *repo.Config) {
		config.Address = "ws://" + proxy.address()
		config.ReconnectInterval = reconnectInterval
	})
	defer cleanup()

	result := make(chan error, 1)

	go func() {
		_, err := client.Raw(ctx, "SLEEP 5s", nil)
		result <- err
	}()

	time.Sleep(500 * time.Millisecond)

	proxy.cut()

	select {

	case err := <-result:
		if err == nil {
			t.Fatal("expected the in-flight request to fail after the connection was lost")
		}

	case <-time.After(3 * time.Second):
		t.Fatal("in-flight request did not return after the connection was lost")
	}
}

func TestReconnectRestoresLiveQuery(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	proxy := startFlakyProxy(t, sharedEndpoint)

	client, cleanup := prepareDatabaseWith(ctx, t, func(config *repo.Config) {
		config.Address = "ws://" + proxy.address()
		config.ReconnectInterval = reconnectInterval
	})
	defer cleanup()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	liveChan, err := client.AllTypesRepo().Query().Live(ctx)
	if err != nil {
		t.Fatal(err)
	}

	proxy.cut()

	record := &model.AllTypes{FieldString: "created after reconnect", FieldMonth: time.January}

	waitFor(t, "record to be created after reconnecting", func() bool {
		return client.AllTypesRepo().Create(ctx, record) == nil
	})

	select {

	case res, more := <-liveChan:
		if !more {
			t.Fatal("live channel closed instead of being restored")
		}

		created, ok := res.(query.LiveCreate[*model.AllTypes])
		if !ok {
			t.Fatalf("expected LiveCreate event, got %T", res)
		}

		createdRecord, err := created.Get()
		if err != nil {
			t.Fatal(err)
		}

		if createdRecord.FieldString != record.FieldString {
			t.Fatalf("expected record %q, got %q", record.FieldString, createdRecord.FieldString)
		}

	case <-time.After(reconnectTimeout):
		t.Fatal("timeout waiting for live notification after reconnecting")
	}
}

// TestReconnectWithExpiredToken checks that a reconnect signs in with the
// configured credentials when the token of the original sign-in has expired.
func TestReconnectWithExpiredToken(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	admin, cleanupAdmin := prepareDatabase(ctx, t)
	defer cleanupAdmin()

	username := "short_" + gofakeit.LetterN(12)
	password := gofakeit.Password(true, true, true, false, false, 32)

	defineUser := fmt.Sprintf(
		"DEFINE USER %s ON ROOT PASSWORD '%s' ROLES OWNER DURATION FOR TOKEN 1s, FOR SESSION NONE",
		username, password,
	)

	if _, err := admin.Raw(ctx, defineUser, nil); err != nil {
		t.Fatal(err)
	}

	defer func() {
		_, _ = admin.Raw(context.Background(), "REMOVE USER "+username+" ON ROOT", nil)
	}()

	proxy := startFlakyProxy(t, sharedEndpoint)

	client, cleanup := prepareDatabaseWith(ctx, t, func(config *repo.Config) {
		config.Address = "ws://" + proxy.address()
		config.Username = username
		config.Password = password
		config.ReconnectInterval = reconnectInterval
	})
	defer cleanup()

	time.Sleep(2 * time.Second)

	proxy.cut()

	record := &model.AllTypes{FieldString: "after token expiry", FieldMonth: time.January}

	waitFor(t, "record to be created after reconnecting with an expired token", func() bool {
		return client.AllTypesRepo().Create(ctx, record) == nil
	})
}

func pingWithTimeout(ctx context.Context, client repo.Client) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	return client.Ping(ctx)
}

func waitFor(t *testing.T, description string, condition func() bool) {
	t.Helper()

	deadline := time.Now().Add(reconnectTimeout)

	for time.Now().Before(deadline) {
		if condition() {
			return
		}

		time.Sleep(reconnectInterval)
	}

	t.Fatalf("timeout waiting for %s", description)
}

// flakyProxy forwards TCP connections to a target and can cut them on demand,
// to simulate a database that becomes unreachable.
type flakyProxy struct {
	listener net.Listener
	target   string

	mu          sync.Mutex
	connections []net.Conn
	isDown      bool
}

func startFlakyProxy(t *testing.T, target string) *flakyProxy {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start proxy: %v", err)
	}

	proxy := &flakyProxy{
		listener: listener,
		target:   target,
	}

	go proxy.acceptLoop()

	t.Cleanup(func() {
		_ = listener.Close()
		proxy.cut()
	})

	return proxy
}

func (p *flakyProxy) address() string {
	return p.listener.Addr().String()
}

// setDown cuts all open connections and refuses new ones while isDown is set.
func (p *flakyProxy) setDown(isDown bool) {
	p.mu.Lock()
	p.isDown = isDown
	p.mu.Unlock()

	if isDown {
		p.cut()
	}
}

func (p *flakyProxy) cut() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, conn := range p.connections {
		_ = conn.Close()
	}

	p.connections = nil
}

func (p *flakyProxy) acceptLoop() {
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

func (p *flakyProxy) forward(client net.Conn) {
	upstream, err := net.Dial("tcp", p.target)
	if err != nil {
		_ = client.Close()
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isDown {
		_ = client.Close()
		_ = upstream.Close()
		return
	}

	p.connections = append(p.connections, client, upstream)

	go pipe(upstream, client)
	go pipe(client, upstream)
}

func pipe(destination, source net.Conn) {
	_, _ = io.Copy(destination, source)
	_ = destination.Close()
	_ = source.Close()
}
