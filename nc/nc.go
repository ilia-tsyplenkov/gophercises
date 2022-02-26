package nc

import "net"

type ncServer struct {
	connType string
	addr     string
	listener net.Listener
	conn     net.Conn
}

func (s *ncServer) Listen() (net.Listener, error) {
	var err error
	s.listener, err = net.Listen(s.connType, s.addr)
	return s.listener, err
}

func (s *ncServer) Accept() (net.Conn, error) {
	var err error
	s.conn, err = s.listener.Accept()
	return s.conn, err
}

func (s *ncServer) Close() {
	s.conn.Close()
	s.listener.Close()
}

func (s *ncServer) Read(p []byte) (int, error) {
	return s.conn.Read(p)
}
