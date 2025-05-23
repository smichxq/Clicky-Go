package serversuite

import (
	"clicky.wesite/clicky/common/mtl"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	consulapi "github.com/hashicorp/consul/api"
	prometheus "github.com/kitex-contrib/monitor-prometheus"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
)

type CommonServerSuite struct {
	CurrentServiceName string
	RegistryAddr       string
	ConsulHealthAddr   string
}

func (s CommonServerSuite) Options() []server.Option {
	opts := []server.Option{
		// RPC layer add metadata(HTTP2）
		server.WithMetaHandler(transmeta.ClientHTTP2Handler),
		// set server base info（ServiceName...）
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.CurrentServiceName,
		}),
		// add Prometheus Server Tracer
		server.WithTracer(prometheus.NewServerTracer(
			"",
			"",
			// disable Kitex internal metrics
			prometheus.WithDisableServer(true),
			// register prometheus metrics
			prometheus.WithRegistry(mtl.Registry),
		),
		),
		// kitex-opentelemetry
		server.WithSuite(tracing.NewServerSuite()),
	}

	// initialize consul registry
	r, err := consul.NewConsulRegister(s.RegistryAddr, consul.WithCheck(&consulapi.AgentServiceCheck{
		HTTP:                           "http://" + s.ConsulHealthAddr + "/health",
		Interval:                       "1s",
		Timeout:                        "1s",
		DeregisterCriticalServiceAfter: "1s",
	}))
	if err != nil {
		panic(err)
	}

	// register component in opts
	opts = append(opts, server.WithRegistry(r))

	return opts
}
