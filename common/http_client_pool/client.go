package httpclientpool

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"clicky.website/clicky/security/biz/model"
	"clicky.website/clicky/security/conf"
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

		// wxmini response
		customRetryIfFunc := func(req *protocol.Request, resp *protocol.Response, err error) bool {
			if resp == nil {
				klog.Errorf("HTTP request failed, uri: %s, method: %s, err: %v", req.RequestURI(), req.Method(), err)
				return true // retry on error
			}

			if resp.StatusCode() != 200 {
				klog.Errorf("HTTP request failed, status code: %d, uri: %s, method: %s", resp.StatusCode(), req.RequestURI(), req.Method())
				return true // retry on non-200 status codes
			}

			klog.Debugf("response body: %s", string(resp.Body()))

			// buzz code
			var buzz_resp *model.Code2SessionResp

			err = json.Unmarshal([]byte(string(resp.Body())), &buzz_resp)
			if err != nil {
				klog.Errorf("unmarshal response failed, err: %v", err)
				return true
			}

			api_retry_codes := conf.GetConf().WxMini.ApiRetryCodes

			for _, retry_code := range api_retry_codes {
				if buzz_resp.ErrCode == retry_code {

					klog.Errorf("wxmini response error, errcode: %d, errmsg: %s, uri: %s, method: %s",
						buzz_resp.ErrCode, buzz_resp.ErrMsg, req.RequestURI(), req.Method())
					return true // retry on specific error codes
				}
			}

			return false
		}
		c.SetRetryIfFunc(customRetryIfFunc)

		if err != nil {
			panic("HertzClient initialization failed: " + err.Error())
		} else {
			HertzClient = c
		}
	})
}

// for code2session
func Do(uri string, method string, body []byte) (*protocol.Response, error) {
	if HertzClient == nil {
		klog.Error("HertzClient is not initialized")
		return nil, errors.New("HertzClient is not initialized")
	}

	req := protocol.AcquireRequest()
	res := protocol.AcquireResponse()
	defer func() {
		protocol.ReleaseRequest(req)
		// Release the response body after use in caller
		// protocol.ReleaseResponse(res)
	}()

	req.SetOptions(
		config.WithRequestTimeout(5000 * time.Millisecond),
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
