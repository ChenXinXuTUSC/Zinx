package service

import (
	"fmt"
	"net"
	"time"
	"zinx/utils"
	zinf "zinx/interface"
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
			utils.Erro(fmt.Sprintf("failed to resolve %s network address: %s:%d", sp.AF, sp.IP, sp.Port))
			panic(resolveErr.Error())
		}
		// start listening
		listenner, listenErr := net.ListenTCP(sp.AF, addr)
		if listenErr != nil {
			utils.Erro("open listening error", listenErr.Error())
			panic(listenErr.Error())
		}

		for {
			// wait for client connection on accept sys invoke
			conn, acceptErr := listenner.AcceptTCP()
			if acceptErr != nil {
				utils.Erro("accept error", acceptErr.Error())
				continue
			}

			// TODO: set the maximum number of connection
			// TODO: handler for new connection

			// temporary echo service
			// rx: receive
			// tx: transport
			go func() {
				// start a go routine to serve the connection
				var clientAddrPort = conn.RemoteAddr().String()
				var buf []byte = make([]byte, 512)
				// infinity loop service
				for {
					rxLen, rxErr := conn.Read(buf)
					if rxErr != nil {
						utils.Erro(fmt.Sprintf("[%s] failed to read ", clientAddrPort), rxErr.Error())
						time.Sleep(1 * time.Second)
						break // end loop
					}
					utils.Dbug(clientAddrPort, "rx", rxLen)
					txLen, txErr := conn.Write(buf[:rxLen])
					if txErr != nil {
						utils.Erro(fmt.Sprintf("[%s] failed to write", clientAddrPort), txErr.Error())
						time.Sleep(1 * time.Second)
						continue
					}
					utils.Dbug(clientAddrPort, "tx", txLen)
				}
			}()
		}
	}()

	utils.Info("server [%s] listening at %s:%d", sp.Name, sp.IP, sp.Port)
}

func (sp *Server) Stop() {
	utils.Info(sp.Name, "stop...")
}

func (sp *Server) Serve() {
	sp.Start()
	select {} // block for loop
	sp.Stop()
}

func NewServer(name string, port int) zinf.ZInfServer {
	s := &Server{
		Name: name,
		AF: "tcp4",
		IP: "0.0.0.0",
		Port: port,
	}

	return s // yes, you do can return a pointer as interface type
}
