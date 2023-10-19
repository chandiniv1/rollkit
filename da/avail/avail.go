package avail

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	ds "github.com/ipfs/go-datastore"
	"github.com/rollkit/rollkit/da"
	"github.com/rollkit/rollkit/log"
	"github.com/rollkit/rollkit/types"
)

type Config struct {
	BaseURL string `json:"base_url"`
	//Seed       string  `json:"seed"`
	// ApiURL     string  `json:"api_url"`
	// AppID      int     `json:"app_id"`
	Confidence float64 `json:"confidence"`
}

type DataAvailabilityLayerClient struct {
	namespace types.NamespaceID
	config    Config
	logger    log.Logger
}

type Confidence struct {
	Block                uint32  `json:"block"`
	Confidence           float64 `json:"confidence"`
	SerialisedConfidence *string `json:"serialised_confidence,omitempty"`
}

type SubmitRequest struct {
	Data string `json:"data"`
}

type SubmitResponse struct {
	BlockHash        string `json:"block_hash"`
	TransactionHash  string `json:"hash"`
	TransactionIndex uint32 `json:"index"`
}

type AppData struct {
	Block      uint32   `json:"block"`
	Extrinsics []string `json:"extrinsics"`
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

	c.logger.Info("starting avail Data Availability Layer Client", "baseURL", c.config.BaseURL)

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
	encodedBlock := base64.StdEncoding.EncodeToString(data)

	requestData := SubmitRequest{
		Data: encodedBlock,
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return da.ResultSubmitBlock{
			BaseResult: da.BaseResult{
				Code:    da.StatusError,
				Message: err.Error(),
			},
		}
	}
	// Make a POST request to the /v2/submit endpoint.
	response, err := http.Post(c.config.BaseURL+"/v2/submit", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return da.ResultSubmitBlock{
			BaseResult: da.BaseResult{
				Code:    da.StatusError,
				Message: err.Error(),
			},
		}
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return da.ResultSubmitBlock{
			BaseResult: da.BaseResult{
				Code:    da.StatusError,
				Message: err.Error(),
			},
		}
	}

	var submitResponse SubmitResponse
	err = json.Unmarshal(responseData, &submitResponse)
	if err != nil {
		return da.ResultSubmitBlock{
			BaseResult: da.BaseResult{
				Code:    da.StatusError,
				Message: err.Error(),
			},
		}
	}
	return da.ResultSubmitBlock{
		BaseResult: da.BaseResult{
			Code:     da.StatusSuccess,
			Message:  "tx hash: " + submitResponse.TransactionHash,
			DAHeight: uint64(submitResponse.TransactionIndex),
		},
	}

}

// data, err := block.MarshalBinary()
// if err != nil {
// 	return da.ResultSubmitBlock{
// 		BaseResult: da.BaseResult{
// 			Code:    da.StatusError,
// 			Message: err.Error(),
// 		},
// 	}
// }

// txHash, err := datasubmit.SubmitData(c.config.ApiURL, c.config.Seed, c.config.AppID, data)

// if err != nil {
// 	return da.ResultSubmitBlock{
// 		BaseResult: da.BaseResult{
// 			Code:    da.StatusError,
// 			Message: err.Error(),
// 		},
// 	}
// }

// return da.ResultSubmitBlock{
// 	BaseResult: da.BaseResult{
// 		Code:     da.StatusSuccess,
// 		Message:  "tx hash: " + hex.EncodeToString(txHash[:]),
// 		DAHeight: 1,
// 	},
// }

// CheckBlockAvailability queries DA layer to check data availability of block.
func (c *DataAvailabilityLayerClient) CheckBlockAvailability(ctx context.Context, dataLayerHeight uint64) da.ResultCheckBlock {

	blockNumber := dataLayerHeight
	confidenceURL := fmt.Sprintf(c.config.BaseURL+"/v1/confidence/%d", blockNumber)

	response, err := http.Get(confidenceURL)

	if err != nil {
		return da.ResultCheckBlock{
			BaseResult: da.BaseResult{
				Code:    da.StatusError,
				Message: err.Error(),
			},
		}
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return da.ResultCheckBlock{
			BaseResult: da.BaseResult{
				Code:    da.StatusError,
				Message: err.Error(),
			},
		}
	}

	var confidenceObject Confidence
	err = json.Unmarshal(responseData, &confidenceObject)
	if err != nil {
		return da.ResultCheckBlock{
			BaseResult: da.BaseResult{
				Code:    da.StatusError,
				Message: err.Error(),
			},
		}
	}

	return da.ResultCheckBlock{
		BaseResult: da.BaseResult{
			Code:     da.StatusSuccess,
			DAHeight: uint64(confidenceObject.Block),
		},
		DataAvailable: confidenceObject.Confidence > float64(c.config.Confidence),
	}
}

//RetrieveBlocks gets the block from DA layer.

func (c *DataAvailabilityLayerClient) RetrieveBlocks(ctx context.Context, dataLayerHeight uint64) da.ResultRetrieveBlocks {
	blocks := []*types.Block{}

Loop:
	blockNumber := dataLayerHeight
	appDataURL := fmt.Sprintf(c.config.BaseURL+"/v1/appdata/%d?decode=true", blockNumber)
	response, err := http.Get(appDataURL)
	if err != nil {
		return da.ResultRetrieveBlocks{
			BaseResult: da.BaseResult{
				Code:    da.StatusError,
				Message: err.Error(),
			},
		}
	}
	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return da.ResultRetrieveBlocks{
			BaseResult: da.BaseResult{
				Code:    da.StatusError,
				Message: err.Error(),
			},
		}
	}

	var appDataObject AppData

	if string(responseData) == "\"Not found\"" {
		appDataObject = AppData{Block: uint32(blockNumber), Extrinsics: []string{}}
	} else if string(responseData) == "\"Processing block\"" {
		goto Loop
	} else {
		err := json.Unmarshal(responseData, &appDataObject)
		if err != nil {
			fmt.Println(string(responseData))
			return da.ResultRetrieveBlocks{
				BaseResult: da.BaseResult{
					Code:    da.StatusError,
					Message: err.Error(),
				},
			}
		}
	}

	txnsByteArray := []byte{}
	for _, extrinsic := range appDataObject.Extrinsics {
		txnsByteArray = append(txnsByteArray, []byte(extrinsic)...)
	}

	block := &types.Block{
		SignedHeader: types.SignedHeader{
			Header: types.Header{
				BaseHeader: types.BaseHeader{
					Height: blockNumber,
				},
			}},
		Data: types.Data{
			Txs: types.Txs{txnsByteArray},
		},
	}
	blocks = append(blocks, block)

	return da.ResultRetrieveBlocks{
		BaseResult: da.BaseResult{
			Code:     da.StatusSuccess,
			DAHeight: uint64(appDataObject.Block),
		},
		Blocks: blocks,
	}
}
