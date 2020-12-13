package httpclient

import (
	"bytes"
	"fmt"
	"github.com/goinbox/gomisc"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-jar/golog"
)

type Client struct {
	config  *Config
	logger  golog.ILogger
	retry   int
	traceId []byte

	*http.Client
}

type Request struct {
	Method  string
	Url     string
	Body    []byte
	Ip      string
	Headers map[string]string

	*http.Request
}

type Response struct {
	TimeDuration time.Duration
	Contents     []byte

	*http.Response
}

func NewClient(cfg *Config, logger golog.ILogger) *Client {
	return &Client{
		config:  NewConfig(),
		logger:  logger,
		traceId: []byte("-"),

		Client: &http.Client{
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   cfg.Timeout,
					KeepAlive: cfg.KeepAliveTime,
				}).DialContext,
				DisableKeepAlives:   cfg.DisableKeepAlives,
				MaxIdleConns:        cfg.MaxIdleConns,
				MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
				IdleConnTimeout:     cfg.IdleConnTimeout,
			},
		},
	}
}

func (c *Client) SetTraceId(traceId []byte) *Client {
	c.traceId = traceId

	return c
}

func (c *Client) Get(url string, headers map[string]string, ip string, retry int) (*Response, error) {
	req, err := NewRequest(http.MethodGet, url, nil, headers, ip)
	if err != nil {
		return nil, err
	}

	return c.Do(req, retry)
}

func (c *Client) Post(url string, data []byte, headers map[string]string, ip string, retry int) (*Response, error) {
	req, err := NewRequest(http.MethodGet, url, data, headers, ip)
	if err != nil {
		return nil, err
	}

	return c.Do(req, retry)
}

func (c *Client) MakePostBodyUrlEncode(data map[string]interface{}) []byte {
	values := url.Values{}
	for key, value := range data {
		values.Add(key, fmt.Sprint(value))
	}

	body := []byte(values.Encode())
	return body
}

func (c *Client) Do(req *Request, retry int) (*Response, error) {
	resp, timeDuration, err := c.do(req)
	if err != nil {
		for i := 0; i < retry; i++ {
			resp, timeDuration, err = c.do(req)
			if err == nil && resp.StatusCode == 200 {
				break
			}
		}
	}

	msg := [][]byte{
		[]byte("Host: " + req.Host),
		[]byte("URL: " + req.Url),
		[]byte("TimeDuration: " + timeDuration.String()),
	}
	if err != nil {
		if resp != nil {
			msg = append(msg, []byte("StatusCode: "+strconv.Itoa(resp.StatusCode)))
		}
		msg = append(msg, []byte("Error: "+err.Error()))
		c.logger.Error(c.fmtLog(bytes.Join(msg, []byte("\t"))))
		return nil, err
	}

	msg = append(msg, []byte("StatusCode: "+strconv.Itoa(resp.StatusCode)))
	c.logger.Log(c.config.LogLevel, c.fmtLog(bytes.Join(msg, []byte("\t"))))

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		TimeDuration: timeDuration,
		Contents:     contents,
		Response:     resp,
	}, nil
}

func (c *Client) fmtLog(msg []byte) []byte {
	return gomisc.AppendBytes(
		c.traceId, []byte("\t"),
		[]byte("[HttpClient]\t"),
		msg,
	)
}

func (c *Client) do(req *Request) (*http.Response, time.Duration, error) {
	start := time.Now()
	resp, err := c.Client.Do(req.Request)
	t := time.Since(start)

	if err != nil {
		return resp, t, err
	}

	return resp, t, nil
}

func NewRequest(method string, url string, body []byte, headers map[string]string, ip string) (*Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Host = req.URL.Host

	if ip != "" {
		s := strings.Split(req.URL.Host, ":")
		s[0] = ip
		req.URL.Host = strings.Join(s, ":")
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	return &Request{
		Method:  method,
		Url:     url,
		Body:    body,
		Headers: headers,
		Ip:      ip,
		Request: req,
	}, nil
}
