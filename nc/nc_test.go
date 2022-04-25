package nc

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
)

const (
	testAddr = ":8080"
)

func TestServerListen(t *testing.T) {
	server := prepareServer("tcp", testAddr, t)
	server.listener.Close()

}

func TestServerAccept(t *testing.T) {
	server := prepareServer("tcp", testAddr, t)
	clientConn := prepareConnectionWithDefaultClient(server, t)

	clientConn.Close()
	server.Close()

}

func TestServerRead(t *testing.T) {
	server := prepareServer("tcp", testAddr, t)
	clientConn := prepareConnectionWithDefaultClient(server, t)

	defer func() {
		t.Log("closing client connection")
		clientConn.Close()
		t.Log("closing server ")
		server.Close()
		t.Log("all resources are free")
	}()

	want := "hello\n"

	fmt.Fprint(clientConn, want)
	got, err := bufio.NewReader(server).ReadString('\n')
	if err != nil {
		t.Fatalf("unexpected error while reading client request: %s\n", err)
	}
	if got != want {
		t.Fatalf("want to have %q but got %q\n", want, got)
	}
}

func TestClientDial(t *testing.T) {
	client := &Client{connType: "tcp", addr: testAddr}
	conn, srvConn, srvListener := prepareConnectionWithDefaultServer(client, t)
	defer func() {
		conn.Close()
		srvConn.Close()
		srvListener.Close()
	}()
	if conn == nil {
		t.Fatal("expected to have non nil connection but got nil\n")
	}
}

func TestClientWrite(t *testing.T) {
	client := &Client{connType: "tcp", addr: testAddr}
	conn, srvConn, srvListener := prepareConnectionWithDefaultServer(client, t)
	defer func() {
		conn.Close()
		srvConn.Close()
		srvListener.Close()
	}()
	want := "hello\n"
	fmt.Fprintf(conn, want)
	got, err := bufio.NewReader(srvConn).ReadString('\n')
	if err != nil {
		t.Fatalf("unexpected error while reading client request: %s\n", err)
	}
	if got != want {
		t.Fatalf("want to have %q but got %q\n", want, got)
	}

}

func prepareConnectionWithDefaultClient(server *Server, t *testing.T) net.Conn {
	t.Helper()
	var err error
	done := make(chan struct{})
	go func() {
		_, err = server.Accept()
		if err != nil {
			t.Fatalf("unexpected fail while accept connection: %s\n", err)
		}
		close(done)

	}()
	clientConn, err := net.Dial("tcp", testAddr)
	if err != nil {
		t.Fatalf("error of establishing connection to %s: %s\n", testAddr, err)
	}

	select {
	case <-time.After(10 * time.Millisecond):
		t.Fatal("connection hasn't been accepted by 'accept' function")
	case <-done:
		{
		}
	}
	return clientConn

}

func prepareConnectionWithDefaultServer(client *Client, t *testing.T) (net.Conn, net.Conn, net.Listener) {
	t.Helper()
	var err error
	done := make(chan struct{})
	ln, err := net.Listen("tcp", testAddr)
	var serverConn net.Conn
	if err != nil {
		t.Fatalf("error listening %q addr - %s\n", testAddr, err)
	}
	go func() {
		serverConn, err = ln.Accept()
		if err != nil {
			t.Fatalf("unexpected fail while accept connection: %s\n", err)
		}
		close(done)

	}()
	clientConn, err := client.Dial()
	if err != nil {
		t.Fatalf("error of establishing connection to %s: %s\n", testAddr, err)
	}

	select {
	case <-time.After(10 * time.Millisecond):
		t.Fatal("connection hasn't been accepted by 'accept' function")
	case <-done:
		{
		}
	}
	return clientConn, serverConn, ln

}
func prepareServer(connType, addr string, t *testing.T) *Server {
	t.Helper()
	server := &Server{connType: "tcp", addr: testAddr}
	_, err := server.Listen()
	if err != nil {
		t.Fatalf("unexpected fail to create a listener: %s\n", err)
	}
	return server
}
