package smtp

import (
	"io"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/doctolib/MailHog/pkg/config"
)

func Listen(cfg *config.Config, exitCh chan int) *net.TCPListener {
	log.Infof("[SMTP] Binding to address: %s\n", cfg.SMTPBindAddr)
	ln, err := net.Listen("tcp", cfg.SMTPBindAddr)
	if err != nil {
		log.Fatalf("[SMTP] Error listening on socket: %s\n", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Errorf("[SMTP] Error accepting connection: %s\n", err)
			continue
		}

		if cfg.Monkey != nil {
			ok := cfg.Monkey.Accept(conn)
			if !ok {
				conn.Close()
				continue
			}
		}

		go Accept(
			conn.(*net.TCPConn).RemoteAddr().String(),
			io.ReadWriteCloser(conn),
			cfg.Storage,
			cfg.MessageChan,
			cfg.Hostname,
			cfg.Monkey,
		)
	}
}
