package suite

import (
	"context"

	"clicky.website/clicky/gateway/conf"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/fallback"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
)

func circuitbreakSuite(metricsKeyCfg map[string]circuitbreak.CBConfig) []client.Option {
	// circuit breaker suite
	cbs := circuitbreak.NewCBSuite(func(ri rpcinfo.RPCInfo) string {
		// "fromServiceName/ToServiceName/method"
		return circuitbreak.RPCInfo2Key(ri)
	})

	// dynamic store kv
	for k, v := range metricsKeyCfg {
		cbs.UpdateServiceCBConfig(k, v)
	}

	opts := []client.Option{
		client.WithCircuitBreaker(cbs),
		client.WithFallback(
			fallback.NewFallbackPolicy(
				fallback.UnwrapHelper(func(ctx context.Context, req, resp interface{}, err error) (fbResp interface{}, fbErr error) {
					if err == nil {
						return resp, nil
					}

					/* CircuitBreaker happend */
					return &generic.HTTPResponse{}, nil
				}),
			),
		),
	}

	return opts
}

func baseSuite() []client.Option {
	return []client.Option{
		client.WithResolver(*conf.ConsulResolver),
		client.WithTransportProtocol(transport.TTHeader),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
	}
}

// metricsKey: "fromServiceName/ToServiceName/method"
func GenericSuite(metricsKeyCfg map[string]circuitbreak.CBConfig) []client.Option {
	return append(baseSuite(), circuitbreakSuite(metricsKeyCfg)...)
}
