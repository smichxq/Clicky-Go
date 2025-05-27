package httpclientpool

import (
	"encoding/json"
	"reflect"
	"testing"

	"clicky.website/clicky/security/biz/model"
	"github.com/cloudwego/hertz/pkg/protocol"
)

func TestDo(t *testing.T) {
	Init()
	if HertzClient == nil {
		t.Fatal("HertzClient is not initialized")
		return
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
				uri:    "http://192.168.3.6:4523/m2/5010042-4669468-default/261010813",
				method: "GET",
				body:   nil,
			},
			want:    &protocol.Response{}, // Expecting a non-nil response, but actual content will depend on the server response
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Do(tt.args.uri, tt.args.method, tt.args.body)
			if err != nil {
				t.Errorf("Do() error = %v", err)
			}

			t.Logf("Response status code: %d", got.StatusCode())
			t.Logf("Response: %d", got.Body())

			var resp *model.Code2SessionResp

			err = json.Unmarshal([]byte(string(got.Body())), &resp)

			t.Logf("Response body: %s", string(got.Body()))

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
