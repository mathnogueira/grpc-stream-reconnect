package grpcstreamreconnect_test

import (
	"testing"

	toxiproxy "github.com/Shopify/toxiproxy/v2/client"
	"github.com/stretchr/testify/require"
)

func TestConnectionShouldWork(t *testing.T) {
	client := toxiproxy.NewClient("localhost:8474")
	proxy, err := client.CreateProxy("grpc-server", "localhost:26379", "localhost:26379")
	require.NoError(t, err)
	defer proxy.Delete()

	require.NotNil(t, proxy)
}
