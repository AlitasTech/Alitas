package proxy

import (
	"github.com/AlitasTech/Alitas/src/hashgraph"
	"github.com/AlitasTech/Alitas/src/node/state"
)

/*
These types are exported and need to be implemented and used by the calling
application.
*/

// ProxyHandler is used to instanciate an InmemProxy
type ProxyHandler interface {
	// CommitHandler is called when Alitas commits a block to the application
	CommitHandler(block hashgraph.Block) (response CommitResponse, err error)

	// SnapshotHandler is called by Alitas to retrieve a snapshot corresponding
	// to a particular block
	SnapshotHandler(blockIndex int) (snapshot []byte, err error)

	// RestoreHandler is called by Alitas to restore the application to a
	// specific state
	RestoreHandler(snapshot []byte) (stateHash []byte, err error)

	// StateChangeHandler is called by onStateChanged to notify that a Alitas
	// node entered a certain state
	StateChangeHandler(state.State) error
}
