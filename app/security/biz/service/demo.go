package service

import (
	common "clicky.website/clicky/security/kitex_gen/common"
	"context"
)

type DemoService struct {
	ctx context.Context
} // NewDemoService new DemoService
func NewDemoService(ctx context.Context) *DemoService {
	return &DemoService{ctx: ctx}
}

// Run create note info
func (s *DemoService) Run(req *common.EmptyResp) (resp *common.EmptyResp, err error) {
	// Finish your business logic.

	return
}
