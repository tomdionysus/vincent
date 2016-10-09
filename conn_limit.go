package vincent

import (
	"net"
	"time"
)

// A TCPListener that limits concurrent connections
type ConnLimitListener struct {
	listener *net.TCPListener
	pool     chan bool
}

// Return a new ConnLimitListener using the supplied connection limit and underlying TCPListener
func NewConnLimitListener(count int, l *net.TCPListener) net.Listener {
	pool := make(chan bool, count)
	for i := 0; i < count; i++ {
		pool <- true
	}

	return &ConnLimitListener{
		listener: l,
		pool:     pool,
	}
}

// Return the underlying listener Address
func (cll *ConnLimitListener) Addr() net.Addr { return cll.listener.Addr() }

// Close the underlying listener
func (cll *ConnLimitListener) Close() error { return cll.listener.Close() }

// Block until a connection is available and the limit has not been reached, then
// accpt the connection and return it
func (cll *ConnLimitListener) Accept() (net.Conn, error) {
	<-cll.pool
	tc, err := cll.listener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(10 * time.Second)
	conn := ConnLimitConn{conn: tc, closechan: cll.pool}
	return &conn, nil
}

// A net.Conn for use with ConnLimitListener
type ConnLimitConn struct {
	conn      net.Conn
	closechan chan bool
}

// Read from the underlying connection
func (cll *ConnLimitConn) Read(b []byte) (int, error) { return cll.conn.Read(b) }

// Write to the underlying connection
func (cll *ConnLimitConn) Write(b []byte) (int, error) { return cll.conn.Write(b) }

// Close the underlying connection
func (cll *ConnLimitConn) Close() error {
	err := cll.conn.Close()
	cll.closechan <- true
	if err != nil {
		return err
	}
	return nil
}

// Return the LocalAddr of the underlying connection
func (cll *ConnLimitConn) LocalAddr() net.Addr { return cll.conn.LocalAddr() }

// Return the RemoteAddr of the underlying connection
func (cll *ConnLimitConn) RemoteAddr() net.Addr { return cll.conn.RemoteAddr() }

// Set the timeout deadline of the underlying connection
func (cll *ConnLimitConn) SetDeadline(t time.Time) error { return cll.conn.SetDeadline(t) }

// Set the read timeout deadline of the underlying connection
func (cll *ConnLimitConn) SetReadDeadline(t time.Time) error { return cll.conn.SetReadDeadline(t) }

// Set the write timeout deadline of the underlying connection
func (cll *ConnLimitConn) SetWriteDeadline(t time.Time) error { return cll.conn.SetWriteDeadline(t) }
