package app

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"github.com/AlitasTech/Alitas/src/hashgraph"
	"github.com/AlitasTech/Alitas/src/node/state"
	"github.com/AlitasTech/Alitas/src/proxy"
	"github.com/sirupsen/logrus"
)

// SocketAppProxyClient is the component of the AppGateway that sends RPC requests to the App
type SocketAppProxyClient struct {
	clientAddr string
	timeout    time.Duration
	logger     *logrus.Entry
	rpc        *rpc.Client
}

// NewSocketAppProxyClient creates a new SocketAppProxyClient
func NewSocketAppProxyClient(clientAddr string, timeout time.Duration, logger *logrus.Entry) *SocketAppProxyClient {
	return &SocketAppProxyClient{
		clientAddr: clientAddr,
		timeout:    timeout,
		logger:     logger,
	}
}

func (p *SocketAppProxyClient) getConnection() error {
	if p.rpc == nil {
		conn, err := net.DialTimeout("tcp", p.clientAddr, p.timeout)

		if err != nil {
			return err
		}

		p.rpc = jsonrpc.NewClient(conn)
	}

	return nil
}

// CommitBlock implements the AppGateway interface
func (p *SocketAppProxyClient) CommitBlock(block hashgraph.Block) (proxy.CommitResponse, error) {
	if err := p.getConnection(); err != nil {
		return proxy.CommitResponse{}, err
	}

	var commitResponse proxy.CommitResponse

	if err := p.rpc.Call("State.CommitBlock", block, &commitResponse); err != nil {
		p.rpc = nil

		return commitResponse, err
	}

	p.logger.WithFields(logrus.Fields{
		"block":           block.Index(),
		"commit_response": commitResponse,
	}).Debug("AppProxyClient.CommitBlock")

	return commitResponse, nil
}

// GetSnapshot implementes the AppGateway interface
func (p *SocketAppProxyClient) GetSnapshot(blockIndex int) ([]byte, error) {
	if err := p.getConnection(); err != nil {
		return []byte{}, err
	}

	var snapshot []byte

	if err := p.rpc.Call("State.GetSnapshot", blockIndex, &snapshot); err != nil {
		p.rpc = nil

		return []byte{}, err
	}

	p.logger.WithFields(logrus.Fields{
		"block":    blockIndex,
		"snapshot": snapshot,
	}).Debug("AppProxyClient.GetSnapshot")

	return snapshot, nil
}

// Restore implements the AppGateway interface
func (p *SocketAppProxyClient) Restore(snapshot []byte) error {
	if err := p.getConnection(); err != nil {
		return err
	}

	var stateHash []byte

	if err := p.rpc.Call("State.Restore", snapshot, &stateHash); err != nil {
		p.rpc = nil

		return err
	}

	p.logger.WithFields(logrus.Fields{
		"state_hash": stateHash,
	}).Debug("AppProxyClient.Restore")

	return nil
}

// OnStateChanged implements the AppGateway interface
func (p *SocketAppProxyClient) OnStateChanged(state state.State) error {
	if err := p.getConnection(); err != nil {
		return err
	}

	if err := p.rpc.Call("State.OnStateChanged", state, nil); err != nil {
		p.rpc = nil

		return err
	}

	p.logger.WithFields(logrus.Fields{
		"state": state.String,
	}).Debug("AppProxyClient.OnStateChanged")

	return nil
}
