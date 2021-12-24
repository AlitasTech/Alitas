// Package alitas implements a component that the TCP AppProxy can connect to.
package alitas

import (
	"fmt"
	"time"

	"github.com/AlitasTech/Alitas/src/proxy"
	"github.com/sirupsen/logrus"
)

// SocketAlitasProxy is a Golang implementation of a service that binds to a
// remote Alitas over an RPC/TCP connection. It implements handlers for the RPC
// requests sent by the SocketAppProxy, and submits transactions to Alitas via
// an RPC request. A SocketAlitasProxy can be implemented in any programming
// language as long as it implements the AppProxy interface over RPC.
type SocketAlitasProxy struct {
	nodeAddress string
	bindAddress string

	handler proxy.ProxyHandler

	client *SocketAlitasProxyClient
	server *SocketAlitasProxyServer
}

// NewSocketAlitasProxy creates a new SocketAlitasProxy
func NewSocketAlitasProxy(
	nodeAddr string,
	bindAddr string,
	handler proxy.ProxyHandler,
	timeout time.Duration,
	logger *logrus.Entry,
) (*SocketAlitasProxy, error) {

	if logger == nil {
		log := logrus.New()
		log.Level = logrus.DebugLevel
		logger = logrus.NewEntry(log)
	}

	client := NewSocketAlitasProxyClient(nodeAddr, timeout)

	server, err := NewSocketAlitasProxyServer(bindAddr, handler, timeout, logger)

	if err != nil {
		return nil, err
	}

	proxy := &SocketAlitasProxy{
		nodeAddress: nodeAddr,
		bindAddress: bindAddr,
		handler:     handler,
		client:      client,
		server:      server,
	}

	go proxy.server.listen()

	return proxy, nil
}

// SubmitTx submits a transaction to Alitas
func (p *SocketAlitasProxy) SubmitTx(tx []byte) error {
	ack, err := p.client.SubmitTx(tx)

	if err != nil {
		return err
	}

	if !*ack {
		error := "Failed to deliver transaction to Alitas"
		return fmt.Errorf("%s", error)
	}

	return nil
}
