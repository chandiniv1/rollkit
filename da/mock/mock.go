package mock

import (
	"encoding/binary"

	"github.com/celestiaorg/optimint/da"
	"github.com/celestiaorg/optimint/log"
	"github.com/celestiaorg/optimint/store"
	"github.com/celestiaorg/optimint/types"
)

// MockDataAvailabilityLayerClient is intended only for usage in tests.
// It does actually ensures DA - it stores data in-memory.
type MockDataAvailabilityLayerClient struct {
	logger log.Logger

	Blocks     map[[32]byte]*types.Block
	BlockIndex map[uint64][32]byte

	dalcKV store.KVStore
}

var _ da.DataAvailabilityLayerClient = &MockDataAvailabilityLayerClient{}
var _ da.BlockRetriever = &MockDataAvailabilityLayerClient{}

// Init is called once to allow DA client to read configuration and initialize resources.
func (m *MockDataAvailabilityLayerClient) Init(config []byte, dalcKV store.KVStore, logger log.Logger) error {
	m.logger = logger
	m.Blocks = make(map[[32]byte]*types.Block)
	m.BlockIndex = make(map[uint64][32]byte)
	m.dalcKV = dalcKV
	return nil
}

// Start implements DataAvailabilityLayerClient interface.
func (m *MockDataAvailabilityLayerClient) Start() error {
	m.logger.Debug("Mock Data Availability Layer Client starting")
	return nil
}

// Stop implements DataAvailabilityLayerClient interface.
func (m *MockDataAvailabilityLayerClient) Stop() error {
	m.logger.Debug("Mock Data Availability Layer Client stopped")
	return nil
}

// SubmitBlock submits the passed in block to the DA layer.
// This should create a transaction which (potentially)
// triggers a state transition in the DA layer.
func (m *MockDataAvailabilityLayerClient) SubmitBlock(block *types.Block) da.ResultSubmitBlock {
	m.logger.Debug("Submitting block to DA layer!", "height", block.Header.Height)

	hash := block.Header.Hash()
	blob, err := block.MarshalBinary()
	if err != nil {
		return da.ResultSubmitBlock{
			DAResult: da.DAResult{
				Code:    da.StatusError,
				Message: err.Error(),
			},
		}
	}

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, block.Header.Height)

	m.dalcKV.Set(b, hash[:])
	m.dalcKV.Set(hash[:], blob)

	m.Blocks[hash] = block
	m.BlockIndex[block.Header.Height] = hash

	return da.ResultSubmitBlock{
		DAResult: da.DAResult{
			Code:    da.StatusSuccess,
			Message: "OK",
		},
	}
}

// CheckBlockAvailability queries DA layer to check data availability of block corresponding to given header.
func (m *MockDataAvailabilityLayerClient) CheckBlockAvailability(header *types.Header) da.ResultCheckBlock {
	_, ok := m.Blocks[header.Hash()]
	return da.ResultCheckBlock{DAResult: da.DAResult{Code: da.StatusSuccess}, DataAvailable: ok}
}

// RetrieveBlock returns block at given height from data availability layer.
func (m *MockDataAvailabilityLayerClient) RetrieveBlock(height uint64) da.ResultRetrieveBlock {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, height)
	hash, err := m.dalcKV.Get(b)
	if err != nil {
		return da.ResultRetrieveBlock{DAResult: da.DAResult{Code: da.StatusError, Message: err.Error()}}
	}
	blob, err := m.dalcKV.Get(hash)
	if err != nil {
		return da.ResultRetrieveBlock{DAResult: da.DAResult{Code: da.StatusError, Message: err.Error()}}
	}

	block := &types.Block{}
	err = block.UnmarshalBinary(blob)
	if err != nil {
		return da.ResultRetrieveBlock{DAResult: da.DAResult{Code: da.StatusError, Message: err.Error()}}
	}

	return da.ResultRetrieveBlock{DAResult: da.DAResult{Code: da.StatusSuccess}, Block: block}
}