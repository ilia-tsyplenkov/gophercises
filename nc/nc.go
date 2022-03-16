package nc

import "net"

type NCserver struct {
	connType string
	addr     string
	listener net.Listener
	conn     net.Conn
}

func NewServer(addr, connType string) *NCserver {
	return &NCserver{connType: connType, addr: addr}
}
func (s *NCserver) Listen() (net.Listener, error) {
	var err error
	s.listener, err = net.Listen(s.connType, s.addr)
	return s.listener, err
}

func (s *NCserver) Accept() (net.Conn, error) {
	var err error
	s.conn, err = s.listener.Accept()
	return s.conn, err
}

func (s *NCserver) Close() {
	s.conn.Close()
	s.listener.Close()
}

func (s *NCserver) Read(p []byte) (int, error) {
	return s.conn.Read(p)
}
