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
	clientConn := prepareConnections(server, t)

	clientConn.Close()
	server.Close()

}

func TestServerRead(t *testing.T) {
	server := prepareServer("tcp", testAddr, t)
	clientConn := prepareConnections(server, t)

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

func prepareConnections(server *NCserver, t *testing.T) net.Conn {
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
func prepareServer(connType, addr string, t *testing.T) *NCserver {
	t.Helper()
	server := &NCserver{connType: "tcp", addr: testAddr}
	_, err := server.Listen()
	if err != nil {
		t.Fatalf("unexpected fail to create a listener: %s\n", err)
	}
	return server
}
