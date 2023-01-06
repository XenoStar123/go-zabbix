package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type (
	Params map[string]interface{}
)

type request struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Auth    string      `json:"auth,omitempty"`
	Id      int32       `json:"id"`
}

type Response struct {
	Jsonrpc string          `json:"jsonrpc"`
	Error   *Error          `json:"error"`
	Result  json.RawMessage `json:"result"`
	Id      int32           `json:"id"`
}

func (resp *Response) Bind(v interface{}) (err error) {
	return json.Unmarshal(resp.Result, v)
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d (%s): %s", e.Code, e.Message, e.Data)
}

type API struct {
	auth   string      // auth token, filled by Login()
	Logger *log.Logger // request/response logger, nil by default
	url    string
	c      *http.Client
	id     int32
}

// Creates new API access object.
// Typical URL is http://host/api_jsonrpc.php or http://host/zabbix/api_jsonrpc.php.
// It also may contain HTTP basic auth username and password like
// http://username:password@host/api_jsonrpc.php.
func NewAPI(url string) *API {
	return &API{url: url, c: http.DefaultClient, id: 1}
}

// Allows one to use specific http.Client, for example with InsecureSkipVerify transport.
// This method modifies API structure and should not be called concurrently with other methods.
func (api *API) SetClient(c *http.Client) {
	api.c = c
}

func (api *API) printf(format string, v ...interface{}) {
	if api.Logger != nil {
		api.Logger.Printf(format, v...)
	}
}

func (api *API) callBytes(method string, params interface{}) ([]byte, error) {
	// id := atomic.AddInt32(&api.id, 1)
	jsonobj := request{"2.0", method, params, api.auth, api.id}
	b, err := json.Marshal(jsonobj)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %v", err)
	}
	api.printf("Request (%s): %s", http.MethodPost, b)

	req, err := http.NewRequest(http.MethodPost, api.url, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %v", err)
	}
	// req.Close = true
	req.ContentLength = int64(len(b))
	req.Header.Add("Content-Type", "application/json-rpc")
	// req.Header.Add("User-Agent", "github.com/XenoStar123/go-zabbix")

	resp, err := api.c.Do(req)
	if err != nil {
		api.printf("Error   : %s", err)
		return nil, fmt.Errorf("api.c.Do: %v", err)
	}
	defer resp.Body.Close()

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %v", err)
	}
	api.printf("Response (%d): %s", resp.StatusCode, b)
	return b, nil
}

// Calls specified API method. Uses api.Auth if not empty.
// err is something network or marshaling related. Caller should inspect response.Error to get API error.
func (api *API) Call(method string, params interface{}) (*Response, error) {
	b, err := api.callBytes(method, params)
	if err != nil {
		return nil, fmt.Errorf("api.callBytes: %v", err)
	}

	resp := &Response{}
	err = json.Unmarshal(b, resp)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}
	return resp, nil
}

// Uses Call() and then sets err to response.Error if former is nil and latter is not.
func (api *API) CallWithError(method string, params interface{}) (*Response, error) {
	resp, err := api.Call(method, params)
	if err != nil {
		return nil, fmt.Errorf("api.Call: %v", err)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("resp.Error: %v", resp.Error)
	}
	return resp, nil
}

// Calls "user.login" API method and fills api.Auth field.
// This method modifies API structure and should not be called concurrently with other methods.
func (api *API) Login(user, password string) (string, error) {
	api.auth = ""
	params := map[string]string{"user": user, "password": password}
	resp, err := api.CallWithError("user.login", params)
	if err != nil {
		return "", fmt.Errorf("api.CallWithError: %v", err)
	}

	auth := ""
	err = resp.Bind(&auth)
	if err != nil {
		return "", fmt.Errorf("resp.Bind: %v", err)
	}
	api.auth = auth
	return auth, nil
}

// Calls "APIInfo.version" API method.
// This method temporary modifies API structure and should not be called concurrently with other methods.
func (api *API) Version() (string, error) {
	auth := api.auth
	api.auth = ""
	resp, err := api.CallWithError("apiinfo.version", Params{})
	api.auth = auth
	if err != nil {
		return "", fmt.Errorf("api.CallWithError: %v", err)
	}

	v := ""
	err = resp.Bind(&v)
	if err != nil {
		return "", fmt.Errorf("resp.Bind: %v", err)
	}
	return v, nil
}
