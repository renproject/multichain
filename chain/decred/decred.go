package decred

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrjson/v3"
	"github.com/decred/dcrd/dcrutil/v3"
	"github.com/decred/dcrd/rpc/jsonrpc/types/v2"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/utxo"
	"github.com/renproject/pack"
)

const (
	// DefaultClientTimeout used by the Client.
	DefaultClientTimeout = time.Minute
	// DefaultClientTimeoutRetry used by the Client.
	DefaultClientTimeoutRetry = time.Second
	// DefaultClientHost used by the Client. This should only be used for local
	// deployments of the multichain.
	DefaultClientHost = "https://127.0.0.1:19556"
	// DefaultWalletHost used by dcrwallet. This should only be used for local
	// deployments of the multichain.
	DefaultWalletHost = "https://127.0.0.1:19557"
	// DefaultClientUser used by the Client. This is insecure, and should only
	// be used for local — or publicly accessible — deployments of the
	// multichain.
	DefaultClientUser = "user"
	// DefaultClientPassword used by the Client. This is insecure, and should
	// only be used for local — or publicly accessible — deployments of the
	// multichain.
	DefaultClientPassword = "password"
	DefaultClientNoTLS    = false
	// Authorization types.
	DefaultClientAuthTypeBasic      = "basic"
	DefaultClientAuthTypeClientCert = "clientcert"
	DefaultClientTLSSkipVerify      = false
	DefaultClientCert               = "rpc.cert"
)

// ClientOptions are used to parameterise the behaviour of the Client.
type ClientOptions struct {
	Timeout       time.Duration
	TimeoutRetry  time.Duration
	NoTLS         bool
	Host          string
	WalletHost    string
	User          string
	Password      string
	TLSSkipVerify bool
	AuthType      string
	ClientCert    string
	ClientKey     string
	RPCCert       string
	WalletRPCCert string
}

type ClientSetting struct {
	user       string
	password   string
	host       string
	httpClient http.Client
}

// DefaultClientOptions returns ClientOptions with the default settings. These
// settings are valid for use with the default local deployment of the
// multichain. In production, the host, user, and password should be changed.
func DefaultClientOptions() ClientOptions {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	defaultCertFile := dir + "/../../infra/decred/" + DefaultClientCert

	return ClientOptions{
		Timeout:       DefaultClientTimeout,
		TimeoutRetry:  DefaultClientTimeoutRetry,
		NoTLS:         DefaultClientNoTLS,
		Host:          DefaultClientHost,
		WalletHost:    DefaultWalletHost,
		User:          DefaultClientUser,
		Password:      DefaultClientPassword,
		TLSSkipVerify: DefaultClientTLSSkipVerify,
		AuthType:      DefaultClientAuthTypeBasic,
		RPCCert:       defaultCertFile,
		WalletRPCCert: defaultCertFile,
	}
}

// WithHost sets the URL of the dcrd node.
func (opts ClientOptions) WithHost(host string) ClientOptions {
	opts.Host = host
	return opts
}

// WithRPCCert sets the path of the dcrd RPC cert.
func (opts ClientOptions) WithRPCCert(certPath string) ClientOptions {
	opts.RPCCert = certPath
	return opts
}

// WithWalletHost sets the URL of the dcrwallet node.
func (opts ClientOptions) WithWalletHost(host string) ClientOptions {
	opts.WalletHost = host
	return opts
}

// WithHost sets the path of the dcrwallet RPC cert.
func (opts ClientOptions) WithWalletRPCCert(certPath string) ClientOptions {
	opts.WalletRPCCert = certPath
	return opts
}

// WithUser sets the username that will be used to authenticate with the dcrd
// node.
func (opts ClientOptions) WithUser(user string) ClientOptions {
	opts.User = user
	return opts
}

// WithPassword sets the password that will be used to authenticate with the
// dcrd node.
func (opts ClientOptions) WithPassword(password string) ClientOptions {
	opts.Password = password
	return opts
}

type client struct {
	opts         ClientOptions
	httpClient   http.Client
	walletClient http.Client
}

