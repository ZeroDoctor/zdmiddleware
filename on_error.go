package zdmiddleware

import (
	"fmt"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

func DumpRequestOnError() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request
		ctx.Next()
		if ctx.Writer.Status() < 200 || ctx.Writer.Status() > 299 {
			data, err := httputil.DumpRequest(req, true)
			if err != nil {
				fmt.Printf("[%s] failed to dump request [error=%s]", red.Render("ERROR"), err.Error())
				return
			}

			fmt.Println(string(data))
		}
	}
}

func DumpRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := ctx.Request
		data, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Printf("[%s] failed to dump request [error=%s]", red.Render("ERROR"), err.Error())
			return
		}

		fmt.Println(string(data))
		ctx.Next()
	}
}
