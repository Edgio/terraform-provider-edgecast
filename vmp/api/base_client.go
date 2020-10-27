// Copyright Verizon Media, Licensed under the terms of the Apache 2.0 license . See LICENSE file in project root for terms.

package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
)

type ApiClient struct {
	BaseUrl   *url.URL
	UserAgent string
	Token     string

	HttpClient *retryablehttp.Client
}

func NewApiClient(apiBaseUri string, apiToken string) (*ApiClient, error) {
	baseUrl, err := url.Parse(apiBaseUri)

	if err != nil {
		return nil, err
	}

	return &ApiClient{
		BaseUrl:    baseUrl,
		Token:      apiToken,
		HttpClient: retryablehttp.NewClient(),
	}, nil
}

func (c *ApiClient) BuildRequest(method, path string, body interface{}) (*retryablehttp.Request, error) {

	rel, pathErr := url.Parse(path)

	if pathErr != nil {
		return nil, pathErr
	}

	u := c.BaseUrl.ResolveReference(rel)

	var buf io.ReadWriter

	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := retryablehttp.NewRequest(method, u.String(), buf)

	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "TOK:"+c.Token)
	req.Header.Set("User-Agent", "verizonmedia/terraform:1.0.0")

	return req, nil

}

func (c *ApiClient) SendRequest(req *retryablehttp.Request, v interface{}) (*http.Response, error) {

	resp, err := c.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}

		bodyAsString := string(body)
		return nil, errors.New(bodyAsString)
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err

}
