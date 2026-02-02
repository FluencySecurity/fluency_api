package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	fsb "github.com/SecurityDo/fluency_api/fsb"
)

type HTTPService struct {
	url       string
	client    *http.Client
	DebugFlag bool
	token     string
	logger    *slog.Logger
}

func NewHTTPService(url string, logger *slog.Logger) *HTTPService {

	s := &HTTPService{
		url:    url,
		logger: logger,
	}
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: 10 * time.Second,
	}
	timeout := time.Duration(600 * time.Second)
	s.client = &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	return s
}

func (r *HTTPService) Close() {
	r.client.CloseIdleConnections()
}
func (r *HTTPService) GetClient() *http.Client {
	return r.client
}

func (r *HTTPService) GetUrl() string {
	return r.url
}

func (r *HTTPService) SetToken(token string) {
	r.token = token
}

func (r *HTTPService) Call(prefix string, functionName string, input interface{}) (result *fsb.JNode, err error) {
	remoteReq := new(fsb.CallRequest)
	remoteReq.Function = functionName

	if input == nil {
		remoteReq.Kargs = fsb.NewJNodeString("{}")
	} else {
		remoteReq.Kargs, _ = fsb.NewJNodeInterface(input)
	}

	reqStr, err := json.Marshal(remoteReq)
	if err != nil {
		return nil, fmt.Errorf("args is not a valid json structure")
	}

	fullUrl := fmt.Sprintf("%s/%s/%s", r.url, prefix, functionName)
	if prefix == "" {
		fullUrl = fmt.Sprintf("%s/%s", r.url, functionName)
	}
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(reqStr))
	req.Header.Set("Content-Type", "application/json")
	if r.token != "" {
		// bearer token
		req.Header.Set("Authorization", "Bearer "+r.token)
	}

	if r.DebugFlag {
		if dump, err := httputil.DumpRequestOut(req, false); err == nil {
			r.logger.Debug("Request:\n-----------------------------------------\n")
			r.logger.Debug(string(dump))
			pretty, _ := json.MarshalIndent(remoteReq, "", "   ")
			r.logger.Debug(string(pretty))
		}
	}
	resp, err := r.client.Do(req)
	if err != nil {
		r.logger.Error("Failed to call local http service %s: %s\n", r.url, err.Error())
		return result, fmt.Errorf("failed to call local http service %s: %s", r.url, err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.logger.Error("Failed to read response body", "Error", err.Error())
		return result, fmt.Errorf("failed to read response body: %w", err)
	}
	if r.DebugFlag {
		if dump, err := httputil.DumpResponse(resp, false); err == nil {
			r.logger.Debug("Response:\n-----------------------------------------\n")
			r.logger.Debug(string(dump))
		}
		r.logger.Debug("Response body:\n-----------------------------------------\n")
		r.logger.Debug(string(body))
	}
	if resp.StatusCode != 200 {
		r.logger.Error("HTTP ERROR from local http service %s: %s\n", r.url, resp.Status)
		return result, fmt.Errorf("HTTP Error from Local HTTP service %s: %s", r.url, resp.Status)
	}
	var res fsb.CallResponse
	//var obj map[string]interface{}
	//err = json.Unmarshal(body,&res)
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	err = decoder.Decode(&res)
	if err != nil {
		r.logger.Error("Failed to parse response body -> ", "Error", err.Error())
		return result, fmt.Errorf("HTTP Error from Local HTTP service: %s", resp.Status)
	}
	if r.DebugFlag {
		pretty, _ := json.MarshalIndent(res, "", "   ")
		r.logger.Debug(string(pretty))
	}

	if res.Verdict == "ERROR" {
		r.logger.Error("RPC call return with ERROR", "prefix", prefix, "functionName", functionName, "Error", res.Error)
		return result, fmt.Errorf("RPC call return with ERROR: %s", res.Error)
	} else if res.Verdict == "EXCEPTION" {
		r.logger.Debug("RPC call return with EXCEPTION: ", "exception", res.Exception)
		return result, fmt.Errorf("RPC call return with EXCEPTION: %s", res.Exception)
	}

	return res.Response, nil

}

type FluencyClient struct {
	serviceClient *HTTPService
	logger        *slog.Logger
}

func NewFluencyClient(siteURL string, token string, debugFlag bool, logger *slog.Logger) *FluencyClient {
	if logger == nil {
		// Using os.Stderr by default is safe for libraries
		logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	}

	s := &FluencyClient{
		serviceClient: NewHTTPService(siteURL, logger),
		logger:        logger,
	}
	s.serviceClient.DebugFlag = debugFlag
	s.serviceClient.SetToken(token)

	return s

}

// SetDebug enables or disables HTTP request/response dump logging.
func (r *FluencyClient) SetDebug(debug bool) {
	r.serviceClient.DebugFlag = debug
}

func (r *FluencyClient) GenericCall(prefix string, functionName string, x interface{}) (res *fsb.JNode, err error) {
	res, err = r.serviceClient.Call(prefix, functionName, x)

	if err != nil {
		fmt.Printf("%s/%s/%s call failed", r.serviceClient.GetUrl(), prefix, functionName)
		return res, err
	}
	return res, err
}
