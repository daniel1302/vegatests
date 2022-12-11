package mainnetapi_test

import "testing"

var (
	BeURLs = []string{
		"be0.mainnet.vega.xyz",
		"be1.mainnet.vega.xyz",
		// "be.explorer.vega.xyz",
	}
)

func TestQueryForBeOptions(t *testing.T) {
	verifyOptionsResponse(t, BeURLs, "")
	verifyOptionsResponse(t, BeURLs, "query")
	verifyOptionsResponse(t, BeURLs, "rest")
	verifyOptionsResponse(t, BeURLs, "genesis")
	verifyOptionsResponse(t, BeURLs, "blockchain")
	verifyOptionsResponse(t, BeURLs, "validators")
	verifyOptionsResponse(t, BeURLs, "block")
	verifyOptionsResponse(t, BeURLs, "tx_search")
	verifyOptionsResponse(t, BeURLs, "unconfirmed_txs")
	verifyOptionsResponse(t, BeURLs, "websocket")
	verifyOptionsResponse(t, BeURLs, "websocket")
	verifyOptionsResponse(t, BeURLs, "websocket")
}

func TestQueryForBeGet(t *testing.T) {
	verifyGetResponse(t, BeURLs, "grpc")
	verifyGetResponse(t, BeURLs, "rest/transactions")
	verifyGetResponse(t, BeURLs, "rest/transactions/A1C0E3AE3B630C3F0D171623104125E3B29DB552BA3E74D5DE10EB7778C20465")
	verifyGetResponse(t, BeURLs, "genesis")
	verifyGetResponse(t, BeURLs, "blockchain")
	verifyGetResponse(t, BeURLs, "validators")
	verifyGetResponse(t, BeURLs, "block")
	// verifyGetResponse(t, BeURLs, "tx_search") // Currently returns 500
	verifyGetResponse(t, BeURLs, "unconfirmed_txs")
	verifyGetResponse(t, BeURLs, "websocket")
	verifyGetResponse(t, BeURLs, "websocket")
	verifyGetResponse(t, BeURLs, "websocket")
}
