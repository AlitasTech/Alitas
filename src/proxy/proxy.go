package proxy

import (
	"github.com/AlitasTech/Alitas/src/hashgraph"
	"github.com/AlitasTech/Alitas/src/node/state"
)

// AppGateway defines the interface which is used by Alitas to communicate with the App
type AppGateway interface {
	SubmitCh() chan []byte
	CommitBlock(block hashgraph.Block) (CommitResponse, error)
	GetSnapshot(blockIndex int) ([]byte, error)
	Restore(snapshot []byte) error
	OnStateChanged(state.State) error
}
