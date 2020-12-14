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
	"strconv"

	"github.com/hashicorp/go-retryablehttp"
)

type ApiClient struct {
	BaseUrl         *url.URL
	UserAgent       string
	Token           string
	IdsClientId     string
	IdsClientSecret string
	IdsScope        string

	HttpClient *retryablehttp.Client
}

func NewApiClient(apiBaseUri string, apiToken string, idsClientId string, idsClientSecret string, idsScope string) (*ApiClient, error) {
	baseUrl, err := url.Parse(apiBaseUri)

	if err != nil {
		return nil, err
	}

	return &ApiClient{
		BaseUrl:         baseUrl,
		Token:           apiToken,
		IdsClientId:     idsClientId,
		IdsClientSecret: idsClientSecret,
		IdsScope:        idsScope,
		HttpClient:      retryablehttp.NewClient(),
	}, nil

}

func (c *ApiClient) BuildRequest(method, path string, body interface{}, isUsingIdsToken bool) (*retryablehttp.Request, error) {

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
	if isUsingIdsToken {
		idsToken, _ := c.GetIdsToken()
		req.Header.Set("Authorization", "Bearer "+idsToken["access_token"].(string))
	} else {
		req.Header.Set("Authorization", "TOK:"+c.Token)
	}

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

func (c *ApiClient) GetIdsToken() (map[string]interface{}, error) {
	tokUrl := "https://id-dev.vdms.io/connect/token"
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Add("scope", c.IdsScope)
	data.Add("client_id", c.IdsClientId)
	data.Add("client_secret", c.IdsClientSecret)

	r, _ := http.NewRequest("POST", tokUrl, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Cache-Control", "no-cache")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode()))) //this does not matter, tried without

	var resp *http.Response
	client := &http.Client{}
	resp, err := client.Do(r)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, err
}