// NewClient returns a new Client.
func NewClient(opts ClientOptions) *client {
	httpClient := http.Client{}
	httpClient.Timeout = opts.Timeout

	// Configure TLS if needed.
	var tlsConfig *tls.Config
	if !opts.NoTLS {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: opts.TLSSkipVerify,
		}
		if !opts.TLSSkipVerify && opts.AuthType == DefaultClientAuthTypeClientCert {
			serverCAs := x509.NewCertPool()
			serverCert, err := ioutil.ReadFile(opts.RPCCert)
			if err != nil {
				fmt.Printf("Error: %s \n", err)
				return nil
			}
			if !serverCAs.AppendCertsFromPEM(serverCert) {
				fmt.Println("Cannot Append Sever cert")
				return nil
			}
			keypair, err := tls.LoadX509KeyPair(opts.ClientCert, opts.ClientKey)
			if err != nil {
				fmt.Printf("Error: %s \n", err)
				return nil
			}

			tlsConfig.Certificates = []tls.Certificate{keypair}
			tlsConfig.RootCAs = serverCAs

		}
		if !opts.TLSSkipVerify && opts.RPCCert != "" {
			pem, err := ioutil.ReadFile(opts.RPCCert)
			if err != nil {
				fmt.Printf("Error: %s \n", err)
				return nil
			}

			pool := x509.NewCertPool()
			if ok := pool.AppendCertsFromPEM(pem); !ok {
				return nil
			}
			tlsConfig.RootCAs = pool
		}
	}

	httpClient.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// Wallet Client.
	walletClient := http.Client{}
	walletClient.Timeout = opts.Timeout

	// Configure TLS if needed.
	var tlsConf *tls.Config
	if !opts.NoTLS {
		tlsConf = &tls.Config{
			InsecureSkipVerify: opts.TLSSkipVerify,
		}
		if !opts.TLSSkipVerify && opts.AuthType == DefaultClientAuthTypeClientCert {
			serverCAs := x509.NewCertPool()
			serverCert, err := ioutil.ReadFile(opts.WalletRPCCert)
			if err != nil {
				fmt.Printf("Error: %s \n", err)
				return nil
			}
			if !serverCAs.AppendCertsFromPEM(serverCert) {
				return nil
			}
			keypair, err := tls.LoadX509KeyPair(opts.ClientCert, opts.ClientKey)
			if err != nil {
				fmt.Printf("Error: %s \n", err)
				return nil
			}

			tlsConf.Certificates = []tls.Certificate{keypair}
			tlsConf.RootCAs = serverCAs

		}
		if !opts.TLSSkipVerify && opts.WalletRPCCert != "" {
			pem, err := ioutil.ReadFile(opts.WalletRPCCert)
			if err != nil {
				fmt.Printf("Error: %s \n", err)
				return nil
			}

			pool := x509.NewCertPool()
			if ok := pool.AppendCertsFromPEM(pem); !ok {
				return nil
			}
			tlsConf.RootCAs = pool
		}
	}

	walletClient.Transport = &http.Transport{
		TLSClientConfig: tlsConf,
	}

	return &client{
		opts:         opts,
		httpClient:   httpClient,
		walletClient: walletClient,
	}
}

// LatestBlock returns the height of the longest blockchain.
func (client *client) LatestBlock(ctx context.Context) (pack.U64, error) {
	//var resp int64
	var resp dcrjson.Response
	if err := client.send(ctx, &resp, "getbestblock"); err != nil {
		return pack.NewU64(0), fmt.Errorf("get block count: %v", err)
	}

	result := struct {
		Hash   string `json:"hash"`
		Height uint64 `json:"height"`
	}{}

	err := json.Unmarshal(resp.Result, &result)
	if err != nil {
		return pack.NewU64(0), err
	}

	if result.Height < 0 {
		return pack.NewU64(0), fmt.Errorf("unexpected block count, expected > 0, got: %v", result.Height)
	}

	return pack.NewU64(uint64(result.Height)), nil
}

