package sielastic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/go-wonk/si/sicore"
)

func DefaultElasticsearchClient(elasticAddresses string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: strings.Split(elasticAddresses, ","),
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Duration(6) * time.Second,
		},
	}
	return elasticsearch.NewClient(cfg)
}

type Client struct {
	*elasticsearch.Client
}

func NewClient(client *elasticsearch.Client) *Client {
	return &Client{client}
}

func (c *Client) IndexDocument(ctx context.Context, indexName string, docID string, body []byte, refresh bool) (map[string]interface{}, error) {
	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       bytes.NewReader(body),
		Refresh:    strconv.FormatBool(refresh),
	}

	res, err := req.Do(ctx, c)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(res.Status())
	}

	var r map[string]interface{}
	if err := sicore.DecodeJson(&r, res.Body); err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) SearchDocuments(ctx context.Context, indexName string, body map[string]interface{}) (map[string]interface{}, error) {

	buf := sicore.GetBytesBuffer(make([]byte, 0, 128))
	defer sicore.PutBytesBuffer(buf)

	if err := sicore.EncodeJson(buf, body); err != nil {
		return nil, err
	}

	res, err := c.Search(
		c.Search.WithContext(ctx),
		c.Search.WithIndex(indexName),
		c.Search.WithBody(buf),
		c.Search.WithTrackTotalHits(true),
		c.Search.WithTrackScores(true),
		// c.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r map[string]interface{}
	if err := sicore.DecodeJson(&r, res.Body); err != nil {
		return nil, err
	}

	if res.IsError() {
		return r, fmt.Errorf("[%s] %s: %s", res.Status(), r["error"].(map[string]interface{})["type"], r["error"].(map[string]interface{})["reason"])
	}

	return r, nil
}
