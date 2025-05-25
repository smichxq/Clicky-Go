package main

import (
	"clicky.website/clicky/security/biz/service"
	common "clicky.website/clicky/security/kitex_gen/common"
	security "clicky.website/clicky/security/kitex_gen/security"
	"context"
)

// SecurityImpl implements the last service interface defined in the IDL.
type SecurityImpl struct{}

// Code2Session implements the SecurityImpl interface.
func (s *SecurityImpl) Code2Session(ctx context.Context, req *security.Code2SessionReq) (resp *security.Code2SessionResp, err error) {
	resp, err = service.NewCode2SessionService(ctx).Run(req)

	return resp, err
}

// Demo implements the SecurityImpl interface.
func (s *SecurityImpl) Demo(ctx context.Context, req *common.EmptyResp) (resp *common.EmptyResp, err error) {
	resp, err = service.NewDemoService(ctx).Run(req)

	return resp, err
}
