package mainnetapi_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
)

type websocketBlockResponse struct {
	Result struct {
		BlockId struct {
			Hash string `json:"hash"`
		} `json:"block_id"`
		Block struct {
			Header struct {
				Height string `json:"height"`
			} `json:"header"`
		} `json:"block"`
	} `json:"result"`
}

func TestBlockexplorerWebsocket(t *testing.T) {
	assert := assert.New(t)

	for _, url := range BeURLs {
		resp, err := sendWebsocketMessage(t, fmt.Sprintf("ws://%s", url), `{"method": "block", "params": ["5"], "id": 1}\n`, nil)
		assert.Errorf(err, "Expected error from invalid websocket for the %s url", url)
		assert.Empty(resp)

		for _, height := range []int{5, 10, 15, 20, 300, 3000} {

			parsedResponse := &websocketBlockResponse{}
			_, err := sendWebsocketMessage(t,
				fmt.Sprintf("wss://%s/websocket", url),
				fmt.Sprintf(`{"method": "block", "params": ["%d"], "id": 1}\n`, height),
				parsedResponse,
			)
			if !assert.NoErrorf(err, "Expected no error from socket for the %s url", url) {
				continue
			}
			assert.NotEmpty(parsedResponse.Result.BlockId.Hash)
			assert.Equal(fmt.Sprintf("%d", height), parsedResponse.Result.Block.Header.Height)
		}
	}

}

func sendWebsocketMessage(t *testing.T, url, message string, response interface{}) (string, error) {
	origin := "http://localhost/"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return "", fmt.Errorf("Failed to dial websocket: %w", err)
	}

	n, err := ws.Write([]byte(message))
	if err != nil {
		return "", fmt.Errorf("failed to write data into websocket: %w", err)
	}
	t.Logf("Sent %d packages to %s", n, url)

	ws.SetReadDeadline(time.Now().Add(time.Second * 5))
	if response == nil {
		var msg = make([]byte, 10240)
		n, err = ws.Read(msg)
		if err != nil {
			return "", fmt.Errorf("failed to read from the websocket: %w", err)
		}

		return string(msg[:n]), nil
	}

	if err := websocket.JSON.Receive(ws, response); err != nil {
		return "", fmt.Errorf("failed to receive json from websocket: %w", err)
	}

	return "", nil
}
