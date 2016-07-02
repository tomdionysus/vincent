package server

import(
  "net"
  "time"
)

type ConnLimitListener struct {
  listener *net.TCPListener
  pool chan bool
}

func NewConnLimitListener(count int, l *net.TCPListener) net.Listener {
  pool := make(chan bool, count)
  for i := 0; i < count; i++ { pool <- true }

  return &ConnLimitListener{
    listener: l,
    pool: pool,
  }
}

func (me *ConnLimitListener) Addr() net.Addr { return me.listener.Addr() }
func (me *ConnLimitListener) Close() error { return me.listener.Close() }
func (me *ConnLimitListener) Accept() (net.Conn, error) {
  <-me.pool 
  tc, err := me.listener.AcceptTCP()
  if err != nil { return nil, err }
  tc.SetKeepAlive(true)
  tc.SetKeepAlivePeriod(10 * time.Second)
  conn := ConnLimitConn{conn: tc, closechan: me.pool }
  return &conn, nil
}

// Connection

type ConnLimitConn struct {
  conn net.Conn
  closechan chan bool
}

func (me *ConnLimitConn) Read(b []byte) (int, error) { return me.conn.Read(b) }
func (me *ConnLimitConn) Write(b []byte) (int, error) { return me.conn.Write(b) }
func (me *ConnLimitConn) Close() error {
  err := me.conn.Close()
  me.closechan <- true 
  if err!=nil { return err }
  return nil
}
func (me *ConnLimitConn) LocalAddr() net.Addr { return me.conn.LocalAddr() }
func (me *ConnLimitConn) RemoteAddr() net.Addr { return me.conn.RemoteAddr() }
func (me *ConnLimitConn) SetDeadline(t time.Time) error { return me.conn.SetDeadline(t) }
func (me *ConnLimitConn) SetReadDeadline(t time.Time) error { return me.conn.SetReadDeadline(t) }
func (me *ConnLimitConn) SetWriteDeadline(t time.Time) error { return me.conn.SetWriteDeadline(t) }