package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if ctx.Request.Method == http.MethodOptions {
		ctx.Status(http.StatusNoContent)
		return
	}

	ctx.Next()
}
