package handler

import (
	"context"
	"fmt"
	"net/http"

	"clicky.website/clicky/gateway/idl"

	"clicky.website/clicky/gateway/biz/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/pkg/generic"
)

func GenericCall(ctx context.Context, c *app.RequestContext) {
	svcName := c.Param("svc")

	client, exist := idl.SvcMapManagerInstance.GetSvc(svcName)

	if !exist {
		utils.SendErrResponse(ctx, c, consts.StatusNotFound, fmt.Errorf(""))
		return
	}

	hlog.Info("uri: ", c.URI().String())

	req, err := http.NewRequest(http.MethodGet, c.URI().String(), nil)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusNotFound, fmt.Errorf(""))
		hlog.Errorf("build http request fail: ", err)
		return
	}

	customReq, err := generic.FromHTTPRequest(req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusNotFound, fmt.Errorf(""))
		hlog.Errorf("build generic request fail: ", err)
		return
	}

	// The second parameter is the method name mapped from the HTTP request.
	// It can be an empty string if not applicable.
	resp, err := client.GenericCall(ctx, "", customReq)
	if err != nil {
		hlog.Error("GenericCall failed: ", err)
		utils.SendErrResponse(ctx, c, consts.StatusNotFound, fmt.Errorf(""))
		return
	}

	realResp := resp.(*generic.HTTPResponse)

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, realResp)
}
