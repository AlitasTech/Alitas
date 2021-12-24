package alitas

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
) 

// SocketAlitasProxyClient is the client component of the AlitasProxy that sends RPC requests to Alitas
type SocketAlitasProxyClient struct {
	nodeAddr string
	timeout  time.Duration
	rpc      *rpc.Client
}

// NewSocketAlitasProxyClient implements a new SocketAlitasProxyClient
func NewSocketAlitasProxyClient(nodeAddr string, timeout time.Duration) *SocketAlitasProxyClient {
	return &SocketAlitasProxyClient{
		nodeAddr: nodeAddr,
		timeout:  timeout,
	}
}

func (p *SocketAlitasProxyClient) getConnection() error {
	if p.rpc == nil {
		conn, err := net.DialTimeout("tcp", p.nodeAddr, p.timeout)

		if err != nil {
			return err
		}

		p.rpc = jsonrpc.NewClient(conn)
	}

	return nil
}

// SubmitTx submits a transaction to Alitas
func (p *SocketAlitasProxyClient) SubmitTx(tx []byte) (*bool, error) {
	if err := p.getConnection(); err != nil {
		return nil, err
	}

	var ack bool

	err := p.rpc.Call("Alitas.SubmitTx", tx, &ack)

	if err != nil {
		p.rpc = nil

		return nil, err
	}

	return &ack, nil
}
