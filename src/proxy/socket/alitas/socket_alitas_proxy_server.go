package alitas

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

// SocketAlitasProxyServer is the server component of the AlitasProxy which
// responds to RPC requests from the client component of the AppProxy
type SocketAlitasProxyServer struct {
	netListener *net.Listener
	rpcServer   *rpc.Server
	handler     proxy.ProxyHandler
	timeout     time.Duration
	logger      *logrus.Entry
}

// NewSocketAlitasProxyServer creates a new SocketAlitasProxyServer
func NewSocketAlitasProxyServer(
	bindAddress string,
	handler proxy.ProxyHandler,
	timeout time.Duration,
	logger *logrus.Entry,
) (*SocketAlitasProxyServer, error) {

	server := &SocketAlitasProxyServer{
		handler: handler,
		timeout: timeout,
		logger:  logger,
	}

	if err := server.register(bindAddress); err != nil {
		return nil, err
	}

	return server, nil
}

func (p *SocketAlitasProxyServer) register(bindAddress string) error {
	rpcServer := rpc.NewServer()
	rpcServer.RegisterName("State", p)

	p.rpcServer = rpcServer

	l, err := net.Listen("tcp", bindAddress)

	if err != nil {
		return err
	}

	p.netListener = &l

	return nil
}

func (p *SocketAlitasProxyServer) listen() error {
	for {
		conn, err := (*p.netListener).Accept()

		if err != nil {
			return err
		}

		go (*p.rpcServer).ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

// CommitBlock implements the AppProxy interface
func (p *SocketAlitasProxyServer) CommitBlock(block hashgraph.Block, response *proxy.CommitResponse) (err error) {
	*response, err = p.handler.CommitHandler(block)

	p.logger.WithFields(logrus.Fields{
		"block":    block.Index(),
		"response": response,
		"err":      err,
	}).Debug("AlitasProxyServer.CommitBlock")

	return
}

// GetSnapshot implements the AppProxy interface
func (p *SocketAlitasProxyServer) GetSnapshot(blockIndex int, snapshot *[]byte) (err error) {
	*snapshot, err = p.handler.SnapshotHandler(blockIndex)

	p.logger.WithFields(logrus.Fields{
		"block":    blockIndex,
		"snapshot": snapshot,
		"err":      err,
	}).Debug("AlitasProxyServer.GetSnapshot")

	return
}

// Restore implements the AppProxy interface
func (p *SocketAlitasProxyServer) Restore(snapshot []byte, stateHash *[]byte) (err error) {
	*stateHash, err = p.handler.RestoreHandler(snapshot)

	p.logger.WithFields(logrus.Fields{
		"state_hash": stateHash,
		"err":        err,
	}).Debug("AlitasProxyServer.Restore")

	return
}

// OnStateChanged implements the AppProxy interface
func (p *SocketAlitasProxyServer) OnStateChanged(state state.State, obj *struct{}) (err error) {
	err = p.handler.StateChangeHandler(state)

	p.logger.WithFields(logrus.Fields{
		"state": state.String(),
		"err":   err,
	}).Debug("AlitasProxyServer.OnStateChanged")

	return
}
