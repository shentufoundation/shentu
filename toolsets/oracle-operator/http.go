package oracle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/certikfoundation/shentu/toolsets/oracle-operator/types"
)

// packRequest packs payload to http request.
func packRequest(ctx types.Context, endpoint string, payload interface{}) (*http.Request, error) {
	var err error
	var req *http.Request
	logger := ctx.Logger().With("type", "request", "endpiont", endpoint, "payload", payload)
	config := ctx.Config()

	var requestBody []byte

	if strings.EqualFold(config.Method, "GET") {
		u, err := url.Parse(endpoint)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		endpoint = u.String()
		logger.Debug("query endpoint with GET method", "endpoint", endpoint)
	} else if strings.EqualFold(config.Method, "POST") {
		requestBody, err = json.Marshal(payload)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
	}

	req, err = http.NewRequest(strings.ToUpper(config.Method), endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return req, err
}

// packResponse packs http response to single score.
func packResponse(ctx types.Context, response *http.Response) (uint8, error) {
	logger := ctx.Logger().With("type", "response")
	logger.Debug("to pack response", "data", response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("read response body", "error", err.Error(), "body", response.Body)
		return 0, err
	}

	var result uint8
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Error("unmarshal response data", "error", err.Error(), "body", body)
		return 0, err
	}
	return result, nil
}

// handleRequest queries primitive endpoint.
func handleRequest(ctx types.Context, endpoint string, payload interface{}) (uint8, error) {
	var response *http.Response
	var RETRY = ctx.Config().RetryTimes
	var TIMEOUT = ctx.Config().Timeout
	// pack task to be an http request
	req, err := packRequest(ctx, endpoint, payload)
	if err != nil {
		ctx.Logger().Error("packing request", "error", err.Error())
		return 0, err
	}
	client := &http.Client{
		Timeout: time.Duration(TIMEOUT) * time.Second,
	}
	// make http query
	for i := 0; i < RETRY; i++ {
		response, err = client.Do(req)
		if response == nil {
			ctx.Logger().Debug("response empty")
			continue
		}
		if response.Body == nil {
			ctx.Logger().Debug("response Body empty")
			continue
		}
		if err == nil {
			defer response.Body.Close()
			break
		}
		defer response.Body.Close()

		ctx.Logger().Debug("retrying requests", "retries", i, "reason", err.Error())
	}
	if response == nil {
		ctx.Logger().Debug("response empty")
		return 0, fmt.Errorf("response empty")
	}
	if response.StatusCode != 200 {
		ctx.Logger().Debug("request not success", "status", response.Status)
		return 0, fmt.Errorf("request not success")
	}
	// extract result from http response
	result, err := packResponse(ctx, response)
	if err != nil {
		ctx.Logger().Error("pack response", "error", err.Error())
		return 0, err
	}
	ctx.Logger().Debug("result from endpoint", "result", result)
	return result, nil
}
