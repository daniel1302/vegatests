package mainnetapi_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	UserAgent       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:107.0) Gecko/20100101 Firefox/107.0"
	Accept          = "*/*"
	ContentTypeJSON = "application/json"
)

var (
	ApiURLs = []string{
		"api0.mainnet.vega.xyz",
		"api1.mainnet.vega.xyz",
		"api2.mainnet.vega.xyz",
		"api3.mainnet.vega.xyz",
	}

	exampleGQLQuery = `{"operationName":"NetworkStats","variables":{},"query":"query NetworkStats {\nnodeData {\nstakedTotal\n}\nstatistics {\nstatus\n}\n}\n"}`
)

func TestQueryForAPIOptions(t *testing.T) {
	verifyOptionsResponse(t, ApiURLs, "")
	verifyOptionsResponse(t, ApiURLs, "graphql")
	verifyOptionsResponse(t, ApiURLs, "query")
	verifyOptionsResponse(t, ApiURLs, "statistics")
	verifyOptionsResponse(t, ApiURLs, "graphql")
}

func TestQueryForAPIPost(t *testing.T) {
	verifyPostResponse(t, ApiURLs, "query", exampleGQLQuery)
}

func TestQueryForAPIGet(t *testing.T) {
	verifyGetResponse(t, ApiURLs, "graphql")
	verifyGetResponse(t, ApiURLs, "query")
	verifyGetResponse(t, ApiURLs, "statistics")
	verifyGetResponse(t, ApiURLs, "graphql")
}

func verifyOptionsResponse(t *testing.T, hosts []string, endpoint string) {
	expectedHeaders := []string{
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Origin",
		"Access-Control-Max-Age",
		"X-Vega-Node-Id",
	}
	assert := assert.New(t)

	for _, apiNode := range hosts {
		url := fmt.Sprintf("https://%s/%s", apiNode, endpoint)
		optionsHeaders := map[string]string{
			"access-control-request-headers": "content-type",
			"access-control-request-method":  "POST",
		}

		t.Logf("Sending OPTION request to %s", url)
		resp, err := call(url, "OPTIONS", optionsHeaders, "")
		if !assert.NoError(err) {
			return
		}

		if !assert.Equalf(http.StatusNoContent, resp.StatusCode, "Invalid response code for the %s endpoint", endpoint) {
			return
		}

		for _, header := range expectedHeaders {
			if assert.Containsf(resp.Header, header, "The '%s' header received from '%s' is missing for OPTIONS request", header, url) {
				assert.Lenf(resp.Header[header], 1, "Expected exactly 1 value returned for the '%s' header from the '%s' URL for OPTIONS request", header, url)
			}
		}
	}
}

func verifyGetResponse(t *testing.T, hosts []string, endpoint string) {
	expectedHeaders := []string{
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Origin",
		"X-Vega-Node-Id",
	}
	assert := assert.New(t)

	getHeaders := map[string]string{
		"Accept-Language": "en-US,en;q=0.5",
		"Accept-Encoding": "gzip, deflate, br",
		"Referer":         "https://stats.vega.trading/",
		"Origin":          "https://stats.vega.trading",
		"DNT":             "1",
		"Connection":      "keep-alive",
		"Pragma":          "no-cache",
		"Cache-Control":   "no-cache",
		"TE":              "trailers",

		"Sec-Fetch-Dest": "empty",
		"Sec-Fetch-Mode": "cors",
		"Sec-Fetch-Site": "cross-site",
		"Accept":         Accept,
		"User-Agent":     UserAgent,
		"Content-Type":   ContentTypeJSON,
	}

	for _, apiNode := range hosts {
		url := fmt.Sprintf("https://%s/%s", apiNode, endpoint)

		t.Logf("Sending GET request to %s", url)
		resp, err := call(url, "GET", getHeaders, "")
		if !assert.NoError(err) {
			return
		}

		if !assert.Equalf(http.StatusOK, resp.StatusCode, "Invalid response code for the %s endpoint", endpoint) {
			return
		}

		for _, header := range expectedHeaders {
			if assert.Containsf(resp.Header, header, "The '%s' header received from '%s' is missing for GET request", header, url) {
				assert.Lenf(resp.Header[header], 1, "Expected exactly 1 value returned for the '%s' header from the '%s' URL for GET request", header, url)
			}
		}
	}
}

func verifyPostResponse(t *testing.T, hosts []string, endpoint, postData string) {
	expectedHeaders := []string{
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Origin",
		"X-Vega-Node-Id",
	}
	assert := assert.New(t)

	postHeaders := map[string]string{
		"Referer":      "https://explorer.vega.xyz/",
		"Origin":       "https://explorer.vega.xyz",
		"Accept":       Accept,
		"User-Agent":   UserAgent,
		"Content-Type": ContentTypeJSON,
	}

	for _, apiNode := range hosts {
		url := fmt.Sprintf("https://%s/%s", apiNode, endpoint)

		t.Logf("Sending POST request to %s", url)
		resp, err := call(url, "POST", postHeaders, postData)
		if !assert.NoError(err) {
			return
		}

		if !assert.Equalf(http.StatusOK, resp.StatusCode, "Invalid response code for the %s endpoint", endpoint) {
			return
		}

		for _, header := range expectedHeaders {
			if assert.Containsf(resp.Header, header, "The '%s' header received from '%s' is missing", header, url) {
				assert.Lenf(resp.Header[header], 1, "Expected exactly 1 value returned for the '%s' header from the '%s' URL", header, url)
			}
		}
	}
}

func call(url, method string, headers map[string]string, body string) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	if len(body) > 0 {
		req.Body = io.NopCloser(strings.NewReader(body))
	}

	for headerKey, headerValue := range headers {
		req.Header.Set(headerKey, headerValue)
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	return res, nil
}
