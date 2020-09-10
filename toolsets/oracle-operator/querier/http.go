package querier

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

// packRequest packs relayer ServiceRequest to http request.
func packRequest(ctx types.Context, msgCreateTask types.MsgCreateTask, endpoint string) (*http.Request, error) {
	var err error
	var req *http.Request
	logger := ctx.Logger().With("type", "msgCreateTask")
	config := ctx.Config()

	var requestBody []byte

	if strings.EqualFold(config.Querier.Method, "GET") {
		u, err := url.Parse(endpoint)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
		endpoint = u.String()
		logger.Debug("query endpoint with GET method", "endpoint", endpoint)
	} else if strings.EqualFold(config.Querier.Method, "POST") {
		requestBody, err = json.Marshal(msgCreateTask)
		if err != nil {
			logger.Error(err.Error())
			return nil, err
		}
	}

	req, err = http.NewRequest(strings.ToUpper(config.Querier.Method), endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return req, err
}

// packResponse packs http response with previous passed relayer ServiceRequest to ServiceResponse.
func packResponse(ctx types.Context, response *http.Response) (uint8, error) {
	logger := ctx.Logger().With("type", "response")
	logger.Debug("to pack response", "data", response.Body)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("read response body", "error", err.Error(), "body", response.Body)
		return 0, err
	}

	var sResData uint8
	err = json.Unmarshal(body, &sResData)
	if err != nil {
		logger.Error("unmarshal response data", "error", err.Error(), "body", body)
		return 0, err
	}
	return sResData, nil
}

// HandleRequest provides interface for relayer to call service provider endpoint.
func HandleRequest(ctx types.Context, msgCreateTask types.MsgCreateTask, endpoint string) (uint8, error) {
	// TODO: revise to async call
	var response *http.Response
	var RETRY = ctx.Config().Querier.RetryTimes

	// pack ServiceRequest to be an http request
	req, err := packRequest(ctx, msgCreateTask, endpoint)
	if err != nil {
		ctx.Logger().Error("pack request", "error", err.Error())
		return 0, err
	}
	client := &http.Client{
		Timeout: time.Duration(ctx.Config().Querier.Timeout) * time.Second,
	}

	// RETRY is 3 by default
	for i := 0; i < RETRY; i++ {
		response, err = client.Do(req)
		if response == nil {
			ctx.Logger().Info("Response empty")
			continue
		}
		if response.Body == nil {
			ctx.Logger().Info("Response Body empty")
			continue
		}
		if err == nil {
			defer response.Body.Close()
			break
		}
		defer response.Body.Close()

		ctx.Logger().Info("Retrying requests", "retries", i, "reason", err.Error())
	}
	if response == nil {
		ctx.Logger().Info("Response empty")
		return 0, fmt.Errorf("response empty")
	}
	if response.StatusCode != 200 {
		ctx.Logger().Info("Request not success", "status", response.Status)
		return 0, fmt.Errorf("request not success")
	}

	sres, err := packResponse(ctx, response)
	if err != nil {
		ctx.Logger().Error("pack response", "error", err.Error())
		return 0, err
	}
	ctx.Logger().Debug("Query Security Score from endpoint", "score", sres)
	return sres, nil
}
