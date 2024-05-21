package core

import (
	"fmt"
	"net"
	zinf "zinx/interface"
	"zinx/utils/conf"
	"zinx/utils/log"
)

type Server struct {
	Name string // server name
	Host string
	Port uint32
	AF   string // address family [IPv$, IPv6]

	msgHandler zinf.ZinfMsgHandler
}

func (sp *Server) Start() {
	// start a go routine for listening
	go func() {
		// try to resolve ip
		addr, resolveErr := net.ResolveTCPAddr(sp.AF, fmt.Sprintf("%s:%d", sp.Host, sp.Port))
		if resolveErr != nil {
			log.Erro(fmt.Sprintf("failed to resolve %s network address: %s:%d", sp.AF, sp.Host, sp.Port))
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
			newConn := NewConnection(conn, cid, sp.msgHandler)
			cid++

			go newConn.Start()
		}
	}()

	log.Info("server [%s] listening at %s:%d", sp.Name, sp.Host, sp.Port)
}

func (sp *Server) Stop() {
	log.Info(sp.Name, "stop...")
}

func (sp *Server) Serve() {
	sp.Start()
	select {} // block for loop
	sp.Stop()
}

func (sp *Server) AddRouter(msgId uint32, router zinf.ZinfRouter) {
	sp.msgHandler.AddRouter(msgId, router)
}

func NewServer() zinf.ZinfServer {
	s := &Server{
		Name: config.GlobalConfig.Name,
		Host: config.GlobalConfig.Host,
		Port: config.GlobalConfig.Port,
		AF:   "tcp4",

		msgHandler: NewMsgHandler(),
	}
	log.Info("%#v", *s)
	return s // yes, you do can return a pointer as interface type
}
