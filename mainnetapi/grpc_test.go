package mainnetapi_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"

	vegaapipb "code.vegaprotocol.io/vega/protos/vega/api/v1"
)

func TestDataNoderpcEndpoint(t *testing.T) {
	assert := assert.New(t)

	for _, url := range append(ApiURLs, BeURLs...) {
		resp, err := callStatistics(fmt.Sprintf("%s:13007", url), 5*time.Second)
		if !assert.NoErrorf(err, "Failed to dial the %s GRPC node", url) {
			return
		}

		if assert.NotNil(resp.Statistics) {
			assert.Containsf(resp.Statistics.ChainId, "vega-mainnet", "Mainnet chain id seems incorrect")
			assert.Greaterf(resp.Statistics.BlockHeight, uint64(100), "mainnet block height should be higher than 100")
		}
	}
}
func TestCoreGrpcEndpoint(t *testing.T) {
	assert := assert.New(t)

	for _, url := range append(ApiURLs, BeURLs...) {
		resp, err := callStatistics(fmt.Sprintf("%s:13002", url), 5*time.Second)
		if !assert.NoErrorf(err, "Failed to dial the %s GRPC node", url) {
			return
		}

		if assert.NotNil(resp.Statistics) {
			assert.Containsf(resp.Statistics.ChainId, "vega-mainnet", "Mainnet chain id seems incorrect")
			assert.Greaterf(resp.Statistics.BlockHeight, uint64(100), "mainnet block height should be higher than 100")
		}
	}
}

func callStatistics(addr string, timeout time.Duration) (*vegaapipb.StatisticsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to dial gRPC node %s: %w", addr, err)
	}

	if conn.GetState() != connectivity.Ready {
		return nil, fmt.Errorf("grpc connection for data node is not ready: connection not ready")
	}

	c := vegaapipb.NewCoreServiceClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()

	response, err := c.Statistics(ctx, &vegaapipb.StatisticsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node statistics: %w", err)
	}

	return response, nil

}
