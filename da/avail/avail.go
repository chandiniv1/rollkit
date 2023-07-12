package avail

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	//"time"

	ds "github.com/ipfs/go-datastore"
	//openrpc "github.com/rollkit/celestia-openrpc"

	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	"github.com/centrifuge/go-substrate-rpc-client/v4/rpc"
	"github.com/rollkit/rollkit/da"
	"github.com/rollkit/rollkit/da/avail/mock"

	//"github.com/rollkit/rollkit/da/mock"
	"github.com/rollkit/rollkit/log"
	"github.com/rollkit/rollkit/types"
)

// DataAvailabilityLayerClient use celestia-node public API.
type DataAvailabilityLayerClient struct {
	RPC         *rpc.RPC
	Client      client.Client
	namespaceID types.NamespaceID
	config      Config
	logger      log.Logger
}

type Config struct {
	Seed   string `json:"seed"`
	ApiURL string `json:"api_url"`
	Size   int    `json:"size"`
	AppID  int    `json:"app_id"`
	Dest   string `json:"dest"`
	Amount uint64 `json:amount`
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
	// a.logger.Info("starting Avail Data Availability Layer Client", "baseURL", a.config.BaseURL)
	// api, err := gsrpc.NewSubstrateAPI(a.config.BaseURL)
	// if err != nil {
	// 	return fmt.Errorf("cannot create api:%w", err)
	// }
	// fmt.Println("api:", api)
	// return nil
	c.logger.Info("starting avail Data Availability Layer Client", "baseURL", c.config.ApiURL)
	// var err error
	// c.client, err = cnc.NewClient(c.BaseURL, cnc.WithTimeout(c.config.TimeOut))
	return nil
}

// Stop stops DataAvailabilityLayerClient.
func (c *DataAvailabilityLayerClient) Stop() error {
	c.logger.Info("stopping Avail Data Availability Layer Client")
	return nil
}

// SubmitBlock submits a block to DA layer.
func (c *DataAvailabilityLayerClient) SubmitBlock(ctx context.Context, block *types.Block) da.ResultSubmitBlock {

	fmt.Println("submit block method called.......")

	txResponseErr := mock.SubmitData(1000, c.config.ApiURL, c.config.Seed, 0)

	fmt.Println("txResponseErr------------", txResponseErr)

	if txResponseErr != nil {
		return da.ResultSubmitBlock{
			BaseResult: da.BaseResult{
				Code:    da.StatusError,
				Message: txResponseErr.Error(),
			},
		}
	}

	fmt.Println("submit block method2 called.......")

	// blockNumber := 149448
	// confidenceURL := fmt.Sprintf("http://localhost:7000/v1/confidence/%d", blockNumber)

	// response, err := http.Get(confidenceURL)

	// responseData, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	return da.ResultSubmitBlock{
	// 		BaseResult: da.BaseResult{
	// 			Code:    da.StatusError,
	// 			Message: err.Error(),
	// 		},
	// 	}
	// }
	// fmt.Println(string(responseData))

	return da.ResultSubmitBlock{
		BaseResult: da.BaseResult{
			Code:     da.StatusSuccess,
			Message:  "responseData: ", //+ string(responseData),
			DAHeight: uint64(1),
		},
	}
}

// CheckBlockAvailability queries DA layer to check data availability of block at given height.
func (a *DataAvailabilityLayerClient) CheckBlockAvailability(ctx context.Context, dataLayerHeight uint64) da.ResultCheckBlock {

	// type Confidence struct {
	// 	Block                uint32  `json:"block"`
	// 	Confidence           float64 `json:"confidence"`
	// 	SerialisedConfidence *string `json:"serialised_confidence,omitempty"`
	// }

	// fmt.Println("check block availability called.........")
	// var blockNumber int
	// confidenceURL := fmt.Sprintf("http://localhost:7000/v1/confidence/%d", blockNumber)

	// response, err := http.Get(confidenceURL)

	// if err != nil {
	// 	return da.ResultCheckBlock{
	// 		BaseResult: da.BaseResult{
	// 			Code:    da.StatusError,
	// 			Message: err.Error(),
	// 		},
	// 	}
	// }

	// responseData, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	return da.ResultCheckBlock{
	// 		BaseResult: da.BaseResult{
	// 			Code:    da.StatusError,
	// 			Message: err.Error(),
	// 		},
	// 	}
	// }

	// var confidenceObject Confidence
	// json.Unmarshal(responseData, &confidenceObject)

	return da.ResultCheckBlock{
		BaseResult: da.BaseResult{
			Code:     da.StatusSuccess,
			DAHeight: dataLayerHeight,
		},
		DataAvailable: fasle//confidenceObject.Confidence > 92,
	}

}

// RetrieveBlocks gets a batch of blocks from DA layer.
func (c *DataAvailabilityLayerClient) RetrieveBlocks(ctx context.Context, dataLayerHeight uint64) da.ResultRetrieveBlocks {

	blocks := make([]*types.Block, 1)

	blocks[0] = new(types.Block)

	// var blockNumber int
	// confidenceURL := fmt.Sprintf("http://localhost:7000/v1/confidence/%d", blockNumber)

	// response, err := http.Get(confidenceURL)

	// responseData, err := ioutil.ReadAll(response.Body)
	// if err != nil {
	// 	return da.ResultRetrieveBlocks{
	// 		BaseResult: da.BaseResult{
	// 			Code:    da.StatusError,
	// 			Message: err.Error(),
	// 		},
	// 	}
	// }
	// fmt.Println(string(responseData))

	return da.ResultRetrieveBlocks{
		BaseResult: da.BaseResult{
			Code:     da.StatusSuccess,
			DAHeight: dataLayerHeight,
			Message:  "block data: " + string(responseData),
		},
		Blocks: blocks,
	}
}
