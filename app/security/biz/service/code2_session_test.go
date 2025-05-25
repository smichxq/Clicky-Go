package service

import (
	security "clicky.website/clicky/security/kitex_gen/security"
	"context"
	"testing"
)

func TestCode2Session_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCode2SessionService(ctx)
	// init req and assert value

	req := &security.Code2SessionReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
