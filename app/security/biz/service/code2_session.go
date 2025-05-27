package service

import (
	"context"
	"encoding/json"

	"clicky.website/clicky/security/biz/model"
	"clicky.website/clicky/security/conf"
	security "clicky.website/clicky/security/kitex_gen/security"
	httpclientpool "clicky.wesite/clicky/common/http_client_pool"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type Code2SessionService struct {
	ctx context.Context
} // NewCode2SessionService new Code2SessionService
func NewCode2SessionService(ctx context.Context) *Code2SessionService {
	return &Code2SessionService{ctx: ctx}
}

// Run create note info
func (s *Code2SessionService) Run(req *security.Code2SessionReq) (resp *security.Code2SessionResp, err error) {
	// Finish your business logic.

	return getSession(req, s)
}

// access wx server to get login code
func getSession(req *security.Code2SessionReq, s *Code2SessionService) (*security.Code2SessionResp, error) {
	_, span := otel.Tracer("getSession").Start(s.ctx, "security.security/Code2Session/getSession")
	defer span.End()
	code2session := model.Code2SessionReq{
		Base:      "https://api.weixin.qq.com/sns/jscode2session",
		AppId:     conf.GetConf().WxMini.AppId,
		Secret:    conf.GetConf().WxMini.Secret,
		JsCode:    *req.Code,
		GrantType: "authorization_code",
	}

	uri, err := code2session.GetReqUrl()
	if err != nil {
		klog.Errorf("get request url failed", err)
		return nil, err
	}
	uri = "http://192.168.3.6:4523/m2/5010042-4669468-default/261010813"
	raw_resp, err := httpclientpool.Do(uri, "GET", nil)
	if err != nil {
		klog.Errorf("http request failed, uri: %s, err: %v", uri, err)
		return nil, err
	}

	// klog.Debugf("response body: %s", string(raw_resp.Body()))

	// var resp model.Code2SessionResp
	resp := &model.Code2SessionResp{}

	err = json.Unmarshal(raw_resp.Body(), resp)
	if err != nil {
		klog.Errorf("unmarshal response failed, err: %v", err)
		return nil, err
	}

	// release the response body
	defer func() {
		protocol.ReleaseResponse(raw_resp)
	}()

	key := uuid.NewString()

	token := uuid.NewString()

	return &security.Code2SessionResp{Key: &key, Token: &token}, nil
}