// UnspentOutputs spendable by the given address.
func (client *client) UnspentOutputs(ctx context.Context, minConf, maxConf int64, addr address.Address) ([]utxo.Output, error) {
	var outputs []utxo.Output

	//var resp int64
	var resp dcrjson.Response
	if err := client.send(ctx, &resp, "listunspent", minConf, maxConf, []string{string(addr)}); err != nil {
		return []utxo.Output{}, fmt.Errorf("bad \"listunspent\": %v", err)
	}

	//outputs := make([]utxo.Output, len(resp.Result))
	type Result struct {
		TxId          string  `json:"txid"`
		VOut          uint32  `json:"vout"`
		Tree          int     `json:"tree"`
		TxType        int     `json:"txtype"`
		Address       string  `json:"address"`
		Account       string  `json:"account"`
		ScriptPubKey  string  `json:"scriptPubKey"`
		Amount        float64 `json:"amount"`
		Confirmations int64   `json:"confirmations"`
		Spendable     bool    `json"spendable"`
	}

	var result []Result

	err := json.Unmarshal(resp.Result, &result)
	if err != nil {
		return []utxo.Output{}, fmt.Errorf("bad \"listunspent\": %v", err)
	}

	//fmt.Printf("Unspent Output: %+v \n", resp)
	for _, v := range result {
		amount, err := dcrutil.NewAmount(v.Amount)
		if err != nil {
			return []utxo.Output{}, fmt.Errorf("bad amount: %v", err)
		}
		if amount < 0 {
			return []utxo.Output{}, fmt.Errorf("bad amount: %v", amount)
		}
		pubKeyScript, err := hex.DecodeString(v.ScriptPubKey)
		if err != nil {
			return []utxo.Output{}, fmt.Errorf("bad pubkey script: %v", err)
		}
		txid, err := chainhash.NewHashFromStr(v.TxId)
		if err != nil {
			return []utxo.Output{}, fmt.Errorf("bad txid: %v", err)
		}
		o := utxo.Output{
			Outpoint: utxo.Outpoint{
				Hash:  pack.NewBytes(txid[:]),
				Index: pack.NewU32(v.VOut),
			},
			Value:        pack.NewU256FromU64(pack.NewU64(uint64(amount))),
			PubKeyScript: pack.NewBytes(pubKeyScript),
		}

		outputs = append(outputs, o)
	}

	return outputs, nil
}

// Output associated with an outpoint, and its number of confirmations.
func (client *client) Output(ctx context.Context, outpoint utxo.Outpoint) (utxo.Output, pack.U64, error) {

	//var resp int64
	var resp dcrjson.Response
	hash := chainhash.Hash{}
	copy(hash[:], outpoint.Hash)
	if err := client.send(ctx, &resp, "getrawtransaction", hash.String(), 1); err != nil {
		return utxo.Output{}, pack.NewU64(0), fmt.Errorf("bad \"getrawtransaction\": %v", err)
	}

	result := types.TxRawResult{}
	err := json.Unmarshal(resp.Result, &result)
	if err != nil {
		return utxo.Output{}, pack.NewU64(0), fmt.Errorf("bad \"getrawtransaction\": %v", err)
	}

	if outpoint.Index.Uint32() >= uint32(len(result.Vout)) {
		return utxo.Output{}, pack.NewU64(0), fmt.Errorf("bad index: %v is out of range", outpoint.Index)
	}
	vout := result.Vout[outpoint.Index.Uint32()]
	amount, err := dcrutil.NewAmount(vout.Value)
	if err != nil {
		return utxo.Output{}, pack.NewU64(0), fmt.Errorf("bad amount: %v", err)
	}
	if amount < 0 {
		return utxo.Output{}, pack.NewU64(0), fmt.Errorf("bad amount: %v", amount)
	}
	pubKeyScript, err := hex.DecodeString(vout.ScriptPubKey.Hex)
	if err != nil {
		return utxo.Output{}, pack.NewU64(0), fmt.Errorf("bad pubkey script: %v", err)
	}
	output := utxo.Output{
		Outpoint:     outpoint,
		Value:        pack.NewU256FromU64(pack.NewU64(uint64(amount))),
		PubKeyScript: pack.NewBytes(pubKeyScript),
	}
	return output, pack.NewU64(uint64(result.Confirmations)), nil
}

