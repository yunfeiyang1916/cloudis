package network

import (
	"context"
	"net"
)

// Conn 网络连接
type Conn struct {
	net.Conn
	// 连接id
	id int64
	// 父级上下文，当父级取消时，子级也可以取消
	parentCtx context.Context
}

func NewConn(id int64, conn net.Conn, srv *Server) *Conn {
	c := &Conn{
		id:        id,
		Conn:      conn,
		parentCtx: srv.ctx,
	}
	return c
}

func (c *Conn) Start() {

}
