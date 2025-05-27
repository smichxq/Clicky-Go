package service

import (
	"context"
	"reflect"
	"testing"

	security "clicky.website/clicky/security/kitex_gen/security"
	httpclientpool "clicky.wesite/clicky/common/http_client_pool"
)

func TestCode2Session_Run(t *testing.T) {
	httpclientpool.Init()
	ctx := context.Background()
	s := NewCode2SessionService(ctx)
	// init req and assert value

	tk := "tokennnnnn"
	req := &security.Code2SessionReq{Code: &tk}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
}

func Test_getSession(t *testing.T) {
	httpclientpool.Init()
	code := "test_js_code" // Replace with a valid code for testing
	type args struct {
		req *security.Code2SessionReq
		s   *Code2SessionService
	}
	tests := []struct {
		name    string
		args    args
		want    *security.Code2SessionResp
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test getSession with valid request",
			args: args{
				req: &security.Code2SessionReq{Code: &code}, // Replace nil with a valid code if needed
				s:   NewCode2SessionService(context.Background()),
			},
			want:    &security.Code2SessionResp{}, // Expecting a non-nil response, but actual content will depend on the server response
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSession(tt.args.req, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSession() = %v, want %v", got, tt.want)
			}
		})
	}
}