// UnspentOutput returns the unspent transaction output identified by the
// given outpoint. It also returns the number of confirmations for the
// output. If the output cannot be found before the context is done, the
// output is invalid, or the output has been spent, then an error should be
// returned.
func (client *client) UnspentOutput(ctx context.Context, outpoint utxo.Outpoint) (utxo.Output, pack.U64, error) {
	var resp dcrjson.Response
	hash := chainhash.Hash{}
	copy(hash[:], outpoint.Hash)
	if err := client.send(ctx, &resp, "gettxout", hash.String(), outpoint.Index.Uint32()); err != nil {
		return utxo.Output{}, pack.NewU64(0), fmt.Errorf("bad \"gettxout\": %v", err)
	}

	result := types.GetTxOutResult{}
	err := json.Unmarshal(resp.Result, &result)
	if err != nil {
		return utxo.Output{}, pack.NewU64(0), fmt.Errorf("bad \"gettxout\": %v", err)
	}

	amount, err := dcrutil.NewAmount(result.Value)
	if err != nil {
		return utxo.Output{}, pack.NewU64(0), fmt.Errorf("bad amount: %v", err)
	}
	if amount < 0 {
		return utxo.Output{}, pack.NewU64(0), fmt.Errorf("bad amount: %v", amount)
	}
	if result.Confirmations < 0 {
		return utxo.Output{}, pack.NewU64(0), fmt.Errorf("bad confirmations: %v", result.Confirmations)
	}
	pubKeyScript, err := hex.DecodeString(result.ScriptPubKey.Hex)
	if err != nil {
		return utxo.Output{}, pack.NewU64(0), fmt.Errorf("bad pubkey script: %v", err)
	}
	output := utxo.Output{
		Outpoint:     outpoint,
		Value:        pack.NewU256FromU64(pack.NewU64(uint64(amount))),
		PubKeyScript: pack.NewBytes(pubKeyScript),
	}
	return output, pack.NewU64(uint64(result.Confirmations)), nil
}

func (client *client) send(ctx context.Context, resp *dcrjson.Response, method string, params ...interface{}) error {
	// Encode the request.
	data, err := encodeRequest(method, params)
	if err != nil {
		return err
	}

	var clSetting *ClientSetting
	switch method {
	case "getbestblock", "getrawtransaction", "gettxout",
		"sendrawtransaction", "estimatesmartfee", "estimatefee":
		clSetting = &ClientSetting{
			user:       client.opts.User,
			password:   client.opts.Password,
			host:       client.opts.Host,
			httpClient: client.httpClient,
		}
	case "listunspent", "gettransaction":
		clSetting = &ClientSetting{
			user:       client.opts.User,
			password:   client.opts.Password,
			host:       client.opts.WalletHost,
			httpClient: client.walletClient,
		}
	}

	return retry(ctx, client.opts.TimeoutRetry, func() error {
		// Create request and add basic authentication headers. The context is
		// not attached to the request, and instead we all each attempt to run
		// for the timeout duration, and we keep attempting until success, or
		// the context is done.
		req, err := http.NewRequest("POST", clSetting.host, bytes.NewBuffer(data))
		if err != nil {
			return fmt.Errorf("building http request: %v", err)
		}
		req.SetBasicAuth(clSetting.user, clSetting.password)

		// Send the request and decode the response.
		res, err := clSetting.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("sending http request: %v", err)
		}
		// Read the raw bytes and close the response.
		respBytes, err := ioutil.ReadAll(res.Body)
		defer res.Body.Close()
		//if err := decodeResponse(resp, res.Body); err != nil {
		//	return fmt.Errorf("decoding http response: %v", err)
		//}

		if err != nil {
			err = fmt.Errorf("error reading json reply: %s", err)
			return err
		}

		// Handle unsuccessful HTTP responses
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			// Generate a standard error to return if the server body is
			// empty.  This should not happen very often, but it's better
			// than showing nothing in case the target server has a poor
			// implementation.
			if len(respBytes) == 0 {
				return fmt.Errorf("%d %s", res.StatusCode,
					http.StatusText(res.StatusCode))
			}
			return fmt.Errorf("%s", respBytes)
		}

		// Unmarshal the response.
		// var resp dcrjson.Response
		if err := json.Unmarshal(respBytes, resp); err != nil {
			return err
		}

		return nil
	})
}

