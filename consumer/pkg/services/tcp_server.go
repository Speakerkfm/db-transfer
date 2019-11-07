package services

import (
	"bufio"
	"db-transfer/consumer/pkg/flags"
	"db-transfer/consumer/pkg/log"
	"github.com/mkideal/pkg/netutil/protocol"
	"net"
	"strconv"
	"time"
)

type TcpServer struct {
	l      net.Listener
	handle func([]byte) []byte
}

func NewTcpServer(conf *flags.Config, handle func([]byte) []byte) *TcpServer {
	l, err := net.Listen(protocol.TCP, conf.AppHost+":"+conf.AppPort)
	if err != nil {
		log.Error().Msg("Error listen connections")

	}
	defer l.Close()

	s := &TcpServer{
		l:      l,
		handle: handle,
	}

	log.Info().Msg("Listening on " + conf.AppHost + ":" + conf.AppPort)

	for {
		newConn, err := s.l.Accept()
		if err != nil {
			log.Error().Err(err).Msg("tcp error accepting")
		}

		conn := &Conn{
			Conn:          newConn,
			IdleTimeout:   5 * time.Second,
			MaxReadBuffer: 262144,
		}

		go tcpHandle(conn, s.handle)
	}

	return s
}

func tcpHandle(conn net.Conn, handle func([]byte) []byte) {
	msg, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		log.Error().Err(err).Msg("tcp port reading error")
	}

	log.Info().Msg("Message received. Len: " + strconv.Itoa(len(msg)))

	response := handle(msg)

	if _, err := conn.Write(response); err != nil {
		log.Error().Err(err).Msg("writing response error")
	}

	conn.Close()
}

func (c *TcpServer) Shutdown() error {
	if err := c.l.Close(); err != nil {
		return err
	}

	log.Info().Msg("TCP server shutdown OK")

	return nil
}
