package network

import (
	"battle_rabbit/iface"
	"battle_rabbit/service/log"
	"fmt"
	"net"
	"sync"
	"time"
)

// TCPServer tcp服务器
type TCPServer struct {
	addr     string
	ln       net.Listener
	newAgent func(net.Conn) iface.IAgent // 连接代理
	wg       sync.WaitGroup
}

func NewTCPServer(addr string, newAgent func(net.Conn) iface.IAgent) *TCPServer {
	s := &TCPServer{
		addr:     addr,
		newAgent: newAgent,
	}
	s.Init()
	return s
}

func (s *TCPServer) SetNewAgentFn(fn func(c net.Conn) iface.IAgent) {
	s.newAgent = fn
}

func (s *TCPServer) Init() {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		panic(err)
	}

	if s.newAgent == nil {
		panic("NewAgent must not be nil")
	}
	log.Debug("start server  ", " succ, now listenning...")
	s.ln = ln
}
func (s *TCPServer) Start() {
	var tempDelay time.Duration
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Debug("accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return
		}
		tempDelay = 0

		fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())
		agent := s.newAgent(conn)
		go func() {
			defer agent.OnClose()
			s.wg.Add(1)

			agent.Run()
			s.wg.Done()
		}()
	}
}

func (s *TCPServer) Stop() {
	exit := make(chan *struct{})
	do := time.After(time.Second * 60)
	go func() {
		s.wg.Wait()
		exit <- nil
	}()
	select {
	case <-do:
	case <-exit:
	}

	log.Debug("== stop == tcp server !!")
	_ = s.ln.Close()
}
