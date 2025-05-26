package httpclientpool

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/client/retry"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/kitex/pkg/klog"
)

var (
	HertzClient *client.Client
	once        sync.Once
)

func Init() {
	once.Do(func() {
		clientCfg := &tls.Config{
			InsecureSkipVerify: true, // ignore TLS certificate verification
		}
		c, err := client.NewClient(
			client.WithTLSConfig(clientCfg),         // use custom TLS config
			client.WithDialer(standard.NewDialer()), // use standard dialer
			client.WithRetryConfig(
				// Max retry times
				retry.WithMaxAttemptTimes(3),
				// first retry delay
				retry.WithInitDelay(5*time.Millisecond),
				retry.WithMaxDelay(10*time.Millisecond),
			),
		)

		if err != nil {
			klog.Errorf("HertzClient initialization failed: %v", err)
			panic("HertzClient initialization failed: " + err.Error())
		} else {
			HertzClient = c
		}
	})
}

func Do(uri string, method string, body []byte) (*protocol.Response, error) {
	req := protocol.AcquireRequest()
	res := protocol.AcquireResponse()
	defer func() {
		protocol.ReleaseRequest(req)
		protocol.ReleaseResponse(res)
	}()

	req.SetOptions(
		config.WithRequestTimeout(5),
	)

	req.SetMethod(method)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	req.SetRequestURI(uri)
	req.SetBody(body)
	err := HertzClient.Do(context.Background(), req, res)
	if err != nil {
		klog.Errorf("HTTP request failed, uri: %s, method: %s, err: %v", uri, method, err)
		return nil, err
	}

	return res, nil
}
