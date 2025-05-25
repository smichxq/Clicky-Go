package service

import (
	common "clicky.website/clicky/security/kitex_gen/common"
	"context"
	"testing"
)

func TestDemo_Run(t *testing.T) {
	ctx := context.Background()
	s := NewDemoService(ctx)
	// init req and assert value

	req := &common.EmptyResp{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
