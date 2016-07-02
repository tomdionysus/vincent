package server

import(
  "net"
  "time"
)

// A TCPListener that limits concurrent connections
type ConnLimitListener struct {
  listener *net.TCPListener
  pool chan bool
}

// Return a new ConnLimitListener using the supplied connection limit and underlying TCPListener
func NewConnLimitListener(count int, l *net.TCPListener) net.Listener {
  pool := make(chan bool, count)
  for i := 0; i < count; i++ { pool <- true }

  return &ConnLimitListener{
    listener: l,
    pool: pool,
  }
}

// Return the underlying listener Address
func (me *ConnLimitListener) Addr() net.Addr { return me.listener.Addr() }
// Close the underlying listener
func (me *ConnLimitListener) Close() error { return me.listener.Close() }
// Block until a connection is available and the limit has not been reached, then
// accpt the connection and return it
func (me *ConnLimitListener) Accept() (net.Conn, error) {
  <-me.pool 
  tc, err := me.listener.AcceptTCP()
  if err != nil { return nil, err }
  tc.SetKeepAlive(true)
  tc.SetKeepAlivePeriod(10 * time.Second)
  conn := ConnLimitConn{conn: tc, closechan: me.pool }
  return &conn, nil
}

// A net.Conn for use with ConnLimitListener
type ConnLimitConn struct {
  conn net.Conn
  closechan chan bool
}

// Read from the underlying connection
func (me *ConnLimitConn) Read(b []byte) (int, error) { return me.conn.Read(b) }
// Write to the underlying connection
func (me *ConnLimitConn) Write(b []byte) (int, error) { return me.conn.Write(b) }
// Close the underlying connection
func (me *ConnLimitConn) Close() error {
  err := me.conn.Close()
  me.closechan <- true 
  if err!=nil { return err }
  return nil
}
// Return the LocalAddr of the underlying connection
func (me *ConnLimitConn) LocalAddr() net.Addr { return me.conn.LocalAddr() }
// Return the RemoteAddr of the underlying connection
func (me *ConnLimitConn) RemoteAddr() net.Addr { return me.conn.RemoteAddr() }
// Set the timeout deadline of the underlying connection
func (me *ConnLimitConn) SetDeadline(t time.Time) error { return me.conn.SetDeadline(t) }
// Set the read timeout deadline of the underlying connection
func (me *ConnLimitConn) SetReadDeadline(t time.Time) error { return me.conn.SetReadDeadline(t) }
// Set the write timeout deadline of the underlying connection
func (me *ConnLimitConn) SetWriteDeadline(t time.Time) error { return me.conn.SetWriteDeadline(t) }