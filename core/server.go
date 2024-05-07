package core

import (
	"fmt"
	"net"
	zinf "zinx/interface"
	"zinx/service"
	"zinx/utils/log"
)

type Server struct {
	Name string // server name
	AF   string // address family [IPv$, IPv6]
	IP   string
	Port int
}

func (sp *Server) Start() {
	// start a go routine for listening
	go func() {
		// try to resolve ip
		addr, resolveErr := net.ResolveTCPAddr(sp.AF, fmt.Sprintf("%s:%d", sp.IP, sp.Port))
		if resolveErr != nil {
			log.Erro(fmt.Sprintf("failed to resolve %s network address: %s:%d", sp.AF, sp.IP, sp.Port))
			panic(resolveErr.Error())
		}
		// start listening
		listenner, listenErr := net.ListenTCP(sp.AF, addr)
		if listenErr != nil {
			log.Erro("open listening error", listenErr.Error())
			panic(listenErr.Error())
		}

		// TODO: a algorithmn for generating connID
		var cid uint32 = 0
		for {
			// wait for client connection on accept sys invoke
			conn, acceptErr := listenner.AcceptTCP()
			if acceptErr != nil {
				log.Erro("accept error", acceptErr.Error())
				continue
			}

			// TODO: set the maximum number of connection
			// bind the connection with a service handler function
			newConn := NewConnection(conn, cid, service.EchoBack)
			cid++

			go newConn.Start()
		}
	}()

	log.Info("server [%s] listening at %s:%d", sp.Name, sp.IP, sp.Port)
}

func (sp *Server) Stop() {
	log.Info(sp.Name, "stop...")
}

func (sp *Server) Serve() {
	sp.Start()
	select {} // block for loop
	sp.Stop()
}

func NewServer(name string, port int) zinf.ZinfServer {
	s := &Server{
		Name: name,
		AF:   "tcp4",
		IP:   "0.0.0.0",
		Port: port,
	}

	return s // yes, you do can return a pointer as interface type
}
