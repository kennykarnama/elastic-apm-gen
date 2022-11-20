package fakeclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kennykarnama/elastic-apm-gen/provider"
)

type Resp struct {
	ID string `json:"id"`
}

type Client interface {
	GetData(ctx context.Context) (*Resp, error)
}

type httpClient struct {
	cli     *http.Client
	baseUrl provider.BaseURL
}

func NewHttpClient(cli *http.Client, baseUrl provider.BaseURL) *httpClient {
	return &httpClient{cli: cli, baseUrl: baseUrl}
}

func (c *httpClient) GetData(ctx context.Context) (*Resp, error) {
	targetUrl := fmt.Sprintf("%s%s", c.baseUrl.Get(ctx), "/api/v1/data")
	req, err := http.NewRequestWithContext(ctx, "GET", targetUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("action=fake.getData err=%w", err)
	}
	resp, err := req.GetBody()
	if err != nil {
		return nil, fmt.Errorf("action=fake.getData err=%w", err)
	}
	if resp != nil {
		defer resp.Close()
	}
	var b []byte
	_, err = resp.Read(b)
	if err != nil {
		return nil, fmt.Errorf("action=fake.getData err=%w", err)
	}
	var result Resp
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("action=fake.getData err=%w", err)
	}
	return &result, nil
}
