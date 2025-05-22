package handler

import (
	"context"
	"fmt"
	"net/http"

	"clicky.website/clicky/gateway/biz/idl"
	"clicky.website/clicky/gateway/biz/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"
)

func GenericCall(ctx context.Context, c *app.RequestContext) {
	svcName := c.Param("svc")

	client, exist := idl.SvcMapManagerInstance.GetSvc(svcName)

	if !exist {
		utils.SendErrResponse(ctx, c, consts.StatusOK, fmt.Errorf("service %s not found", svcName))
		return
	}

	fmt.Println("uri:", c.URI().String())

	reqq, errr := http.NewRequest(http.MethodGet, c.URI().String(), nil)
	if errr != nil {
		panic(errr)
	}

	customReq, errr := generic.FromHTTPRequest(reqq)

	if errr != nil {
		panic(errr)
	}

	// The second parameter is the method name mapped from the HTTP request.
	// It can be an empty string if not applicable.
	resp, errr := client.GenericCall(ctx, "", customReq)
	if errr != nil {
		klog.Errorf("GenericCall failed: %v", errr)
		panic(errr)
	}

	realResp := resp.(*generic.HTTPResponse)

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, realResp)
}
