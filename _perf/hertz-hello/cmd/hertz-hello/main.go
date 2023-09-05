package main

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/utils"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func main() {
	h := server.Default(server.WithHostPorts("127.0.0.1:8080"))

	// utils.H is a shortcut for map[string]interface{}
	h.GET("/someJSON", func(_ context.Context, c *app.RequestContext) {
		c.JSON(consts.StatusOK, utils.H{"message": "hey", "status": consts.StatusOK})
	})

	h.GET("/moreJSON", func(_ context.Context, c *app.RequestContext) {
		// You also can use a struct
		var msg struct {
			Company  string `json:"company"`
			Location string
			Number   int
		}
		msg.Company = "company"
		msg.Location = "location"
		msg.Number = 123
		// Note that msg.Company becomes "company" in the JSON
		// Will output  :   {"company": "company", "Location": "location", "Number": 123}
		c.JSON(consts.StatusOK, msg)
	})

	h.GET("/pureJson", func(_ context.Context, c *app.RequestContext) {
		c.PureJSON(consts.StatusOK, utils.H{
			"html": "<p> Hello World </p>",
		})
	})

	h.GET("/someData", func(_ context.Context, c *app.RequestContext) {
		c.Data(consts.StatusOK, "text/plain; charset=utf-8", []byte("hello"))
	})

	h.Spin()
}
