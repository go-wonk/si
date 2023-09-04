package sielastic

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/go-wonk/si/v2/sicore"
)

func DefaultElasticsearchClient(elasticAddresses, userName, password string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: strings.Split(elasticAddresses, ","),
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Duration(6) * time.Second,
		},
		Username: userName,
		Password: password,
	}
	return elasticsearch.NewClient(cfg)
}

type Client struct {
	*elasticsearch.Client
}

func NewClient(client *elasticsearch.Client) *Client {
	return &Client{client}
}

func (c *Client) IndexDocument(ctx context.Context, indexName string, body []byte) (map[string]interface{}, error) {
	req := esapi.IndexRequest{
		Index: indexName,
		// DocumentID: docID,
		Body: bytes.NewReader(body),
		// Refresh:    strconv.FormatBool(refresh),
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

// func (c *Client) SearchDocuments(ctx context.Context, indexName string, body map[string]interface{}) (map[string]interface{}, error) {

// 	buf := sicore.GetBytesBuffer(nil)
// 	defer sicore.PutBytesBuffer(buf)

// 	if err := sicore.EncodeJson(buf, body); err != nil {
// 		return nil, err
// 	}

// 	res, err := c.Search(
// 		c.Search.WithContext(ctx),
// 		c.Search.WithIndex(indexName),
// 		c.Search.WithBody(buf),
// 		c.Search.WithTrackTotalHits(true),
// 		c.Search.WithTrackScores(true),
// 		// c.Search.WithPretty(),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()

// 	var r map[string]interface{}
// 	if err := sicore.DecodeJson(&r, res.Body); err != nil {
// 		return nil, err
// 	}

// 	if res.IsError() {
// 		return r, fmt.Errorf("[%s] %s: %s", res.Status(), r["error"].(map[string]interface{})["type"], r["error"].(map[string]interface{})["reason"])
// 	}

// 	return r, nil
// }

var ErrElasticResponseHasError = errors.New("response has error")

type RespRootCause struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
	Line   int    `json:"line"`
	Col    int    `json:"col"`
}
type RespError struct {
	RootCause []RespRootCause `json:"root_cause"`
	Reason    string          `json:"reason"`
}

type Resp struct {
	Error  RespError `json:"error"`
	Status int       `json:"status"`
}

func (c *Client) SearchDocuments(ctx context.Context, indexName string, body map[string]interface{}, dest any) error {

	buf := sicore.GetBytesBuffer(nil)
	defer sicore.PutBytesBuffer(buf)

	if err := sicore.EncodeJson(buf, body); err != nil {
		return err
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
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		resErr := Resp{}
		copied, err := sicore.DecodeJsonCopied(&resErr, res.Body)
		if err != nil {
			return err
		}

		if err := sicore.DecodeJson(&dest, copied); err != nil {
			return err
		}
		return errors.New(resErr.Error.Reason)
	}

	if err := sicore.DecodeJson(dest, res.Body); err != nil {
		return err
	}

	return nil
}

/*
{
  "took": 3,
  "timed_out": false,
  "_shards": {
    "total": 1,
    "successful": 1,
    "skipped": 0,
    "failed": 0
  },
  "hits": {
    "total": {
      "value": 36,
      "relation": "eq"
    },
    "max_score": null,
    "hits": []
  },
  "aggregations": {
    "idCount": {
      "doc_count_error_upper_bound": 0,
      "sum_other_doc_count": 0,
      "buckets": [
        {
          "key": "http-adpos-nginx",
          "doc_count": 36
        }
      ]
    }
  }
}

{
  "error": {
    "root_cause": [
      {
        "type": "parsing_exception",
        "reason": "Found two aggregation type definitions in [idCount]: [terms] and [having]",
        "line": 26,
        "col": 19
      }
    ],
    "type": "parsing_exception",
    "reason": "Found two aggregation type definitions in [idCount]: [terms] and [having]",
    "line": 26,
    "col": 19
  },
  "status": 400
}
*/
