package solana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// Request defines a JSON-RPC 2.0 request object. See
// https://www.jsonrpc.org/specification for more information. A Request should
// not be explicitly created, but instead unmarshaled from JSON.
type Request struct {
	Version string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response defines a JSON-RPC 2.0 response object. See
// https://www.jsonrpc.org/specification for more information. A Response is
// usually marshaled into bytes and returned in response to a Request.
type Response struct {
	Version string           `json:"jsonrpc"`
	ID      interface{}      `json:"id"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Error   *Error           `json:"error,omitempty"`
}

// Error defines a JSON-RPC 2.0 error object. See
// https://www.jsonrpc.org/specification for more information.
type Error struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data"`
}

// SendData sends data to method via jsonrpc
func SendData(method string, data []byte, url string) (Response, error) {
	request := Request{
		Version: "2.0",
		ID:      1,
		Method:  method,
		Params:  data,
	}
	// Send request to lightnode
	response, err := SendRequest(request, url)
	if err != nil {
		return Response{}, err
	}

	var resp Response
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		return Response{}, fmt.Errorf("cannot decode %v response body = %s, err = %v", method, buf.String(), err)
	}
	if resp.Error != nil {
		return Response{}, fmt.Errorf("got err back from %v request, err = %v", method, resp.Error)
	}
	return resp, nil
}

// SendDataWithRetry is the same as SendData but will retry if sending the request failed
func SendDataWithRetry(method string, data []byte, url string) (Response, error) {
	request := Request{
		Version: "2.0",
		ID:      1,
		Method:  method,
		Params:  data,
	}
	// Send request to lightnode with retry (max 10 times)
	response, err := SendRequestWithRetry(request, url, 10, 10)
	if err != nil {
		return Response{}, fmt.Errorf("failed to send request, err = %v", err)
	}

	var resp Response
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		return Response{}, fmt.Errorf("cannot decode %v response body = %s, err = %v", method, buf.String(), err)
	}
	if resp.Error != nil {
		return Response{}, fmt.Errorf("got err back from %v request, err = %v", method, resp.Error)
	}
	return resp, nil
}

// SendRequest sends the JSON-2.0 request to the target url and returns the response and any error.
func SendRequest(request Request, url string) (*http.Response, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	resp, err := SendRawPost(data, url)
	if err != nil {
		fmt.Printf("Sending %s to %s resulted in an error: %v\n", string(data), url, err)
		return nil, err
	}
	return resp, nil
}

// SendRawPost sends a raw bytes as a POST request to the URL specified
func SendRawPost(data []byte, url string) (*http.Response, error) {
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	client := newClient(10 * time.Second)
	buff := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", url, buff)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

// SendRequestWithRetry calls SendRequest but with configurable retry logic
func SendRequestWithRetry(request Request, url string, timeoutInSecs int, retries int) (response *http.Response, err error) {
	failures := 0
	for failures < retries {
		response, err = SendRequest(request, url)
		if err != nil {
			failures++
			if failures >= retries {
				return nil, err
			}
			fmt.Printf("%s errored: %v. Retrying after %d seconds\n", url, err, timeoutInSecs)
			time.Sleep(time.Duration(timeoutInSecs) * time.Second)
			continue
		}
		break
	}
	return
}

func newClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 4 * time.Second,
			ResponseHeaderTimeout: 3 * time.Second,
		},
		Timeout: timeout,
	}
}
