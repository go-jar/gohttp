package httpclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-jar/golog"
)

const (
	DEFAULT_RETRY = 3
)

type Client struct {
	config *Config
	logger golog.ILogger
	retry  int

	*http.Client
}

type Response struct {
	TimeDuration time.Duration
	Contents     []byte

	*http.Response
}

func NewClient(cfg *Config, l golog.ILogger) *Client {
	return &Client{
		config: NewConfig(),
		logger: l,
		retry:  DEFAULT_RETRY,
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

func (c *Client) SetRetry(retry int) {
	c.retry = retry
}

func (c *Client) Get(url string) (*Response, error) {
	return c.Do(http.MethodGet, url, nil)
}

func (c *Client) Post(ur string, data map[string]interface{}) (*Response, error) {
	values := url.Values{}
	for key, value := range data {
		values.Add(key, fmt.Sprint(value))
	}

	body := []byte(values.Encode())

	return c.Do(http.MethodPost, ur, body)
}

func (c *Client) Do(methodType string, url string, body []byte) (*Response, error) {
	resp, timeDuration, err := c.do(methodType, url, body)
	if err != nil {
		for i := 0; i < c.retry; i++ {
			resp, timeDuration, err = c.do(methodType, url, body)
			if err == nil && resp.StatusCode == 200 {
				break
			}
		}
	}

	msg := [][]byte{
		[]byte("Url: " + url),
		[]byte("TimeDuration: " + timeDuration.String()),
	}
	if err != nil {
		if resp != nil {
			msg = append(msg, []byte("StatusCode: "+strconv.Itoa(resp.StatusCode)))
		}
		msg = append(msg, []byte("Error: "+err.Error()))
		c.logger.Error(bytes.Join(msg, []byte("\t")))
		return nil, err
	}

	msg = append(msg, []byte("StatusCode: "+strconv.Itoa(resp.StatusCode)))
	c.logger.Info(bytes.Join(msg, []byte("\t")))

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

func (c *Client) do(methodType string, url string, body []byte) (*http.Response, time.Duration, error) {
	req, _ := http.NewRequest(methodType, url, bytes.NewReader(body))

	start := time.Now()
	resp, err := c.Client.Do(req)
	t := time.Since(start)

	if err != nil {
		return resp, t, err
	}

	return resp, t, nil
}
