package icon

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/icon-project/goloop/server/jsonrpc"
	"github.com/renproject/pack"
)

// Estimator ...
type Estimator struct {
	Client Client
	Tx     Tx
}

type RP struct {
	Version string         `json:"jsonrpc" validate:"required"`
	ID      interface{}    `json:"id"`
	Result  jsonrpc.HexInt `json:"result"`
}

type getMaxStepLimitParams struct {
	To       jsonrpc.Address `json:"to" validate:"required,t_addr_score"`
	DataType string          `json:"dataType" validate:"required,call"`
	Data     interface{}     `json:"data"`
}

// EstimateGasPrice ...
func (es Estimator) EstimateGasPrice(context.Context) (pack.U256, error) {
	rq := &Request{
		Version: jsonprcVersion,
		Method:  "icx_call",
		ID:      time.Now().UnixNano() / int64(time.Millisecond),
	}

	params := getMaxStepLimitParams{
		To:       "cx0000000000000000000000000000000000000001",
		DataType: "call",
	}

	type paramData struct {
		ContextType string `json:"contextType"`
	}

	type dataParams struct {
		Method string    `json:"method"`
		Params paramData `json:"params"`
	}

	pData := paramData{
		ContextType: "invoke",
	}

	dParams := dataParams{
		Method: "getMaxStepLimit",
		Params: pData,
	}

	params.Data = dParams

	mParams, _ := json.Marshal(params)

	rq.Params = json.RawMessage(mParams)

	buf, err := json.Marshal(rq)
	if err != nil {
		log.Fatal(err)
	}

	result, err := http.Post(es.Client.v3.Endpoint, "application/json; charset=utf-8", bytes.NewReader(buf))
	if err != nil {
		return pack.U256{}, err
	}
	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Panic(err)
		return pack.U256{}, err
	}

	var resp RP
	json.Unmarshal(body, &resp)

	return pack.NewU256FromU64(pack.NewU64(uint64(resp.Result.Value()))), nil
}