func encodeRequest(method string, params []interface{}) ([]byte, error) {
	rawParams, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("encoding params: %v", err)
	}
	req := struct {
		Version string          `json:"jsonrpc"`
		ID      int             `json:"id"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params"`
	}{
		Version: "1.0",
		ID:      1,
		Method:  method,
		Params:  rawParams,
	}
	rawReq, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encoding request: %v", err)
	}
	return rawReq, nil
}

func decodeResponse(resp interface{}, r io.Reader) error {
	res := struct {
		Version string           `json:"version"`
		ID      int              `json:"id"`
		Result  *json.RawMessage `json:"result"`
		Error   *json.RawMessage `json:"error"`
	}{}
	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return fmt.Errorf("decoding response: %v", err)
	}
	if res.Error != nil {
		return fmt.Errorf("decoding response: %v", string(*res.Error))
	}
	if res.Result == nil {
		return fmt.Errorf("decoding result: result is nil")
	}
	if err := json.Unmarshal(*res.Result, resp); err != nil {
		return fmt.Errorf("decoding result: %v", err)
	}
	return nil
}

func retry(ctx context.Context, dur time.Duration, f func() error) error {
	ticker := time.NewTicker(dur)
	err := f()
	for err != nil {
		log.Printf("retrying: %v", err)
		select {
		case <-ctx.Done():
			return fmt.Errorf("%v: %v", ctx.Err(), err)
		case <-ticker.C:
			err = f()
		}
	}
	return nil
}

// SubmitTx to the Decred network.
func (client *client) SubmitTx(ctx context.Context, tx utxo.Tx) error {
	serial, err := tx.Serialize()
	if err != nil {
		return fmt.Errorf("bad tx: %v", err)
	}
	var resp dcrjson.Response
	if err := client.send(ctx, &resp, "sendrawtransaction", hex.EncodeToString(serial), true); err != nil {
		return fmt.Errorf("bad \"sendrawtransaction\": %v", err)
	}
	fmt.Printf("Response: %+v \n", string(resp.Result))
	return nil
}

func (client *client) EstimateSmartFee(ctx context.Context, numBlocks int64) (float64, error) {
	resp := dcrjson.Response{}

	if err := client.send(ctx, &resp, "estimatesmartfee", numBlocks); err != nil {
		return 0.0, fmt.Errorf("estimating smart fee: %v", err)
	}

	if resp.Error != nil {
		return 0.0, resp.Error
	}

	result, err := strconv.ParseFloat(string(resp.Result), 10)
	if err != nil {
		return 0.0, err
	}

	return result, nil
}

func (client *client) EstimateFeeLegacy(ctx context.Context, numBlocks int64) (float64, error) {
	resp := dcrjson.Response{}

	if err := client.send(ctx, &resp, "estimatefee", numBlocks); err != nil {
		return 0.0, fmt.Errorf("estimating fee: %v", err)
	}

	if resp.Error != nil {
		return 0.0, resp.Error
	}

	result, err := strconv.ParseFloat(string(resp.Result), 10)
	if err != nil {
		return 0.0, err
	}

	return result, nil
}

// Confirmations of a transaction in the decred network.
func (client *client) Confirmations(ctx context.Context, txHash pack.Bytes) (int64, error) {
	var resp dcrjson.Response

	size := len(txHash)
	txHashReversed := make([]byte, size)
	copy(txHashReversed[:], txHash[:])
	for i := 0; i < size/2; i++ {
		txHashReversed[i], txHashReversed[size-1-i] = txHashReversed[size-1-i], txHashReversed[i]
	}

	if err := client.send(ctx, &resp, "gettransaction", hex.EncodeToString(txHashReversed)); err != nil {
		return 0, fmt.Errorf("bad \"gettransaction\": %v", err)
	}

	result := types.TxRawResult{}
	err := json.Unmarshal(resp.Result, &result)
	if err != nil {
		return 0, fmt.Errorf("bad \"gettransaction\": %v", err)
	}

	confirmations := result.Confirmations
	if confirmations < 0 {
		confirmations = 0
	}
	return confirmations, nil
}
