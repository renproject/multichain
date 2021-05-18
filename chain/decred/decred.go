package decred

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/decred/dcrd/dcrjson/v3"
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
	User          string
	Password      string
	TLSSkipVerify bool
	AuthType      string
	ClientCert    string
	ClientKey     string
	RPCCert       string
}

// DefaultClientOptions returns ClientOptions with the default settings. These
// settings are valid for use with the default local deployment of the
// multichain. In production, the host, user, and password should be changed.
func DefaultClientOptions() ClientOptions {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	defaultCertFile := dir + "/decred/" + DefaultClientCert

	return ClientOptions{
		Timeout:       DefaultClientTimeout,
		TimeoutRetry:  DefaultClientTimeoutRetry,
		NoTLS:         DefaultClientNoTLS,
		Host:          DefaultClientHost,
		User:          DefaultClientUser,
		Password:      DefaultClientPassword,
		TLSSkipVerify: DefaultClientTLSSkipVerify,
		AuthType:      DefaultClientAuthTypeBasic,
		RPCCert:       defaultCertFile,
	}
}

// WithHost sets the URL of the Bitcoin node.
func (opts ClientOptions) WithHost(host string) ClientOptions {
	opts.Host = host
	return opts
}

// WithUser sets the username that will be used to authenticate with the Bitcoin
// node.
func (opts ClientOptions) WithUser(user string) ClientOptions {
	opts.User = user
	return opts
}

// WithPassword sets the password that will be used to authenticate with the
// Bitcoin node.
func (opts ClientOptions) WithPassword(password string) ClientOptions {
	opts.Password = password
	return opts
}

type client struct {
	opts       ClientOptions
	httpClient http.Client
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
				return nil
			}
			if !serverCAs.AppendCertsFromPEM(serverCert) {
				return nil
			}
			keypair, err := tls.LoadX509KeyPair(opts.ClientCert, opts.ClientKey)
			if err != nil {
				return nil
			}

			tlsConfig.Certificates = []tls.Certificate{keypair}
			tlsConfig.RootCAs = serverCAs

		}
		if !opts.TLSSkipVerify && opts.RPCCert != "" {
			pem, err := ioutil.ReadFile(opts.RPCCert)
			if err != nil {
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

	return &client{
		opts:       opts,
		httpClient: httpClient,
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

func (client *client) send(ctx context.Context, resp *dcrjson.Response, method string, params ...interface{}) error {
	// Encode the request.
	data, err := encodeRequest(method, params)
	if err != nil {
		return err
	}

	return retry(ctx, client.opts.TimeoutRetry, func() error {
		// Create request and add basic authentication headers. The context is
		// not attached to the request, and instead we all each attempt to run
		// for the timeout duration, and we keep attempting until success, or
		// the context is done.
		req, err := http.NewRequest("POST", client.opts.Host, bytes.NewBuffer(data))
		if err != nil {
			return fmt.Errorf("building http request: %v", err)
		}
		req.SetBasicAuth(client.opts.User, client.opts.Password)

		// Send the request and decode the response.
		res, err := client.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("sending http request: %v", err)
		}
		// Read the raw bytes and close the response.
		respBytes, err := ioutil.ReadAll(res.Body)
		//fmt.Printf("Response: %+v \n", res)
		defer res.Body.Close()
		//if err := decodeResponse(resp, res.Body); err != nil {
		//	return fmt.Errorf("decoding http response: %v", err)
		//}

		if err != nil {
			err = fmt.Errorf("error reading json reply: %w", err)
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
