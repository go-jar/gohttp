package httpclient

import (
	"fmt"
	"testing"

	"github.com/go-jar/golog"
)

func TestArgs(t *testing.T) {
	client := newClient()
	resp, err := client.Get("http://127.0.0.1:8010/demo/args/123", nil, "", 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(resp.Contents), resp.TimeDuration.String())

	resp, err = client.Get("http://127.0.0.1:8010/test/args/123", nil, "", 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(resp.Contents), resp.TimeDuration.String())
}

func TestGet(t *testing.T) {
	client := newClient()
	resp, err := client.Get("http://127.0.0.1:8010/demo/DescribeDemo", nil, "", 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(resp.Contents), resp.TimeDuration.String())
}

func TestPost(t *testing.T) {
	client := newClient()
	params := map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	body := client.MakePostBodyUrlEncode(params)
	resp, err := client.Post("http://127.0.0.1:8010/demo/ProcessPost", body, nil, "", 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(resp.Contents), resp.TimeDuration.String())
}

func newClient() *Client {
	logger, _ := golog.NewConsoleLogger(golog.LevelDebug)
	config := NewConfig()
	client := NewClient(config, logger)
	return client
}
