package mongo

import (
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

type failingCloseConn struct {
	net.Conn
	calls atomic.Int32
	err   error
}

func (conn *failingCloseConn) Close() error {
	conn.calls.Add(1)
	return conn.err
}

func TestConnConcurrentClose(t *testing.T) {
	closeErr := errors.New("close failed")
	underlying := &failingCloseConn{err: closeErr}
	conn := newConnWith(underlying, nil)
	var workers sync.WaitGroup
	for range 32 {
		workers.Go(func() {
			if err := conn.Close(); !errors.Is(err, closeErr) {
				t.Errorf("Close() = %v, want %v", err, closeErr)
			}
		})
	}
	workers.Wait()
	if calls := underlying.calls.Load(); calls != 1 {
		t.Fatalf("underlying Close called %d times, want 1", calls)
	}
}

func TestServerStopReleasesListenerOnConnectionError(t *testing.T) {
	server := NewServer().(*server)
	server.SetAddress("127.0.0.1")
	server.SetPort(0)
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = server.Stop() })
	addr := server.tcpListener.Addr().String()
	closeErr := errors.New("close failed")
	conn := newConnWith(&failingCloseConn{err: closeErr}, nil)
	server.AddConn(conn)
	if err := server.Stop(); !errors.Is(err, closeErr) {
		t.Fatalf("Stop() = %v, want %v", err, closeErr)
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatalf("listener was not released: %v", err)
	}
	defer listener.Close()
	_ = server.RemoveConn(conn)
}

func TestServerRestartWithActiveConnections(t *testing.T) {
	server := NewServer().(*server)
	server.SetAddress("127.0.0.1")
	server.SetPort(0)
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := server.Stop(); err != nil {
			t.Error(err)
		}
	})
	addr := server.tcpListener.Addr().String()
	for range 20 {
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err != nil {
			t.Fatal(err)
		}
		if err := server.Restart(); err != nil {
			_ = conn.Close()
			t.Fatal(err)
		}
		if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			t.Fatal(err)
		}
		if _, err := conn.Read(make([]byte, 1)); !errors.Is(err, io.EOF) && !errors.Is(err, syscall.ECONNRESET) {
			t.Errorf("read after Stop = %v, want EOF or connection reset", err)
		}
		_ = conn.Close()
	}
}

func TestServerStopWithIdleConnection(t *testing.T) {
	server := NewServer().(*server)
	// A peer that sends no request must not keep shutdown waiting.
	server.SetAddress("127.0.0.1")
	server.SetPort(0)
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	conn, err := net.DialTimeout("tcp", server.tcpListener.Addr().String(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	done := make(chan error, 1)
	go func() { done <- server.Stop() }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Stop blocked on an idle connection")
	}
}
