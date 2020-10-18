package httpclient

import (
	"fmt"
	"testing"

	"github.com/go-jar/golog"
)

func TestGet(t *testing.T) {
	client := newClient()
	resp, err := client.Get("http://127.0.0.1:8010/demo/DescribeDemo")
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

	resp, err := client.Post("http://127.0.0.1:8010/demo/ProcessPost", params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(resp.Contents), resp.TimeDuration.String())
}

func newClient() *Client {
	logger, _ := golog.NewConsoleLogger(golog.LEVEL_INFO)
	config := NewConfig()
	client := NewClient(config, logger)
	return client
}
