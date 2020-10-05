package icon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"time"

	"github.com/icon-project/goloop/server/jsonrpc"
)

const (
	headerContentType   = "Content-Type"
	headerAccept        = "Accept"
	typeApplicationJSON = "application/json"
)

type JsonRpcClient struct {
	hc           *http.Client
	Endpoint     string
	CustomHeader map[string]string
	Pre          func(req *http.Request) error
}

type Response struct {
	Version string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *jsonrpc.Error  `json:"error,omitempty"`
	ID      interface{}     `json:"id"`
}

type HttpError struct {
	response string
	message  string
}

func (e *HttpError) Error() string {
	return e.message
}

func (e *HttpError) Response() string {
	return e.response
}

func NewHttpError(r *http.Response) error {
	var response string
	if rb, err := ioutil.ReadAll(r.Body); err != nil {
		response = fmt.Sprintf("Fail to read body err=%+v", err)
	} else {
		response = string(rb)
	}
	return &HttpError{
		message:  "HTTP " + r.Status,
		response: response,
	}
}

func NewJsonRpcClient(hc *http.Client, endpoint string) *JsonRpcClient {
	return &JsonRpcClient{hc: hc, Endpoint: endpoint, CustomHeader: make(map[string]string)}
}

func (c *JsonRpcClient) _do(req *http.Request) (resp *http.Response, err error) {
	if c.Pre != nil {
		if err = c.Pre(req); err != nil {
			return nil, err
		}
	}
	resp, err = c.hc.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http-status(%s) is not StatusOK", resp.Status)
		return
	}
	return
}

func (c *JsonRpcClient) Do(method string, reqPtr, respPtr interface{}) (jrResp *Response, err error) {
	jrReq := &jsonrpc.Request{
		ID:      time.Now().UnixNano() / int64(time.Millisecond),
		Version: jsonrpc.Version,
		Method:  method,
	}
	if reqPtr != nil {
		b, mErr := json.Marshal(reqPtr)
		if mErr != nil {
			err = mErr
			return
		}
		jrReq.Params = json.RawMessage(b)
	}
	reqB, err := json.Marshal(jrReq)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(reqB))
	if err != nil {
		return
	}
	req.Header.Set(headerContentType, typeApplicationJSON)
	req.Header.Set(headerAccept, typeApplicationJSON)
	for k, v := range c.CustomHeader {
		req.Header.Set(k, v)
	}

	var dErr error
	resp, err := c._do(req)
	if err != nil {
		if resp != nil {
			if ct, _, mErr := mime.ParseMediaType(resp.Header.Get(headerContentType)); mErr != nil {
				err = mErr
				return
			} else if ct == typeApplicationJSON {
				if jrResp, dErr = decodeResponseBody(resp); dErr != nil {
					err = dErr
					return
				}
			} else {
				err = NewHttpError(resp)
				return
			}
			err = jrResp.Error
			return
		}
		return
	}

	if jrResp, dErr = decodeResponseBody(resp); dErr != nil {
		err = fmt.Errorf("fail to decode response body err:%+v, jsonrpcResp:%+v",
			dErr, resp)
		return
	}
	if jrResp.Error != nil {
		err = jrResp.Error
		return
	}
	if respPtr != nil {
		err = json.Unmarshal(jrResp.Result, respPtr)
		if err != nil {
			return
		}
	}
	return
}

func (c *JsonRpcClient) Raw(reqB []byte) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(reqB))
	if err != nil {
		return
	}
	req.Header.Set(headerContentType, typeApplicationJSON)
	req.Header.Set(headerAccept, typeApplicationJSON)
	for k, v := range c.CustomHeader {
		req.Header.Set(k, v)
	}

	return c._do(req)
}

func decodeResponseBody(resp *http.Response) (jrResp *Response, err error) {
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&jrResp)
	return
}
