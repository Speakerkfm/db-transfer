package services

import (
	"io"
	"net"
	"time"
)

type Conn struct {
	net.Conn
	IdleTimeout time.Duration
	MaxReadBuffer int64
}

func (c *Conn) Write(p []byte) (int, error) {
	c.updateDeadline()
	return c.Conn.Write(p)
}

func (c *Conn) Read(b []byte) (int, error) {
	c.updateDeadline()
	r := io.LimitReader(c.Conn, c.MaxReadBuffer)
	return r.Read(b)
}

func (c *Conn) updateDeadline() {
	idleDeadline := time.Now().Add(c.IdleTimeout)
	c.Conn.SetDeadline(idleDeadline)
}
