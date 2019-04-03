package xhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/buger/jsonparser"
)

const (
	defaultConnTimeout  = 3000
	defaultRespTimeout  = 5000
	defaultTotalTimeout = 30000

	minHttpUrlLen = len("http://")

	ErrorDefault = -iota - 1 //-1
	ErrorInvalidUrl
	ErrorInvalidMethod
)

var (
	errMsg = map[int]error{
		ErrorDefault:       errors.New("default error"),
		ErrorInvalidUrl:    errors.New("invalid url, lost http/https?"),
		ErrorInvalidMethod: errors.New("invalid method"),
	}
)

var allowMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodPost:    true,
	http.MethodHead:    true,
	http.MethodPut:     true,
	http.MethodDelete:  true,
	http.MethodOptions: true,
}

type XHttp struct {
	body    string
	method  string
	client  *http.Client
	params  map[string]string
	headers map[string]string
}

// set the http method
func (c *XHttp) Method(m string) *XHttp {
	c.method = m
	return c
}

// set the http method to get
func (c *XHttp) Get() *XHttp {
	c.method = "GET"
	return c
}

// set the http method to post
func (c *XHttp) Post() *XHttp {
	c.method = "POST"
	return c
}

// add a http header
func (c *XHttp) AddHeader(k, v string) *XHttp {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}
	c.headers[k] = v
	return c
}

// add some http header
func (c *XHttp) AddHeaders(headerMap map[string]string) *XHttp {
	if c.headers == nil {
		c.headers = make(map[string]string)
	}
	for k, v := range headerMap {
		c.headers[k] = v
	}
	return c
}

// set the http uri param
func (c *XHttp) AddParam(k, v string) *XHttp {
	if c.params == nil {
		c.params = make(map[string]string)
	}
	c.params[k] = v
	return c
}

// add some http params
func (c *XHttp) AddParams(paramMap map[string]string) *XHttp {
	if c.params == nil {
		c.params = make(map[string]string)
	}
	for k, v := range paramMap {
		c.params[k] = v
	}
	return c
}

// set the http body
func (c *XHttp) SetBody(b string) *XHttp {
	c.body = b
	return c
}

// new with default timeout
func New() *XHttp {
	return NewWithTimeout(defaultConnTimeout, defaultRespTimeout, defaultTotalTimeout)
}

// new with timeout
func NewWithTimeout(connectTime, responseHeaderTimeout, totalTimeout int64) *XHttp {
	dialTimeout := func(network, addr string) (net.Conn, error) {
		dialer := &net.Dialer{
			Timeout: time.Millisecond * time.Duration(connectTime),
		}
		conn, err := dialer.Dial(network, addr)
		if err != nil {
			err = fmt.Errorf("net.DialTimeout, addr:%s, err:%v", addr, err)
			if conn != nil {
				err = fmt.Errorf("%v, conn:%v", err, conn.RemoteAddr())
			}
		}
		return conn, err
	}
	return &XHttp{
		client: &http.Client{
			Transport: &http.Transport{
				Dial:                  dialTimeout,
				ResponseHeaderTimeout: time.Millisecond * time.Duration(responseHeaderTimeout),
			},
			Timeout: time.Millisecond * time.Duration(totalTimeout),
		},
	}
}

// get all data from http request
func (c *XHttp) RespToString(url string) (string, error) {
	body, err := c.getRespBody(url)
	if err != nil {
		return "", err
	}
	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// get a json's value from http request
func (c *XHttp) RespGetJsonKey(url string, keys ...string) ([]byte, error) {
	body, err := c.getRespBody(url)
	if err != nil {
		return []byte{}, err
	}
	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return []byte{}, err
	}

	v, _, _, err := jsonparser.Get(data, keys...)

	return v, err
}

// get a http response body
// warning: you must close body at last
func (c *XHttp) getRespBody(url string) (body io.ReadCloser, err error) {
	if len(url) <= minHttpUrlLen || url[0:4] != "http" {
		return nil, errMsg[ErrorInvalidUrl]
	}

	if _, ok := allowMethods[c.method]; !ok {
		return nil, errMsg[ErrorInvalidMethod]
	}

	req, err := http.NewRequest(c.method, url, strings.NewReader(c.body))
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if c.params != nil {
		for key, val := range c.params {
			q.Add(key, val)
		}
		c.params = nil
		req.URL.RawQuery = q.Encode()
	}

	if c.headers != nil {
		for key, val := range c.headers {
			req.Header.Add(key, val)
		}
		c.headers = nil
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// get all data from http request
func (c *XHttp) RespToJson(url string, j interface{}) error {
	body, err := c.getRespBody(url)
	if err != nil {
		return err
	}
	err = json.NewDecoder(body).Decode(j)
	defer body.Close()

	return err
}
