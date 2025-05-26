package httpclientpool

import (
	"encoding/json"
	"reflect"
	"testing"

	"clicky.website/clicky/security/biz/model"
	"clicky.website/clicky/security/conf"
	"github.com/cloudwego/hertz/pkg/protocol"
)

func TestDo(t *testing.T) {
	Init()
	if HertzClient == nil {
		t.Fatal("HertzClient is not initialized")
		return
	}

	wx := conf.GetConf().WxMini
	c2s := model.Code2SessionReq{
		Base:      "https://api.weixin.qq.com/sns/jscode2session",
		AppId:     wx.AppId,
		Secret:    wx.Secret,
		JsCode:    "test_js_code",
		GrantType: "authorization_code",
	}

	uri, err := c2s.GetReqUrl()
	if err != nil {
		t.Fatalf("failed to get request URL: %v", err)
	}

	type args struct {
		uri    string
		method string
		body   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *protocol.Response
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test Do with valid request",
			args: args{
				uri:    uri,
				method: "GET",
				body:   nil,
			},
			want:    &protocol.Response{}, // Expecting a non-nil response, but actual content will depend on the server response
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Do(tt.args.uri, tt.args.method, tt.args.body)

			var resp *model.Code2SessionResp

			err = json.Unmarshal([]byte(string(got.Body())), &resp)

			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Logf("Response body: %s", string(got.Body()))
				t.Errorf("Do() = %v, want %v", resp, tt.want)
			}
		})
	}
}
