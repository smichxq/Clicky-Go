package service

import (
	"context"

	security "clicky.website/clicky/security/kitex_gen/security"
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

	return &security.Code2SessionResp{req.Code, req.Code}, nil
}
