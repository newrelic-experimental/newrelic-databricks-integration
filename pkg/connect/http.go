package connect

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HttpMethod int

const (
	Undefined HttpMethod = iota
	Get
	Post
	Head
	Put
	Delete
	Connect
	Options
	Trace
	Patch
)

const defaultTimeout = 5 * time.Second

type BuilderFunc = func(*HttpConfig) (*http.Request, error)

type HttpConfig struct {
	Method  HttpMethod
	Url     string
	Body    any
	Headers map[string]string
	Builder BuilderFunc
	Timeout time.Duration
}

type HttpConnector struct {
	Config HttpConfig
}

func (c *HttpConnector) SetConfig(config any) {
	conf, ok := config.(HttpConfig)
	if !ok {
		panic("Expected an HttpConfig struct.")
	}
	c.Config = conf
}

func (c *HttpConnector) Request() ([]byte, ConnecError) {
	if c.Config.Builder != nil {
		req, err := c.Config.Builder(&c.Config)
		if err != nil {
			return nil, MakeConnectErr(err, 0)
		}
		res, reqErr := httpRequest(req, c.Config.Timeout)
		return []byte(res), reqErr
	} else {
		switch c.Config.Method {
		case Get:
			res, err := httpGet(c.Config.Url, c.Config.Headers, c.Config.Timeout)
			return []byte(res), err
		case Post:
			switch body := c.Config.Body.(type) {
			case io.Reader:
				res, err := httpPost(c.Config.Url, body, c.Config.Headers, c.Config.Timeout)
				return []byte(res), err
			case string:
				res, err := httpPost(c.Config.Url, strings.NewReader(body), c.Config.Headers, c.Config.Timeout)
				return []byte(res), err
			case []byte:
				res, err := httpPost(c.Config.Url, bytes.NewReader(body), c.Config.Headers, c.Config.Timeout)
				return []byte(res), err
			default:
				return nil, MakeConnectErr(errors.New("Unsupported type for body"), 0)
			}
		default:
			return nil, MakeConnectErr(errors.New("Unsupported HTTP method"), 0)
		}
	}
}

func (c *HttpConnector) ConnectorID() string {
	return "HTTP"
}

func (c *HttpConnector) SetReqBuilder(builder BuilderFunc) {
	c.Config.Builder = builder
}

func (c *HttpConnector) SetMethod(method HttpMethod) {
	c.Config.Method = method
}

func (c *HttpConnector) SetUrl(url string) {
	c.Config.Url = url
}

func (c *HttpConnector) SetBody(body any) {
	c.Config.Body = body
}

func (c *HttpConnector) SetHeaders(headers map[string]string) {
	c.Config.Headers = headers
}

func (c *HttpConnector) SetTimeout(timeout time.Duration) {
	c.Config.Timeout = timeout
}

// Build and Connector for HTTP GET requests.
func MakeHttpGetConnector(url string, headers map[string]string) HttpConnector {
	return HttpConnector{
		Config: HttpConfig{
			Method:  Get,
			Url:     url,
			Headers: headers,
			Timeout: defaultTimeout,
		},
	}
}

// Build and Connector for HTTP POST requests.
func MakeHttpPostConnector(url string, body any, headers map[string]string) HttpConnector {
	return HttpConnector{
		Config: HttpConfig{
			Method:  Post,
			Url:     url,
			Body:    body,
			Headers: headers,
			Timeout: defaultTimeout,
		},
	}
}

// Build and Connector for HTTP GET requests.
func MakeHttpConnectorWithBuilder(builder BuilderFunc) HttpConnector {
	return HttpConnector{
		Config: HttpConfig{
			Builder: builder,
			Timeout: defaultTimeout,
		},
	}
}

func httpGet(url string, headers map[string]string, timeout time.Duration) (string, ConnecError) {
	req, err := http.NewRequest("GET", url, nil)

	for key, val := range headers {
		req.Header.Add(key, val)
	}

	client := http.DefaultClient
	client.Timeout = timeout
	resp, err := client.Do(req)

	if err != nil {
		if resp != nil {
			return "", MakeConnectErr(err, resp.StatusCode)
		} else {
			return "", MakeConnectErr(err, 0)
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", MakeConnectErr(err, resp.StatusCode)
	}

	return string(body), ConnecError{}
}

func httpPost(url string, reqBody io.Reader, headers map[string]string, timeout time.Duration) (string, ConnecError) {
	req, err := http.NewRequest("POST", url, reqBody)

	for key, val := range headers {
		req.Header.Add(key, val)
	}

	client := http.DefaultClient
	client.Timeout = timeout
	resp, err := client.Do(req)

	if err != nil {
		if resp != nil {
			return "", MakeConnectErr(err, resp.StatusCode)
		} else {
			return "", MakeConnectErr(err, 0)
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", MakeConnectErr(err, resp.StatusCode)
	}

	return string(body), ConnecError{}
}

func httpRequest(req *http.Request, timeout time.Duration) (string, ConnecError) {
	client := http.DefaultClient
	client.Timeout = timeout
	resp, err := client.Do(req)

	if err != nil {
		if resp != nil {
			return "", MakeConnectErr(err, resp.StatusCode)
		} else {
			return "", MakeConnectErr(err, 0)
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", MakeConnectErr(err, resp.StatusCode)
	}

	return string(body), ConnecError{}
}
