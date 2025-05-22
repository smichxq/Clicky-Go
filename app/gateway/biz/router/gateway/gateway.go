package gateway

import (
	"context"
	"fmt"
	"net/http"

	"clicky.website/clicky/gateway/biz/handler"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func CustomizedRegister(r *server.Hertz) {
	r.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusOK, "hertz-gateway is running")
	})

	fmt.Println("register gateway")

	registerDynamic(r)
}

// registerGateway registers the router of gateway
func registerDynamic(r *server.Hertz) {
	// group := r.Group("/wxmini")
	// .Use(middleware.GatewayAuth()...)

	r.Any("/:svc/*uri", handler.GenericCall)
}
