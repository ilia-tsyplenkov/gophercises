package nc

import "net"

type Server struct {
	connType string
	addr     string
	listener net.Listener
	conn     net.Conn
}

func NewServer(addr, connType string) *Server {
	return &Server{connType: connType, addr: addr}
}

func NewClient(addr, connType string) *Client {
	return &Client{connType: connType, addr: addr}
}

func (s *Server) Listen() (net.Listener, error) {
	var err error
	s.listener, err = net.Listen(s.connType, s.addr)
	return s.listener, err
}

func (s *Server) Accept() (net.Conn, error) {
	var err error
	s.conn, err = s.listener.Accept()
	return s.conn, err
}

func (s *Server) Close() {
	s.conn.Close()
	s.listener.Close()
}

func (s *Server) Read(p []byte) (int, error) {
	return s.conn.Read(p)
}

type Client struct {
	connType string
	addr     string
	conn     net.Conn
}

func (c *Client) Dial() (net.Conn, error) {
	conn, err := net.Dial(c.connType, c.addr)
	c.conn = conn
	return conn, err
}

func (c *Client) Write(p []byte) (int, error) {
	return c.conn.Write(p)
}
