package avail

import (
	"context"
	"encoding/json"

	// "io/ioutil"
	// "net/http"
	// "time"

	// "github.com/rollkit/celestia-openrpc/types/share"

	ds "github.com/ipfs/go-datastore"
	openrpc "github.com/rollkit/celestia-openrpc"
	openrpcns "github.com/rollkit/celestia-openrpc/types/namespace"

	"github.com/rollkit/rollkit/da"
	"github.com/rollkit/rollkit/da/avail/datasubmit"
	"github.com/rollkit/rollkit/log"
	"github.com/rollkit/rollkit/types"
)

type Config struct {
	Seed   string `json:"seed"`
	ApiURL string `json:"api_url"`
	Size   int    `json:"size"`
	AppID  int    `json:"app_id"`
	Dest   string `json:"dest"`
	Amount uint64 `json:amount`
}

// DataAvailabilityLayerClient use celestia-node public API.
type DataAvailabilityLayerClient struct {
	rpc       *openrpc.Client
	namespace openrpcns.Namespace
	config    Config
	logger    log.Logger
}

var _ da.DataAvailabilityLayerClient = &DataAvailabilityLayerClient{}
var _ da.BlockRetriever = &DataAvailabilityLayerClient{}

// Init initializes DataAvailabilityLayerClient instance.
func (c *DataAvailabilityLayerClient) Init(namespaceID types.NamespaceID, config []byte, kvStore ds.Datastore, logger log.Logger) error {
	c.logger = logger

	if len(config) > 0 {
		return json.Unmarshal(config, &c.config)
	}

	return nil
}

// Start prepares DataAvailabilityLayerClient to work.
func (c *DataAvailabilityLayerClient) Start() error {
	c.logger.Info("starting avail Data Availability Layer Client", "baseURL", c.config.ApiURL)

	return nil
}

// Stop stops DataAvailabilityLayerClient.
func (c *DataAvailabilityLayerClient) Stop() error {
	c.logger.Info("stopping Avail Data Availability Layer Client")
	return nil
}

// SubmitBlock submits a block to DA layer.
func (c *DataAvailabilityLayerClient) SubmitBlock(ctx context.Context, block *types.Block) da.ResultSubmitBlock {

	data, err := block.MarshalBinary()
	if err != nil {
		return da.ResultSubmitBlock{
			BaseResult: da.BaseResult{
				Code:    da.StatusError,
				Message: err.Error(),
			},
		}
	}

	txResponseErr := datasubmit.SubmitData(1000, c.config.ApiURL, c.config.Seed, c.config.AppID, data)

	if txResponseErr != nil {
		return da.ResultSubmitBlock{
			BaseResult: da.BaseResult{
				Code:    da.StatusError,
				Message: txResponseErr.Error(),
			},
		}
	}

	return da.ResultSubmitBlock{
		BaseResult: da.BaseResult{
			Code:     da.StatusSuccess,
			Message:  "data submitted succesfully ",
			DAHeight: 1,
		},
	}

}
