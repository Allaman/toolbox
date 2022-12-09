package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"time"

	"github.com/allaman/toolbox/http-client/sophisticated/logger"
)

// Client provides basic data for the http client
type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

// ClientOptions struct for functional options
type ClientOptions func(*Client)

// CreateClient returns a Client
func CreateClient(opts ...ClientOptions) *Client {
	const (
		defaultBaseURL        = "https://example.com"
		defaultTimeout        = 30
		defaultDumpingEnabled = false
	)
	c := &Client{
		httpClient: &http.Client{},
		baseURL:    defaultBaseURL,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func withInsecureSSL() ClientOptions {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return func(c *Client) {
		c.httpClient.Transport = customTransport
	}
}

func withTimeout(t int) ClientOptions {
	logger.Log.Debug().Msgf("setting timeout to '%d'", t)
	return func(c *Client) {
		c.httpClient.Timeout = time.Duration(time.Duration(t) * time.Second)
	}
}

func withBaseURL(u string) ClientOptions {
	logger.Log.Debug().Msgf("setting baseurl to '%s'", u)
	return func(c *Client) {
		c.baseURL = u
	}
}

func withDumpingEnabled() ClientOptions {
	logger.Log.Debug().Msgf("setting roundtripper")
	return func(c *Client) {
		c.httpClient.Transport = &GeneralRoundTripper{
			transport:   http.DefaultTransport,
			interceptor: DumpingInterceptor,
		}
	}
}

// NewRequest returns a http.Request according to its parameter and header
func (c *Client) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.baseURL, path), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	req.Header.Set("Accept", "application/json")
	logger.Log.Debug().Msgf("Request is %v", req)
	return req, nil
}

// Do performs the http request
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	localReq := req
	if logger.Log.Trace().Enabled() {
		localReq = c.reqWithTracing(req)
	}
	resp, err := c.httpClient.Do(localReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) reqWithTracing(req *http.Request) *http.Request {
	clientTrace := &httptrace.ClientTrace{
		GetConn: func(hostPort string) { logger.Log.Trace().Str("host", hostPort).Msgf("starting to create conn") },
		DNSStart: func(info httptrace.DNSStartInfo) {
			logger.Log.Trace().Str("Lookup", info.Host).Msgf("starting to look up dns")
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			var addresses string
			for _, add := range info.Addrs {
				addresses = fmt.Sprintf("%s, %s", addresses, add.String())
			}
			logger.Log.Trace().Str("Address", addresses).Msgf("done looking up dns")
		},
		ConnectStart: func(network, addr string) {
			logger.Log.Trace().Str("network", network).Str("addr", addr).Msgf("starting tcp connection")
		},
		ConnectDone: func(network, addr string, err error) {
			logger.Log.Trace().Str("network", network).Str("addr", addr).Err(err).Msgf("tcp connection created")
		},
		GotConn: func(info httptrace.GotConnInfo) {
			logger.Log.Trace().Bool("reused", info.Reused).Bool("wasIdle", info.WasIdle).Str("idleTime", info.IdleTime.String()).Msgf("connection established")
		},
	}
	clientTraceCtx := httptrace.WithClientTrace(req.Context(), clientTrace)
	req = req.WithContext(clientTraceCtx)
	return req
}

type Debug struct {
	DNS struct {
		Start   string       `json:"start"`
		End     string       `json:"end"`
		Host    string       `json:"host"`
		Address []net.IPAddr `json:"address"`
		Error   error        `json:"error"`
	} `json:"dns"`
	Dial struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"dial"`
	Connection struct {
		Time string `json:"time"`
	} `json:"connection"`
	WroteAllRequestHeaders struct {
		Time string `json:"time"`
	} `json:"wrote_all_request_header"`
	WroteAllRequest struct {
		Time string `json:"time"`
	} `json:"wrote_all_request"`
	FirstReceivedResponseByte struct {
		Time string `json:"time"`
	} `json:"first_received_response_byte"`
}

type Handler func(*http.Request) (*http.Response, error)

type Interceptor func(*http.Request, Handler) (*http.Response, error)

type GeneralRoundTripper struct {
	transport   http.RoundTripper
	interceptor Interceptor
}

func (general *GeneralRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return general.interceptor(req, general.transport.RoundTrip)
}

func DumpingInterceptor(req *http.Request, handler Handler) (*http.Response, error) {
	bytes, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		logger.Log.Fatal().Err(err).Msgf("error while dumping http request")
	}
	logger.Log.Trace().Str("request", string(bytes)).Msgf("request dump")
	response, err := handler(req)
	if err != nil {
		logger.Log.Fatal().Err(err).Msgf("error while handling http response")
	}
	bytes, err = httputil.DumpResponse(response, true)
	if err != nil {
		logger.Log.Fatal().Err(err).Msgf("error while dumping http response")
	}
	logger.Log.Trace().Str("response", string(bytes)).Msgf("response dump")
	return response, err
}
