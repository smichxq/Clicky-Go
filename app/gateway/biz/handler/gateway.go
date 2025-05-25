package handler

import (
	"clicky.website/clicky/gateway/idl"
	"context"
	"fmt"
	"net/http"

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
		utils.SendErrResponse(ctx, c, consts.StatusOK, fmt.Errorf("service %s not found", svcName))
		return
	}

	hlog.Info("uri: ", c.URI().String())

	req, err := http.NewRequest(http.MethodGet, c.URI().String(), nil)
	if err != nil {
		panic(err)
	}

	customReq, errr := generic.FromHTTPRequest(req)

	if errr != nil {
		panic(errr)
	}

	// The second parameter is the method name mapped from the HTTP request.
	// It can be an empty string if not applicable.
	resp, err := client.GenericCall(ctx, "", customReq)
	if err != nil {
		hlog.Error("GenericCall failed: ", err)
		panic(err)
	}

	realResp := resp.(*generic.HTTPResponse)

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, realResp)
}
