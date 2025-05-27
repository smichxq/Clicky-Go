package model

import (
	"testing"

	"clicky.website/clicky/security/conf"
)

func TestCode2SessionReq_GetReqUrl(t *testing.T) {
	type fields struct {
		Base      string
		AppId     string
		Secret    string
		JsCode    string
		GrantType string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "",
			fields: fields{
				Base:      "https://api.weixin.qq.com/sns/jscode2session",
				AppId:     conf.GetConf().WxMini.AppId,
				Secret:    conf.GetConf().WxMini.Secret,
				JsCode:    "789",
				GrantType: "authorization_code",
			},
			want:    "https://api.weixin.qq.com/sns/jscode2session?appid=123&grant_type=111&js_code=789&secret=456",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c2s := &Code2SessionReq{
				Base:      tt.fields.Base,
				AppId:     tt.fields.AppId,
				Secret:    tt.fields.Secret,
				JsCode:    tt.fields.JsCode,
				GrantType: tt.fields.GrantType,
			}
			got, err := c2s.GetReqUrl()
			if (err != nil) != tt.wantErr {
				t.Errorf("Code2SessionReq.GetReqUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Code2SessionReq.GetReqUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
