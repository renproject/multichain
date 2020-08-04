package ontology

import (
	ontSdk "github.com/ontio/ontology-go-sdk"
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/smartcontract/service/native/ont"
)

type client struct {
	sdk *ontSdk.OntologySdk
}

// NewClient returns a new Client.
func NewClient(host string) *client {
	c := &client{
		sdk: ontSdk.NewOntologySdk(),
	}
	c.sdk.NewRpcClient().SetAddress(host)
	return c
}

func (client *client) GenAddress() common.Address {
	return ontSdk.NewAccount().Address
}

func (client *client) MultiTransferOnt(gasPrice, gasLimit uint64, payer *ontSdk.Account, states []*ont.State,
	signer *ontSdk.Account) (common.Uint256, error) {
	return client.sdk.Native.Ont.MultiTransfer(gasPrice, gasLimit, payer, states, signer)
}

func (client *client) GetEvent(txHash string) (*sdkcom.SmartContactEvent, error) {
	return client.sdk.GetSmartContractEvent(txHash)
}
