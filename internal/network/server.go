package network

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
)

// Handler 每个建立连接后的处理器接口
type Handler interface {
	Handle(context.Context, *Conn)
}

// Option 配置函数，采用功能配置项来设置配置
type Option func(*options)

// options 服务配置项
type options struct {
	// Network 网络协议
	Network string
	// Address 监听地址
	Address string

	onConnect func(conn net.Conn)
}

// WithNetwork 设置网络协议
func WithNetwork(network string) Option {
	return func(o *options) {
		o.Network = network
	}
}

// WithAddress 设置监听地址
func WithAddress(address string) Option {
	return func(o *options) {
		o.Address = address
	}
}

// Server 服务
type Server struct {
	options
	// 监听器
	net.Listener
	sync.Once
	// 连接map
	connMap *sync.Map
	ctx     context.Context
	// identifier 连接最大id标识
	identifier int64
}

// NewServer 构建服务
func NewServer(opts ...Option) *Server {
	var o options
	for _, v := range opts {
		v(&o)
	}
	srv := Server{
		options: o,
	}
	return &srv
}

// Start 开启服务
func (s *Server) Start() error {
	var err error
	// 只执行一次
	s.Once.Do(func() {
		ls, e := net.Listen(s.Network, s.Address)
		if e != nil {
			err = e
			return
		}
		s.Listener = ls
		defer func() {
			// 在出现意外错误时需要关闭
			ls.Close()
		}()
		for {
			// Accept会一直阻塞直到有新的连接建立或者listen中断才会返回
			conn, e := ls.Accept()
			// 通常是由于listener被关闭无法继续监听导致的错误
			if e != nil {
				err = e
				return
			}
			connId := atomic.AddInt64(&s.identifier, 1)
			c := NewConn(connId, conn, s)
			go c.Start()
		}
	})
	return err
}
